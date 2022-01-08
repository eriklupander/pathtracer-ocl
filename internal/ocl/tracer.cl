__constant double PI = 3.14159265359f;
__constant unsigned int MAX_BOUNCES = 4;

typedef struct tag_ray {
  double4 origin;
  double4 direction;
} ray;

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
  double padding1;           // 8 bytes
  double padding2;           // 8 bytes
  double padding3;           // 8 bytes
  double padding4;           // 8 bytes
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

inline bool checkCap(double4 origin, double4 direction, double t) {
	double x = origin.x + t*direction.x;
	double z = origin.z + t*direction.z;
	return pow(x, 2) + pow(z, 2) <= 1.0;
}

inline double2 intersectCaps(double4 origin, double4 direction, object* obj) {
	// !c.closed removed
    if (fabs(direction.y) < 0.0001) {
		return (double2)(0.0,0.0);
	}

    double2 retVal = (double2)(0.0,0.0);

	// check for an intersection with the lower end cap by intersecting
	// the ray with the plane at y=cyl.minimum
	double t1 = (obj->minY - origin.y) / direction.y;
	if (checkCap(origin, direction, t1)) {
		//*xs = append(*xs, NewIntersection(t, c))
        // TODO 
        retVal.x = t1;
	}

	// check for an intersection with the upper end cap by intersecting
	// the ray with the plane at y=cyl.maximum
	double t2 = (obj->maxY - origin.y) / direction.y;
	if (checkCap(origin, direction, t2)) {
		//*xs = append(*xs, NewIntersection(t, c))
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

__kernel void trace(__global ray *rays, __global object *objects,
                    const unsigned int numObjects, __global double *output,
                    __global double *seedX, const unsigned int samples) {
  double colorWeight = 1.0 / samples;
  int i = get_global_id(0);

  float fgi = seedX[i] / numObjects;

  double4 originPoint = (double4)(0.0f, 0.0f, 0.0f, 1.0f);

  double4 colors = (double4)(0, 0, 0, 0);

  for (unsigned int n = 0; n < samples; n++) {
    // Each new sample needs to reset to original ray
    double4 rayOrigin = rays[i].origin;
    double4 rayDirection = rays[i].direction;

    // for each bounce...
    unsigned int actualBounces = 0;
    // Each ray may bounce up to 5 times
    bounce bounces[5] = {};
    for (unsigned int b = 0; b < MAX_BOUNCES; b++) {
     // track up to 16 intersections per ray.
      double intersections[16] = {0};  // t of an intersection
      //unsigned int xsObjects[8] = {0}; // index maps to each xs above, value to objects


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
        if (objType == 0) { // intersect transformed ray with plane
          if (fabs(tRayDirection.y) > 0.0001) {
            double t = -tRayOrigin.y / tRayDirection.y;
            intersections[j] = t;
            //xsObjects[numIntersections] = j;
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
            intersections[j] = t1;
            //xsObjects[numIntersections] = j;
            numIntersections++;
          }
        } else if (objType == 2) {
            // Cylinder intersection
            // double rdx2 = tRayDirection.x * tRayDirection.x;
            // double rdz2 = tRayDirection.z * tRayDirection.z;

            // double a = rdx2 + rdz2;
            // if (fabs(a) < 0.0001) {
            //     //c.intercectCaps(ray, xs)
            //     continue;
            // }

            // double b = 2*tRayOrigin.x*tRayDirection.x + 2*tRayOrigin.z*tRayDirection.z;

            // double rox2 = tRayOrigin.x * tRayOrigin.x;
            // double roz2 = tRayOrigin.z * tRayOrigin.z;

            // double c1 = rox2 + roz2 - 1;

            // double disc = b*b - 4*a*c1;

            // // ray does not intersect the cylinder
            // if (disc < 0.0) {
            //     continue;
            // }

            // double t0 = (-b - sqrt(disc)) / (2 * a);
            // double t1 = (-b + sqrt(disc)) / (2 * a);

            // double y0 = tRayOrigin.y + t0*tRayDirection.y;

            // // BROKEN BELOW!!!
            // if (y0 > objects[j].minY && y0 < objects[j].maxY) {
            //     //*xs = append(*xs, NewIntersection(t0, c))
            //     // add intersection
            //     intersections[j] = t0;
            //     numIntersections++;
            // }

            // double y1 = tRayOrigin.y + t1*tRayDirection.y;
            // if (y1 > objects[j].minY && y1 < objects[j].maxY) {
            //     //*xs = append(*xs, NewIntersection(t1, c))
            //     // add intersection
            //     intersections[j] = t1;
            //     numIntersections++;
            // }

            // // TODO fix caps
            // double2 caps = intersectCaps(tRayOrigin, tRayDirection, xs, &objects[j])
            // if (caps.x > 0.0) {
            //     intersections[j] = caps.x;
            //     numIntersections++;
            // }
            // if (caps.y > 0.0) {
            //     intersections[j] = caps.y;
            //     numIntersections++;
            // }
        }
      }

      // find lowest positive intersection index
      double lowestIntersectionT = 1024.0;
      int lowestIntersectionIndex = -1;
      for (unsigned int x = 0; x < 16; x++) {
        if (intersections[x] > 0.0001) {
          if (intersections[x] < lowestIntersectionT) {
            lowestIntersectionT = intersections[x];
            lowestIntersectionIndex = x;
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
        double4 localPoint = mul(obj.inverse, position);
        double4 objectNormal;

        // PLANE always have its normal UP in local space
        if (obj.type == 0) {
          objectNormal = (double4)(0.0, 1.0, 0.0, 0.0);
        } else if (obj.type == 1) {
          // SPHERE always has its normal from sphere center outwards to the
          // world position.
          objectNormal = localPoint - originPoint;
        } else if (obj.type == 2) {
            // CYLINDER
            // compute the square of the distance from the y axis
            double dist = pow(localPoint.x, 2) + pow(localPoint.z, 2);
            if (dist < 1 && localPoint.y >= obj.maxY - 0.0001) {
                objectNormal = (double4)(0.0, 1.0, 0.0, 0.0);
            } else if (dist < 1 && localPoint.y <= obj.minY + 0.0001) {
                objectNormal = (double4)(0.0, -1.0, 0.0, 0.0);
            } else {
                objectNormal = (double4)(localPoint.x, 0.0, localPoint.z, 0.0);
            }
        }
        // Finish the normal vector by multiplying it back into world coord
        // using the inverse transpose matrix and then normalize it
        double4 normalVec = mul(obj.inverseTranspose, objectNormal);
        normalVec.w = 0.0; // set w to 0
        normalVec = normalize(normalVec);

        // reflection vector  (add when we start doing reflective materials)
        // double dotScalar = dot(rayDirection, normalVec);
        // double4 norm = (normalVec * 2.0) * dotScalar;
        // double4 reflectVec = rayDirection - norm;

        // The "inside" stuff from the old impl will be needed for refraction
        // later comps.Inside = false negate the normal if the normal if facing
        // away from the "eye"
        if (dot(eyeVector, normalVec) < 0.0) {
          normalVec = normalVec * -1.0;
        }

        // Compute the over point, with a slight offset along the normal, in
        // order to avoid self-intersection on the next bounce.
        double4 overPoint = position + normalVec * 0.0001;

        // Prepare the outgoing ray (next bounce), reuse the original ray, just
        // update its origin and direction.
        rayDirection = randomVectorInHemisphere(normalVec, fgi, b, n);
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

/*
      if (i == 6397 && bounces[actualBounces-1].emission.x > 0.0) {
        printf("bounce: %d ", x);
        printf("accum: %f %f %f ", accumColor.x, accumColor.y, accumColor.z);
        printf("mask: %f %f %f ", mask.x, mask.y, mask.z);
        printf("cos: %f ", bounces[x].cos);
        printf("color: %f %f %f ", bounces[x].color.x, bounces[x].color.y, bounces[x].color.z);
        printf("emission: %f %f %f\n", bounces[x].emission.x, bounces[x].emission.y, bounces[x].emission.z);
      }
      */
    }

    // Finish this "sample" by adding the accumulated color to the total
    colors += accumColor;
  }

  // Finish the ray by multiplying each RGB component by its total fraction and
  // store in the output bufer.
  output[i * 4] = colors.x * colorWeight;
  output[i * 4 + 1] = colors.y * colorWeight;
  output[i * 4 + 2] = colors.z * colorWeight;
  output[i * 4 + 3] = 1.0;
}