__constant double PI = 3.14159265359f;
__constant unsigned int MAX_BOUNCES = 4;

typedef struct __attribute__((packed)) tag_camera {
    int width;              // 4 bytes
    int height;             // 4 bytes
    double fov;             // 8 bytes
    double pixelSize;       // 8 bytes
    double halfWidth;       // 8 bytes
    double halfHeight;      // 8 bytes
    double aperture;        // 8 bytes
    double focalLength;     // 8 bytes
    double16 inverse;       // 128 bytes
    char padding[72];       // 72 bytes
} camera;

typedef struct tag_ray {
  double4 origin;
  double4 direction;
} ray;

typedef struct __attribute__((packed)) tag_group {
	double4 bbMin;      // 32 bytes
	double4 bbMax;      // 32 bytes
	double4 color;      // 32 bytes
	double4 emission;   // 32 bytes
	int triOffset;      // 4 bytes
	int triCount;       // 4 bytes
	int childGroupCount;// 4 bytes, should always be 2 or 0
	int children[16];   // 64 bytes, allow up to 16 subgroups.
	char padding[52];   // padding, 52 bytes (can be used as a label)
	// Total 256 bytes
} group;

typedef struct __attribute__((packed)) tag_object {
  double16 transform;        // 128 bytes 16x4
  double16 inverse;          // 128 bytes
  double16 inverseTranspose; // 128 bytes
  double4 color;             // 32 bytes
  double4 emission;          // 32 bytes
  double refractiveIndex;    // 8 bytes
  long type;                 // 8 bytes
  double minY;               // 8 bytes. Used for cylinders.
  double maxY;               // 8 bytes. Used for cylinders.
  double reflectivity;       // 8 bytes
  double padding3;           // 8 bytes
  double padding4;           // 8 bytes                      // 512
  double4 bbMin;             // 32 bytes
  double4 bbMax;             // 32 bytes                     // 576
  int groupOffset;           // 4 bytes. Used for groups to know which "group" that's the root group.
  char padding5[452];        // 452 bytes                    // 1024
} object;

typedef struct tag_intersection {
  unsigned int objectIndex;
  double t;
  double t2;
} intersection;

typedef struct tag_bounce {
  double4 point;
  double cos;
  double4 color;
  double4 emission;
  // diffuse         bool
  // refractiveIndex float64
} bounce;

typedef struct __attribute__((packed)) tag_triangle {
    double4 p1;       // 32 bytes
    double4 p2;       // 32 bytes
    double4 p3;       // 32 bytes
    double4 e1;       // 32 bytes
    double4 e2;       // 32 bytes
    double4 n;        // 32 bytes
    double4 n1;       // 32 bytes
    double4 n2;       // 32 bytes
    double4 n3;       // 32 bytes
    char padding[224];// 224 bytes
} triangle;           // 512 total

inline double maxX(double a, double b, double c) {
    return max(max(a, b), c);
}
inline double minX(double a, double b, double c) {
    return min(min(a, b), c);
}

inline double2 checkAxis(double origin, double direction, double minBB, double maxBB) {
    double2 out = (double2){0,0};
    double tminNumerator = minBB - origin; //-1.0 - origin;
    double tmaxNumerator = maxBB - origin; //1.0 - origin;
    if (fabs(direction) >= 0.0001) {
      out.x = tminNumerator / direction;
      out.y = tmaxNumerator / direction;
    } else {
      out.x = tminNumerator * HUGE_VAL;
      out.y = tmaxNumerator * HUGE_VAL;
    }
    if (out.x > out.y) {
      // swap
      double temp = out.x;
      out.x = out.y;
      out.y = temp;
    }
    return out;
}

inline bool intersectRayWithBox(double4 tRayOrigin, double4 tRayDirection, double4 bbMin, double4 bbMax) {
	// There is supposed  to be a way to optimize this for fewer checks by looking at early values.
	double2 xt = checkAxis(tRayOrigin.x, tRayDirection.x, bbMin.x, bbMax.x);
    double2 yt = checkAxis(tRayOrigin.y, tRayDirection.y, bbMin.y, bbMax.y);
    double2 zt = checkAxis(tRayOrigin.z, tRayDirection.z, bbMin.z, bbMax.z);

    // We use double2 as a poor-mans multi-valued return.
    // If the largest of the min values is greater smallest max value...
    double tmin = maxX(xt.x, yt.x, zt.x); // x == min
    double tmax = minX(xt.y, yt.y, zt.y); // y == max
	return tmin < tmax;
}

inline bool checkCap(double4 origin, double4 direction, double t) {
	double x = origin.x + t*direction.x;
	double z = origin.z + t*direction.z;
	return pow(x, 2) + pow(z, 2) <= 1.0;
}

inline double2 intersectCaps(double4 origin, double4 direction, double minY, double maxY) {
	// !c.closed removed
    if (fabs(direction.y) < 0.0001) {
		return (double2)(0.0,0.0);
	}

    double2 retVal = (double2)(0.0,0.0);

	// check for an intersection with the lower end cap by intersecting
	// the ray with the plane at y=cyl.minimum
	double t1 = (minY - origin.y) / direction.y;
	if (checkCap(origin, direction, t1)) {
		retVal.x = t1;
	}

	// check for an intersection with the upper end cap by intersecting
	// the ray with the plane at y=cyl.maximum
	double t2 = (maxY - origin.y) / direction.y;
	if (checkCap(origin, direction, t2)) {
		retVal.y = t2;
	}
    return retVal;
}

// from https://stackoverflow.com/a/50665114
inline static float noise3D(float x, float y, float z) {
  float ptr = 0.0f;
  return fract(sin(x * 112.9898f + y * 179.233f + z * 237.212f) * 43758.5453f,
               &ptr);
}

// randomVectorInHemisphere is based on
// https://raytracey.blogspot.com/2016/11/opencl-path-tracing-tutorial-2-path.html
//
// but adapted to use another rand function and double4 instead of float4. The
// thing is that using this func for diffuse surfaces produces a good and
// balanced result in the final image, while using the randomConeInHemisphere
// func translated from Hunter Loftis PathTracer produces overexposed
// highlights. I think randomConeInHemisphere distributes the rays more
// "cone-like" while this one distributes them better across the entire
// hemisphere, which is what we want for strictly diffuse surfaces.
inline double4 randomVectorInHemisphere(double4 normalVec, double x, double y,
                                        double z) {
  double rand1 = 2.0 * PI * noise3D(x, y, z);
  double rand2 = noise3D(y, z, x);
  double rand2s = sqrt(rand2);

  /* create a local orthogonal coordinate frame centered at the hitpoint */
  double4 axis;
  if (fabs(normalVec.x) > 0.1) {
    axis = (double4)(0.0, 1.0, 0.0, 0.0);
  } else {
    axis = (double4)(1.0, 0.0, 0.0, 0.0);
  }
  double4 u = normalize(cross(axis, normalVec));
  double4 v = cross(normalVec, u);

  /* use the coordinate frame and random numbers to compute the next ray
   * direction */
  return u * cos(rand1) * rand2s + v * sin(rand1) * rand2s +
         normalVec * sqrt(1.0 - rand2);
}

// mul multiplies the vec by the matrix, producing a new vector.
inline double4 mul(double16 mat, double4 vec) {
  double4 elem1 = mat.s0123 * vec;
  double4 elem2 = mat.s4567 * vec;
  double4 elem3 = mat.s89AB * vec;
  double4 elem4 = mat.sCDEF * vec;
  return (double4)(elem1.x + elem1.y + elem1.z + elem1.w,
                   elem2.x + elem2.y + elem2.z + elem2.w,
                   elem3.x + elem3.y + elem3.z + elem3.w,
                   elem4.x + elem4.y + elem4.z + elem4.w);
}

inline ray rayForPixel(unsigned int x, unsigned int y, camera cam, float rndX, float rndY) {
	double4 pointInView = {0.0, 0.0, -1.0, 1.0};
	double4 originPoint =  {0.0, 0.0, 0.0, 1.0};
	double xOffset = cam.pixelSize * ((double)x + rndX);
    double yOffset = cam.pixelSize * ((double)y + rndY);

	// this feels a little hacky but actually works.
	pointInView.x = cam.halfWidth - xOffset;
	pointInView.y = cam.halfHeight - yOffset;

    double4 pixel = mul(cam.inverse, pointInView);
    double4 origin = mul(cam.inverse, originPoint);

    double4 subVec = pixel - origin;
    double4 direction = normalize(subVec);

    // if DoF...
    if (cam.aperture != 0) {

        double4 pos = origin + direction*cam.focalLength;
        double4 newOrigin={};
        newOrigin.x = origin.x + (-cam.aperture + rndY*cam.aperture*2);
        newOrigin.y = origin.y + (-cam.aperture + rndX*cam.aperture*2);
        newOrigin.z = origin.z;
        newOrigin.w = 1.0;
        direction = pos - newOrigin;
        origin = newOrigin;
    }
    ray r = {origin, direction};
    return r;
}

__kernel void trace(__global object *objects,
                    const unsigned int numObjects,
                    __global triangle *triangles,
                    __global group *groups,
                    __global double *output,
                    __global double *seedX,
                    const unsigned int samples,
                    __global camera *cam,
                    const unsigned int yOffset) {

  //int skipped = 0;
  //int hit = 0;
  double colorWeight = 1.0 / samples;
  int i = get_global_id(0);
  float fgi = seedX[i] / numObjects;
  float fgi2 = seedX[i] / samples;
  double4 originPoint = (double4)(0.0f, 0.0f, 0.0f, 1.0f);
  double4 colors = (double4)(0, 0, 0, 0);

  // get current x,y coordinate from i given image width
  unsigned int x = i % cam->width;
  unsigned int y = yOffset + i / cam->width;

  for (unsigned int n = 0; n < samples; n++) {
    // For each sample, compute a new ray cast through the target (x,y) pixel with random offset within the pixel.
    ray r = rayForPixel(x, y, *cam, noise3D(fgi, n, fgi2), noise3D(fgi, fgi2, n));
    double4 rayOrigin = r.origin;
    double4 rayDirection = r.direction;


    unsigned int actualBounces = 0;
    // Each ray may bounce up to 5 times
    bounce bounces[5] = {};
    for (unsigned int b = 0; b < MAX_BOUNCES; b++) {

      // track up to 16 intersections per ray.
      double intersections[16] = {0};  // t of an intersection
      unsigned int xsObjects[16] = {0}; // index maps to each xs above, value to objects
      double4 xsTriangle[16] = {0};
      // ----------------------------------------------------------
      // Loop through scene objects in order to find intersections
      // ----------------------------------------------------------
      unsigned int numIntersections = 0;
      for (unsigned int j = 0; j < numObjects; j++) {
        long objType = objects[j].type;
        //  translate our ray into object space by multiplying ray pos and dir
        //  with inverse object matrix
        double4 tRayOrigin = mul(objects[j].inverse, rayOrigin);
        double4 tRayDirection = mul(objects[j].inverse, rayDirection);

        // Intersection code
        if (objType == 0) { // PLANE - intersect transformed ray with plane
          if (fabs(tRayDirection.y) > 0.0001) {
            double t = -tRayOrigin.y / tRayDirection.y;
            intersections[numIntersections] = t;
            xsObjects[numIntersections] = j;
            numIntersections++;
          }
        } else if (objType == 1) { // SPHERE
          // this is a vector from the origin of the ray to the center of the
          // sphere at 0,0,0
          double4 vecToCenter = tRayOrigin - originPoint;

          // This dot product is always 1.0 if tRayDirection is normalized. Which it isn't.
          double a = dot(tRayDirection, tRayDirection);

          // Take the dot of the direction and the vector from ray origin to
          // sphere center times 2
          double b = 2.0 * dot(tRayDirection, vecToCenter);

          // Take the dot of the two sphereToRay vectors and decrease by 1 (is
          // that because the sphere is unit length 1?
          double c = dot(vecToCenter, vecToCenter) - 1.0;

          // calculate the discriminant
          double discriminant = (b * b) - 4 * a * c;
          if (discriminant > 0.0) {
            // finally, find the intersection distances on our ray.
            double t1 = (-b - sqrt(discriminant)) / (2 * a);
            // double t2 = (-b + sqrt(discriminant)) / (2*a); // add back in
            // when we do refraction
            intersections[numIntersections] = t1;
            xsObjects[numIntersections] = j;
            numIntersections++;
          }
        } else if (objType == 2) { // CYLINDER
            // Cylinder intersection
             double rdx2 = tRayDirection.x * tRayDirection.x;
             double rdz2 = tRayDirection.z * tRayDirection.z;

             double a = rdx2 + rdz2;
             if (fabs(a) < 0.0001) {
                 //c.intercectCaps(ray, xs)
                 continue;
             }

             double b = 2*tRayOrigin.x*tRayDirection.x + 2*tRayOrigin.z*tRayDirection.z;

             double rox2 = tRayOrigin.x * tRayOrigin.x;
             double roz2 = tRayOrigin.z * tRayOrigin.z;

             double c1 = rox2 + roz2 - 1;

             double disc = b*b - 4*a*c1;

             // ray does not intersect the cylinder
             if (disc < 0.0) {
                 continue;
             }

             double t0 = (-b - sqrt(disc)) / (2 * a);
             double t1 = (-b + sqrt(disc)) / (2 * a);

             double y0 = tRayOrigin.y + t0*tRayDirection.y;

             if (y0 > objects[j].minY && y0 < objects[j].maxY) {
                 // add intersection
                 intersections[numIntersections] = t0;
                 xsObjects[numIntersections] = j;
                 numIntersections++;
             }

             double y1 = tRayOrigin.y + t1*tRayDirection.y;
             if (y1 > objects[j].minY && y1 < objects[j].maxY) {
                 // add intersection
                 intersections[numIntersections] = t1;
                 xsObjects[numIntersections] = j;
                 numIntersections++;
             }

             // TODO fix caps
             double2 caps = intersectCaps(tRayOrigin, tRayDirection, objects[j].minY, objects[j].maxY);
             if (caps.x > 0.0) {
                 intersections[numIntersections] = caps.x;
                 xsObjects[numIntersections] = j;
                 numIntersections++;
             }
             if (caps.y > 0.0) {
                 intersections[numIntersections] = caps.y;
                  xsObjects[numIntersections] = j;
                  numIntersections++;
             }
        } else if (objType == 3) { // BOX
            // There is supposed to be a way to optimize this for fewer checks by looking at early values.
            double2 xt = checkAxis(tRayOrigin.x, tRayDirection.x, -1.0, 1.0);
            double2 yt = checkAxis(tRayOrigin.y, tRayDirection.y, -1.0, 1.0);
            double2 zt = checkAxis(tRayOrigin.z, tRayDirection.z, -1.0, 1.0);

            // Om det största av min-värdena är större än det minsta max-värdet.
            double tmin = maxX(xt.x, yt.x, zt.x);
            double tmax = minX(xt.y, yt.y, zt.y);
            if (tmin > tmax) {
                // No intersection
                continue;
            }

            // assign intersections
            intersections[numIntersections] = tmin;
            xsObjects[numIntersections] = j;
            numIntersections++;

            intersections[numIntersections] = tmax;
            xsObjects[numIntersections] = j;
            numIntersections++;
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
                //skipped++;
                continue;
            }
            //hit++;

            // If the "object" BB was intersected, we take a look at the "object's" groupOffset. If > -1, we
            // need to set up a local stack to traverse the group hierarchy
            if (objects[j].groupOffset > -1) {

                // START PSUEDO-RECURSIVE CODE
                // 1) Create an empty stack. (move to top later)
                int stack[200] = {0}; // = make([]int, MAX_RECURSION_DEPTH)

                // Stack index, i.e. current "depth" of stack
                int currentSIndex = 0;

                // Tree index, i.e. which "node index" we're currently processing
                int currentNodeIndex = objects[j].groupOffset;

                // Initialize current node as root
                group root = groups[currentNodeIndex];
                group *current = &root;

                for (;current != 0 || currentSIndex > -1;) {
                    for (;current != 0 && intersectRayWithBox(tRayOrigin, tRayDirection, current->bbMin, current->bbMax);) { // HERE ADD BB CHECK //&& current.Inside) { && intersectRayWithBox(tRayOrigin, tRayDirection, current->bbMin, current->bbMax)

                    // Iterate over all triangles and record triangle/ray intersections...
                    for (int n = current->triOffset;n < current->triOffset+current->triCount;n++) {

                        double4 dirCrossE2 = cross(tRayDirection, triangles[n].e2);
                        double determinant = dot(triangles[n].e1, dirCrossE2);
                        if (fabs(determinant) < 0.0001) {
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
                        if (v < 0 || (u+v) > 1) {
                            continue;
                        }
                        double t = f * dot(triangles[n].e2, originCrossE1);
                        intersections[numIntersections] = t;
                        xsObjects[numIntersections] = j;

                        // assume we have vertex normals. If not, assume N in n1,n2,n3
                        // stored the computed normal in a list using the same indexing as xsObjects so
                        // if a ray intersects several triangles in the group, we'll get an intersection per triangle
                        // but can separate their normals and then only use the one for the nearest intersection
                        xsTriangle[numIntersections] = triangles[n].n2 * u +
                                        triangles[n].n3 * v +
                                        triangles[n].n1  * (1.0-u-v);

                        numIntersections++;
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
                    }

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
                        current = &root; //groups[current->children[1]];
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

      // find lowest positive intersection index
      double lowestIntersectionT = 1024.0;
      int lowestIntersectionIndex = -1;
      int normalIndex = -1;
      for (unsigned int x = 0; x < numIntersections; x++) {
        if (intersections[x] > 0.0001) {
          if (intersections[x] < lowestIntersectionT) {
            lowestIntersectionT = intersections[x];
            lowestIntersectionIndex = xsObjects[x];
            normalIndex = x;                          // while only used for triangles, we track computed normal by x.
          }
        }
      }

      if (lowestIntersectionIndex > -1) {
        object obj = objects[lowestIntersectionIndex];
        // Remember that we use the untransformed ray here!

        // Position gives us the intersection position along RAY at T
        double4 position = rayOrigin + rayDirection * lowestIntersectionT;

        // The vector to the eye (or last bounce origin) is exactly the opposite
        // of the ray direction
        double4 eyeVector = -rayDirection;

        // object normal at intersection: Transform point from world to object
        // space

        double4 objectNormal;

        // PLANE always have its normal UP in local space
        if (obj.type == 0) {
          objectNormal = (double4)(0.0, 1.0, 0.0, 0.0);
        } else if (obj.type == 1) {
          // SPHERE always has its normal from sphere center outwards to the
          // world position.
           double4 localPoint = mul(obj.inverse, position);
          objectNormal = localPoint - originPoint;
        } else if (obj.type == 2) {
            // CYLINDER
            // compute the square of the distance from the y axis
             double4 localPoint = mul(obj.inverse, position);
            double dist = pow(localPoint.x, 2) + pow(localPoint.z, 2);
            if (dist < 1 && localPoint.y >= obj.maxY - 0.0001) {
                objectNormal = (double4)(0.0, 1.0, 0.0, 0.0);
            } else if (dist < 1 && localPoint.y <= obj.minY + 0.0001) {
                objectNormal = (double4)(0.0, -1.0, 0.0, 0.0);
            } else {
                objectNormal = (double4)(localPoint.x, 0.0, localPoint.z, 0.0);
            }
        } else if (obj.type == 3) {
            // CUBE
            // NormalAtLocal for a cube uses the fact that given a unit cube, the point of the surface axis X,Y or Z is
            // always either 1.0 for positive XYZ and -1.0 for negative XYZ. I.e - if the point is 0.4, 1, -0.5,
            // we know that the point is on the top Y surface and we can return a 0,1,0 normal.
             double4 localPoint = mul(obj.inverse, position);
            double maxc = maxX(fabs(localPoint.x), fabs(localPoint.y), fabs(localPoint.z));
            if (maxc == fabs(localPoint.x)) {
                objectNormal = (double4)(localPoint.x, 0.0, 0.0, 0.0);
            } else if (maxc == fabs(localPoint.y)) {
                objectNormal = (double4) (0.0, localPoint.y, 0.0, 0.0);
            } else {
                objectNormal = (double4) (0.0, 0.0, localPoint.z, 0.0);
            }
        } else if (obj.type == 4) {
            // GROUP, which in practice means a triangle, whose normal is typically pre-populated in N and stored in xsTriangles
            objectNormal = xsTriangle[normalIndex];
            //printf("object normal: %f %f %f\n", objectNormal.x, objectNormal.y, objectNormal.z);
        }
        // Finish the normal vector by multiplying it back into world coord
        // using the inverse transpose matrix and then normalize it
        double4 normalVec = mul(obj.inverseTranspose, objectNormal);
        normalVec.w = 0.0; // set w to 0
        normalVec = normalize(normalVec);

        // The "inside" stuff from the old impl will be needed for refraction
        // later comps.Inside = false

        // negate the normal if the normal if facing
        // away from the "eye"
        if (dot(eyeVector, normalVec) < 0.0) {
          normalVec = normalVec * -1.0;
        }

        // Compute the over point, with a slight offset along the normal, in
        // order to avoid self-intersection on the next bounce.
        double4 overPoint = position + normalVec * 0.0001;

        // Prepare the outgoing ray (next bounce) by reusing the original ray, just
        // update its origin and direction.

        // Impl here supports either diffuse or reflected, but for obj.reflectivity > 0 a proportionate portion of samples
        // will diffuse instead of reflect. Poor-man's BRDF
        if (obj.reflectivity == 0.0 || noise3D(fgi, n, b) > obj.reflectivity) {
            // Diffuse
            rayDirection = randomVectorInHemisphere(normalVec, fgi, b, n);
        } else {
            // Reflected, calculate reflection vector and set as rayDirection
            double dotScalar = dot(rayDirection, normalVec);
            double4 norm = (normalVec * 2.0) * dotScalar;
            rayDirection = rayDirection - norm;
        }

        rayOrigin = overPoint;

        // Calculate the cosine of the OUTGOING ray in relation to the surface
        // normal.
        double cosine = dot(rayDirection, normalVec);

        // Finish this iteration by storing the bounce.
        bounce bnce = {position, cosine, obj.color, obj.emission};
        bounces[b] = bnce;
        actualBounces++;
      }
    }

    // ------------------------------------ //
    // Calculate final color using bounces! //
    // ------------------------------------ //
    double4 accumColor = (double4)(0.0, 0.0, 0.0, 0.0);
    double4 mask = (double4)(1.0, 1.0, 1.0, 1.0);
    for (unsigned int x = 0; x < actualBounces; x++) {

      // Start by dealing with diffuse surfaces.
      // First, ADD current accumulated color with the hadamard of the current
      // mask and the emission properties of the hit object.
      accumColor = accumColor + mask * bounces[x].emission;

      // Update the mask by multiplying it with the hit object's color
      mask *= bounces[x].color;

      // perform cosine-weighted importance sampling by multiplying the mask
      // with the cosine
      mask *= bounces[x].cos;
    }

    // Finish this "sample" by adding the accumulated color to the total
    colors += accumColor;
  }

  // Finish the pixel by multiplying each RGB component by its total fraction and
  // store in the output bufer.
  output[i * 4] = colors.x * colorWeight;
  output[i * 4 + 1] = colors.y * colorWeight;
  output[i * 4 + 2] = colors.z * colorWeight;
  output[i * 4 + 3] = 1.0;

  //printf("%d rays missing the BB, %d hits\n", skipped, hit);
}