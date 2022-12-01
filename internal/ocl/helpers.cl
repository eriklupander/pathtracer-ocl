
// findFirstIntersectionCloserThan returns true if an intersection is found that's closer than minT and which is not
// ignoreIndex. This function is very similar to findClosestIntersection and should be refactored into something more
// common. However, the triangle intersection code is somewhat tricky to generalize.
inline bool findFirstIntersectionCloserThan(__global object *objects, unsigned int numObjects, __global group *groups, __global triangle *triangles, double4 rayOrigin, double4 rayDirection, double minT, unsigned int ignoreIndex) {

    // ----------------------------------------------------------
    // Loop through scene objects in order to find intersections
    // ----------------------------------------------------------
    for (unsigned int j = 0; j < numObjects; j++) {
        if (j == ignoreIndex) {
            continue;
        }
        long objType = objects[j].type;
        //  translate our ray into object space by multiplying ray pos and dir
        //  with inverse object matrix
        double4 tRayOrigin = mul(objects[j].inverse, rayOrigin);
        double4 tRayDirection = mul(objects[j].inverse, rayDirection);

        // Intersection code
        if (objType == 0) { // PLANE - intersect transformed ray with plane
            double t = intersectPlane(tRayOrigin, tRayDirection);
            if (t > 0.0 && t < minT) {
                return true;
            }
        } else if (objType == 1) { // SPHERE

            // finally, find the intersection distances on our ray.
            double2 t = intersectSphere(tRayOrigin, tRayDirection);
            // double t2 = (-b + sqrt(discriminant)) / (2*a); // add back in
            // when we do refraction
            if (t.x > 0.0 && t.x < minT) {
                return true;
            }
//            if (t.y > 0.0 && t.y < minT) {
//                return true;
//            }
        } else if (objType == 2) { // CYLINDER
            double4 out = intersectCylinder(tRayOrigin, tRayDirection, objects[j]);
            for (unsigned int a = 0; a < 4; a++) {
                if (out[a] > 0.0 && out[a] < minT) {
                    return true;
                }
            }
        } else if (objType == 3) { // BOX
            double2 out = intersectCube(tRayOrigin, tRayDirection);

            // assign intersections
            if (out.x > 0.0 && out.x < minT) {
               return true;
            }
            if (out.y > 0.0 && out.y < minT) {
               return true;
            }

        } else if (objType == 4) { // GROUPS

            // Group with triangles experiment
            // Groups MUST have their bounds computed. Start by checking if ray intersects bounds.
            // Remember: At this point in the code, the group's transform has already modified the ray.
            // However, the cube intersection is based on transform/rotate/scale to unit cube. Our BB does not
            // really work that way...
            // Note!! BB must have extent in all 3-axises. I.e two triangles forming a wall facing the Z axis will have 0
            // depth which breaks the intersect code. (typically, use this for models that's rarely flat, or fake something if 0.)
            // Using this BB only reduces teapot with 8 samples from 3m29.753546781s to 31.606680099s.
            // Further, adding the BB check for each node in the tree further reduces the time taken to 4.037422895s
            if (!intersectRayWithBox(tRayOrigin, tRayDirection, objects[j].bbMin, objects[j].bbMax)) {
                // skipped++;
                continue;
            }
            // hit++;

            // If the "object" BB was intersected, we take a look at the "object's" groupOffset. If > -1, we
            // need to set up a local stack to traverse the group hierarchy
            if (objects[j].childCount > 0) {

                // this is somewhat ugly, but since a "parent" obj (from objects) may have up to 64 children
                // (references to indexes in "groups"), we must use a for-statement here.
                for (int childIndex = 0; childIndex < objects[j].childCount; childIndex++) {
                    // START PSUEDO-RECURSIVE CODE
                    // 1) Create an empty stack. (move to top later)
                    int stack[64] = {0};

                    // Stack index, i.e. current "depth" of stack
                    int currentSIndex = 0;

                    // Tree index, i.e. which "node index" we're currently processing
                    int currentNodeIndex = objects[j].children[childIndex];

                    // Initialize current node as root. Note the ugly code to get a pointer to the current node...
                    group root = groups[currentNodeIndex];
                    group *current = &root;

                    for (; current != 0 || currentSIndex > -1;) {
                        for (; current != 0 && intersectRayWithBox(tRayOrigin, tRayDirection, current->bbMin, current->bbMax);) {

                            // Iterate over all triangles and record triangle/ray intersections...
                            for (int n = current->triOffset; n < current->triOffset + current->triCount; n++) {

                                double4 dirCrossE2 = cross(tRayDirection, triangles[n].e2);
                                double determinant = dot(triangles[n].e1, dirCrossE2);
                                if (fabs(determinant) < EPSILON) {
                                    continue;
                                }

                                // Triangle misses over P1-P3 edge
                                double f = 1.0 / determinant;
                                double4 p1ToOrigin = tRayOrigin - triangles[n].p1;
                                double u = f * dot(p1ToOrigin, dirCrossE2);
                                if (u < 0 || u > 1) {
                                    continue;
                                }

                                double4 originCrossE1 = cross(p1ToOrigin, triangles[n].e1);
                                double v = f * dot(tRayDirection, originCrossE1);
                                if (v < 0 || (u + v) > 1) {
                                    continue;
                                }
                                double t = f * dot(triangles[n].e2, originCrossE1);
                                if (t > 0.0 && t < minT) {
                                    return true;
                                }
                            }

                            // Push the current node index to the Stack, i.e. add at current index and then increment the stack depth.
                            stack[currentSIndex] = currentNodeIndex;
                            currentSIndex++;

                            // if the left child is populated (i.e. > -1), update currentNodeIndex with left child and
                            // update the pointer to the current node
                            if (current->children[0] > 0) {
                                currentNodeIndex = current->children[0];
                                root = groups[current->children[0]];
                                current = &root;
                            } else {
                                // If no left child, mark current as nil, so we can exit the inner for.
                                current = 0;
                            }
                        } // exit of inner for loop, i.e. carry on to the right side

                        // We pop our stack by decrementing (remember, the last iteration above resulting an increment, but no push. (Fix?)
                        currentSIndex--;
                        if (currentSIndex == -1) {
                            goto done;
                        }

                        // get the popped item by fetching the node index from the current stack index.
                        root = groups[stack[currentSIndex]];
                        current = &root;

                        // we're done with the left subtree, check if there's a right-hand node.
                        if (current->children[1] > 0) {
                            // if there's a right-hand node, update the node index and the current node.
                            currentNodeIndex = current->children[1];
                            root = groups[current->children[1]];
                            current = &root;
                        } else {
                            // if no right-hand side, set current to nil. In a binary tree, we should
                            // always get a right side if we got a left side...
                            current = 0;
                        }
                    }
                    // END PSUEDO-RECURSIVE CODE
                done:
                    current = 0;
                }
            }
        }
    }
    return false;
}
