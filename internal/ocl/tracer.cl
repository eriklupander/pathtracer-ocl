__constant double PI = 3.14159265359f;
__constant double PI_X2 = 6.28318530718f;
__constant unsigned int MAX_EFFECTIVE_BOUNCES = 4;
__constant unsigned int MAX_BOUNCES = 10;
__constant double EPSILON = 0.0001;

typedef struct __attribute__((packed)) tag_camera {
    int width;          // 4 bytes
    int height;         // 4 bytes
    double fov;         // 8 bytes
    double pixelSize;   // 8 bytes
    double halfWidth;   // 8 bytes
    double halfHeight;  // 8 bytes
    double aperture;    // 8 bytes
    double focalLength; // 8 bytes
    double16 inverse;   // 128 bytes
    char padding[72];   // 72 bytes
} camera;

typedef struct tag_ray {
    double4 origin;
    double4 direction;
} ray;

typedef struct __attribute__((packed)) tag_group {
    double4 bbMin;       // 32 bytes
    double4 bbMax;       // 32 bytes
    double4 color;       // 32 bytes
    double4 emission;    // 32 bytes
    int triOffset;       // 4 bytes
    int triCount;        // 4 bytes
    int childGroupCount; // 4 bytes, should always be 2 or 0
    int children[2];     // 8 bytes, we only allow binary trees.
    char padding[108];   // padding, 108 bytes (can be used as a label)
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
    double textureScaleX;
    double textureScaleY;
    double textureScaleXNM;
    double textureScaleYNM;
    double4 bbMin;             // 32 bytes
    double4 bbMax;             // 32 bytes                     // 576
    int childCount;            // 4 bytes. Used for groups to know which "group" that's the root group.
    int children[64];          // 256 bytes
    bool isTextured;           // 1 byte
    unsigned char textureIndex;// 1 byte
    bool isTexturedNM;           // 1 byte
    unsigned char textureIndexNM;// 1 byte
    bool isRefraction;               // 1 byte
    char label[8];               // 8 bytes
    char padding5[167];          // ==> 1024
} object;

typedef struct tag_intersection_old {
    unsigned int objectIndex;
    double t;
    double4 color;      // while color and emission can be read from the "object" referenced by objectIndex,
    double4 emission;   // 3D models organized into BVH trees needs to get their material from the intersected group of the tree.
} intersection_old;

typedef struct tag_bounce {
    double4 point;
    double cos;
    double4 color;
    double4 emission;
    double4 normal;
    double refractiveIndex;
    bool isRefraction;
} bounce;

typedef struct __attribute__((packed)) tag_triangle {
    double4 p1;           // 32 bytes
    double4 p2;           // 32 bytes
    double4 p3;           // 32 bytes
    double4 e1;           // 32 bytes
    double4 e2;           // 32 bytes
    double4 n1;           // 32 bytes
    double4 n2;           // 32 bytes
    double4 n3;           // 32 bytes
    double4 color;        // 32 bytes (288 bytes)
    char	padding[224]; // 224 bytes
} triangle;               // 512 total

// used as an internal data structure
typedef struct tag_context {
    double intersections[64]; // = {0};   // t of an intersection (MOVE TO LOCAL)
    unsigned int xsObjects[64]; // = {0}; // index maps to each xs above, value to objects
    double4 xsTriangle[64]; // = {0};
    double4 xsTriangleColor[64]; // = {0};
    double4 xsTriangleEmission[64]; // = {0};
} context;

typedef struct intersection_tag {
    double t;
    int lowestIntersectionIndex;
    int normalIndex;
} intersection;

inline double maxX(double a, double b, double c) { return max(max(a, b), c); }
inline double minX(double a, double b, double c) { return min(min(a, b), c); }

inline double2 cubeUVFrontCross(double4 point) {
	double u = fmod(point.x+1.0, 2) / 2.0;
	double v = fmod(point.y+1.0, 2) / 2.0;
	double2 uv = (double2)(0.25 + u*0.25, 0.6666666-v*0.333333);
	return uv;
}
inline double2 cubeUVBackCross(double4 point) {
	double u = fmod(1.0-point.x, 2) / 2.0;
	double v = fmod(point.y+1.0, 2) / 2.0;
	double2 uv = (double2)(0.75+u*0.25, 0.6666666-v*0.333333);
	return uv;
}
inline double2 cubeUVLeftCross(double4 point) {
	double u = fmod(point.z+1.0, 2) / 2.0;
	double v = fmod(point.y+1.0, 2) / 2.0;
	double2 uv = (double2)(u*0.25, 0.6666666-v*0.333333);
	return uv;
}
inline double2 cubeUVRightCross(double4 point) {
	double u = fmod(1.0-point.z, 2) / 2.0;
	double v = fmod(point.y+1.0, 2) / 2.0;
	double2 uv = (double2)(0.5+u*0.25, 0.6666666-v*0.333333);
	return uv;
}
inline double2 cubeUVTopCross(double4 point) {
	double u = fmod(point.x+1.0, 2) / 2.0;
	double v = fmod(1.0-point.z, 2) / 2.0;
	double2 uv = (double2)(0.25+u*0.25, 1.0-v*0.333333);
	return uv;
}
inline double2 cubeUVBottomCross(double4 point) {
	double u = fmod(point.x+1.0, 2) / 2.0;
	double v = fmod(point.z+1.0, 2) / 2.0;
	double2 uv = (double2)(0.25+u*0.25, v*0.333333);
	return uv;
}


inline double2 cubeUV(double4 point) {

	double absX = fabs(point[0]);
	double absY = fabs(point[1]);
	double absZ = fabs(point[2]);
	double coord = maxX(absX, absY, absZ);

	if (coord == point[0]) {
		return cubeUVRightCross(point); // right
	}
	if (coord == -point[0]) {
	     return cubeUVLeftCross(point); //"left"
	}
	if (coord == point[1]) {
		return cubeUVTopCross(point); //"up"
	}
	if (coord == -point[1]) {
		return cubeUVBottomCross(point); //"down"
	}
	if (coord == point[2]) {
	  	return cubeUVFrontCross(point); // "front"
	}

	return cubeUVBackCross(point); //"back"
}


inline double2 sphericalMap(double4 p) {

	// compute the azimuthal angle
	// -π < theta <= π
	// angle increases clockwise as viewed from above,
	// which is opposite of what we want, but we'll fix it later.
	double theta = atan2(p.x, p.z);

	// vec is the vector pointing from the sphere's origin (the world origin)
	// to the point, which will also happen to be exactly equal to the sphere's
	// radius.
	double4 vec = (double4)(p.x, p.y, p.z, 0.0);
	double radius = length(vec);

	// compute the polar angle
	// 0 <= phi <= π
	double phi = acos(p.y / radius);

	// -0.5 < raw_u <= 0.5
	double rawU = theta / (2.0 * PI);

	// 0 <= u < 1
	// here's also where we fix the direction of u. Subtract it from 1,
	// so that it increases counterclockwise as viewed from above.
	double u = 1 - (rawU + 0.5);

	// we want v to be 0 at the south pole of the sphere,
	// and 1 at the north pole, so we have to "flip it over"
	// by subtracting it from 1.
	double v = 1 - phi/PI;

    double2 res;
    res.x = u;
    res.y = v;
	return res;
}

inline int round2(double number) {
   int sign = (int)((number > 0) - (number < 0));
   int odd = ((int)number % 2); // odd -> 1, even -> 0
   return ((int)(number-sign*(0.5-odd)));
}

inline double sunflowerRadius(double i, double n, double b) {
  double r = 1.0; // put on boundary
   if (i <= (n - b)) {
      r = sqrt(i-0.5) / sqrt(n-(b+1.0)/2.0); // apply square root
   }
   return r;
}

// Distributes n points evenly within a circle with sunflowerRadius 1
// alpha controls point distribution on the edge. Typical values 1-2, higher values more points on the edge.
// i is the index of a point. It is in the range [1,n] .
// https://stackoverflow.com/questions/28567166/uniformly-distribute-x-points-inside-a-circle
//
// example: amountPoints=500, alpha=2, pointNumber=[1..amountPoints]
inline double2 sunflower(int amountPoints, double alpha, int pointNumber, bool randomize, double rand) {
   double pointIndex = (double) pointNumber;  //float64(pointNumber)
   if (randomize) {
      pointIndex += rand - 0.5;
   }

   double sqp = sqrt(convert_double(amountPoints));
   double b = round(alpha * sqp); // number of boundary points
   double phi = (sqrt(5.0) + 1.0) / 2.0;                                // golden ratio
   double r = sunflowerRadius(pointIndex, amountPoints, b);
   double theta = 2.0 * PI * pointIndex / (phi * phi);

   return (double2)(r * cos(theta), r * sin(theta));
}

inline double2 checkAxis(double origin, double direction, double minBB, double maxBB) {
    double2 out = (double2){0, 0};
    double tminNumerator = minBB - origin; //-1.0 - origin;
    double tmaxNumerator = maxBB - origin; // 1.0 - origin;
    if (fabs(direction) >= EPSILON) {
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

    // If the largest of the min values is greater smallest max value...
    double tmin = maxX(xt.x, yt.x, zt.x); // x == min
    double tmax = minX(xt.y, yt.y, zt.y); // y == max
    return tmin < tmax;
}

inline bool checkCap(double4 origin, double4 direction, double t) {
    double x = origin.x + t * direction.x;
    double z = origin.z + t * direction.z;
    return pow(x, 2) + pow(z, 2) <= 1.0;
}

inline double2 intersectCaps(double4 origin, double4 direction, double minY, double maxY) {
    // !c.closed removed
    if (fabs(direction.y) < EPSILON) {
        return (double2)(0.0, 0.0);
    }

    double2 retVal = (double2)(0.0, 0.0);

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

// from https://stackoverflow.com/a/50665114. Perhaps try to find a better way to create random numbers? Maybe
// it's possible to preload all random numbers and/or cycle a finite set?
inline static float noise3D(float x, float y, float z) {
    float ptr = 0.0f;
    return fract(sin(x * 112.9898f + y * 179.233f + z * 237.212f) * 43758.5453f, &ptr);
}

// from https://math.stackexchange.com/questions/1585975/how-to-generate-random-points-on-a-sphere
// note that we're exchanging y and z since y is up for us, while the formula above uses z as up.
inline double4 randomPointOnSphere(double r, double u1, double u2) {
    //latitude: 𝜆=arccos(2𝑢1−1)−𝜋2 OR arcsin(2𝑎−1)
    //longitude:𝜙=2𝜋𝑢2
    double lat = acos(2*u1 - 1) - PI*2; // asin(2*u1-1);
    double lon = 2*PI*u2;

    // 𝑥=cos𝜆cos𝜙
    // 𝑦=cos𝜆sin𝜙
    // 𝑧=sin𝜆
    double4 out = (double4)(0.0, 0.0, 0.0, 1.0);
    out.x = cos(lat) * cos(lon) * r;
    out.y = (sin(lat) - PI*0.25)  * r;
    out.z = cos(lat) * sin(lon) * r;

    return out;
}

inline double4 randomSphereDirection(double x, double y, double z) {
    double rnd1 = noise3D(y, z, x);
    double rnd2 = noise3D(x, y, z);
    double2 h = double2(rnd1, rnd2) * double2(2.0, PI_X2) - double2(1.0, 0.0);
    float phi = h.y;

    // vec3(sqrt(1.-h.x*h.x)*vec2(sin(phi),cos(phi)),h.x);
    double2 tmp = sqrt(1.0 - h.x * h.x) * double2(sin(phi), cos(phi));
	return double4(tmp.x, tmp.y, h.x, 0.0);
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
inline double4 randomVectorInHemisphere(double4 normalVec, double x, double y, double z) {
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

    // use the coordinate frame and random numbers to compute the next ray direction
    return u * cos(rand1) * rand2s + v * sin(rand1) * rand2s + normalVec * sqrt(1.0 - rand2);
}

// mul multiplies the vec by the matrix, producing a new vector.
inline double4 mul(double16 mat, double4 vec) {
    double4 elem1 = mat.s0123 * vec;
    double4 elem2 = mat.s4567 * vec;
    double4 elem3 = mat.s89AB * vec;
    double4 elem4 = mat.sCDEF * vec;
    return (double4)(elem1.x + elem1.y + elem1.z + elem1.w, elem2.x + elem2.y + elem2.z + elem2.w, elem3.x + elem3.y + elem3.z + elem3.w,
                     elem4.x + elem4.y + elem4.z + elem4.w);
}

inline double2 intersectCube(double4 tRayOrigin, double4 tRayDirection) {
    double2 out = (0,0);
    // There is supposed to be a way to optimize this for fewer checks by looking at early values.
    double2 xt = checkAxis(tRayOrigin.x, tRayDirection.x, -1.0, 1.0);
    double2 yt = checkAxis(tRayOrigin.y, tRayDirection.y, -1.0, 1.0);
    double2 zt = checkAxis(tRayOrigin.z, tRayDirection.z, -1.0, 1.0);

    // Om det största av min-värdena är större än det minsta max-värdet.
    double tmin = maxX(xt.x, yt.x, zt.x);
    double tmax = minX(xt.y, yt.y, zt.y);
    if (tmin > tmax) {
        return out;
    }
    out.x = tmin;
    out.y = tmax;
    return out;
}

inline double4 intersectCylinder(double4 tRayOrigin, double4 tRayDirection, object obj) {
    double4 out={0,0,0,0};
    double rdx2 = tRayDirection.x * tRayDirection.x;
    double rdz2 = tRayDirection.z * tRayDirection.z;

    double a = rdx2 + rdz2;
    if (fabs(a) < EPSILON) {
        // c.intercectCaps(ray, xs)
        return out;
    }

    double b = 2 * tRayOrigin.x * tRayDirection.x + 2 * tRayOrigin.z * tRayDirection.z;

    double rox2 = tRayOrigin.x * tRayOrigin.x;
    double roz2 = tRayOrigin.z * tRayOrigin.z;

    double c1 = rox2 + roz2 - 1;

    double disc = b * b - 4 * a * c1;

    // ray does not intersect the cylinder
    if (disc < 0.0) {
        return out;
    }

    double t0 = (-b - sqrt(disc)) / (2 * a);
    double t1 = (-b + sqrt(disc)) / (2 * a);

    double y0 = tRayOrigin.y + t0 * tRayDirection.y;

    if (y0 > obj.minY && y0 < obj.maxY) {
        // add intersection
        out.x = t0;
    }

    double y1 = tRayOrigin.y + t1 * tRayDirection.y;
    if (y1 > obj.minY && y1 < obj.maxY) {
        // add intersection
        out.y = t1;
    }

    // TODO fix so caps can be enabled/disabled... for now, disable.
//    double2 caps = intersectCaps(tRayOrigin, tRayDirection, obj.minY, obj.maxY);
//    if (caps.x > 0.0) {
//        out.z = caps.x;
//    }
//    if (caps.y > 0.0) {
//        out.w = caps.y;
//    }
    return out;
}

inline double2 intersectSphere(double4 tRayOrigin, double4 tRayDirection) {
    // this is a vector from the origin of the ray to the center of the
    // sphere at 0,0,0
    double4 vecToCenter = tRayOrigin - ((double4)(0.0, 0.0, 0.0, 1.0));

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
        double t2 = (-b + sqrt(discriminant)) / (2*a);
        double2 t;
        t.x = t1;
        t.y = t2;
        return t;
    }
    return (double2)(0.0, 0.0);
}

inline double intersectPlane(double4 tRayOrigin, double4 tRayDirection) {
    if (fabs(tRayDirection.y) > EPSILON) {
            return -tRayOrigin.y / tRayDirection.y;
    }
    return 0.0;
}

inline double schlick(double4 eyeVec, double4 normalVec, double n1, double n2) {

    // find the cosine of the angle between the eye and normal vectors using Dot
    double cos = dot(eyeVec, normalVec);
    // total internal reflection can only occur if n1 > n2
    if (n1 > n2) {
        double n = n1 / n2;
        double sin2Theta = (n * n) * (1.0 - (cos * cos));
        if (sin2Theta > 1.0) {
            return 1.0;
        }
        // compute cosine of theta_t using trig identity
        double cosTheta = sqrt(1.0 - sin2Theta);

        // when n1 > n2, use cos(theta_t) instead
        cos = cosTheta;
    }
    double temp = (n1 - n2) / (n1 + n2);
    double r0 = temp * temp;
    return r0 + (1-r0)*pow(1-cos, 5);
}

inline double4 computeRefractedRay(double4 eyeVector, double4 normalVec, double n1, double n2) {
    // Find the ratio of first index of refraction to the second.
	double nRatio = n1 / n2;

	// cos(theta_i) is the same as the dot product of the two vectors
	double cosI = dot(eyeVector, normalVec);

	// Find sin(theta_t)^2 via trigonometric identity
	double sin2Theta = (nRatio * nRatio) * (1.0 - (cosI * cosI));
	if (sin2Theta > 1.0) {
	    // was black, how to handle?? This is probably that famous total reflectance?
	    // In the original ray-tracer, this meant that the refraction did not contribute any "color" to
	    // the final pixel color.
		return (double4)(0,0,0,0);
	}

	// Find cos(theta_t) via trigonometric identity
	double cosTheta = sqrt(1.0 - sin2Theta);

	// Compute the direction of the refracted ray
	double4 direction = (normalVec * ((nRatio*cosI)-cosTheta)) - eyeVector * nRatio;

    // Return the refracted ray direction vector (use underpoint at callsite)
    return direction;

	//refractRay := mat.NewRay(comps.UnderPoint, direction)
}

// findClosestIntersection returns the closest intersection. NOTE! It possible we could optimize this for shadow rays,
// if we pass some kind of maxT - if
inline intersection findClosestIntersection(__local object *objects, unsigned int numObjects, __global group *groups, __global triangle *triangles, double4 rayOrigin, double4 rayDirection, context *ctx) {
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
            double t = intersectPlane(tRayOrigin, tRayDirection);
            if (t != 0.0) {
                ctx->intersections[numIntersections] = t;
                ctx->xsObjects[numIntersections] = j;
                numIntersections++;
            }
        } else if (objType == 1) { // SPHERE

            // finally, find the intersection distances on our ray.
            double2 t = intersectSphere(tRayOrigin, tRayDirection);
             // required for refraction and possibly to detect when the camera starts inside a sphere

            if (t.x != 0.0) {
                ctx->intersections[numIntersections] = t.x;
                ctx->xsObjects[numIntersections] = j;
                numIntersections++;
            }
            if (t.y != 0.0) {
                ctx->intersections[numIntersections] = t.y;
                ctx->xsObjects[numIntersections] = j;
                numIntersections++;
            }
        } else if (objType == 2) { // CYLINDER
            double4 out = intersectCylinder(tRayOrigin, tRayDirection, objects[j]);
            for (unsigned int a = 0; a < 4; a++) {
                if (out[a] != 0) {
                    ctx->intersections[numIntersections] = out[a];
                    ctx->xsObjects[numIntersections] = j;
                    numIntersections++;
                }
            }
        } else if (objType == 3) { // BOX
            double2 out = intersectCube(tRayOrigin, tRayDirection);

            // assign intersections
            if (out.x != 0.0) {
                ctx->intersections[numIntersections] = out.x;
                ctx->xsObjects[numIntersections] = j;
                numIntersections++;
            }
            if (out.y != 0.0) {
                ctx->intersections[numIntersections] = out.y;
                ctx->xsObjects[numIntersections] = j;
                numIntersections++;
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
                                ctx->intersections[numIntersections] = t;
                                ctx->xsObjects[numIntersections] = j;

                                // assume we have vertex normals. If not, assume N in n1,n2,n3
                                // stored the computed normal in a list using the same indexing as xsObjects so
                                // if a ray intersects several triangles in the group, we'll get an intersection per triangle
                                // but can separate their normals and then only use the one for the nearest intersection
                                ctx->xsTriangle[numIntersections] = triangles[n].n2 * u + triangles[n].n3 * v + triangles[n].n1 * (1.0 - u - v);

                                // experiment: record the color and emission of the intersection
                                ctx->xsTriangleColor[numIntersections] = triangles[n].color;
                                ctx->xsTriangleEmission[numIntersections] = (double4){0,0,0,0}; //triangles[n].emission;
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

    if  (numIntersections == 0) {
        return (intersection){0.0, -1, -1};
    }

    // find lowest positive intersection index
    double lowestIntersectionT = 1024.0;
    int lowestIntersectionIndex = -1;
    int normalIndex = -1;
    for (unsigned int x = 0; x < numIntersections; x++) {
        if (ctx->intersections[x] > EPSILON) {
            if (ctx->intersections[x] < lowestIntersectionT) {
                lowestIntersectionT = ctx->intersections[x];
                lowestIntersectionIndex = ctx->xsObjects[x];
                normalIndex = x; // while only used for triangles, we track computed normal by x.
            }
        }
    }
    intersection ixs = {lowestIntersectionT, lowestIntersectionIndex, normalIndex};
    return ixs;
}


inline ray rayForPixel(unsigned int x, unsigned int y, camera cam, float rndX, float rndY, int sample, int totalSamples) {
    double4 pointInView = {0.0, 0.0, -1.0, 1.0};
    double4 originPoint = {0.0, 0.0, 0.0, 1.0};
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

        double4 pos = origin + direction * cam.focalLength;
        double4 newOrigin = {};
        double2 xy = sunflower(totalSamples, 2, sample, false, rndX);
        //printf("X: %f Y: %f ", xy.x, xy.y);
        newOrigin.x = origin.x + (xy.y * cam.aperture);
        newOrigin.y = origin.y + (xy.x * cam.aperture);
//        newOrigin.x = origin.x + (-cam.aperture + xy.y * cam.aperture * 2);
//        newOrigin.y = origin.y + (-cam.aperture + xy.x * cam.aperture * 2);
        newOrigin.z = origin.z;
        newOrigin.w = 1.0;
        direction = pos - newOrigin;
        origin = newOrigin;
    }
    ray r = {origin, direction};
    return r;
}

// nextEventEstimation is a very efficient method of reducing noise by directly sampling all light sources for each bounce,
// checking for line of sight to a random point on every lightsource. However, NEE only works reasonably well with diffuse
// materials.
//
// This function operates on a
inline void nextEventEstimation(__local object *objects, unsigned int numObjects, __global group *groups, __global triangle *triangles, bounce *b, double fgi, double fgi2, double n, double4 mask, unsigned int x, double4 *accumColor) {
    for (unsigned int l = 0; l < numObjects;l++) {
        if (objects[l].emission.x > 0.0) { // Note: handle if we have a light source without red emission...

            double4 lightOriginPosition = (double4)(objects[l].transform[3], objects[l].transform[7], objects[l].transform[11], 0.0); // note .w will be == 1 after next line
            double scaleBy = max(max(objects[l].transform[0], objects[l].transform[5]), objects[l].transform[10]);
            double4 lightScale = (double4)(scaleBy, scaleBy, scaleBy, 1.0);
            double4 rpos = randomPointOnSphere(1.0, noise3D(fgi, n+x*l, fgi2), noise3D(fgi, fgi2, n+x*x*l));
            double4 lightPosition = lightOriginPosition + (rpos * lightScale);

            double4 shadowRayDirection = normalize(lightPosition - b->point);
            double4 shadowRayOrigin = b->point + (shadowRayDirection*EPSILON); // take a slight overpos

            double lightDotNormal = dot(shadowRayDirection, b->normal);
            if (lightDotNormal > 0.0) {

                // now, we need to check if the shadowRay intersects any scene object EXCEPT our light source...
                context ctx = {{0},{0},{0},{0},{0}};
                intersection ixs = findClosestIntersection(objects, numObjects, groups, triangles, shadowRayOrigin, shadowRayDirection, &ctx);
                if (ixs.lowestIntersectionIndex == l && ixs.t > EPSILON) {
                    double4 effectiveColor = b->color * objects[l].emission;

                    // I've seen this as well:
                    // l += light.getPower() * cos * cosp * rectangle.getArea() / lengthSquared;
                    // perhaps use the surface area of the light's hemisphere and divide by t*t?
                    // 2*Pi*r2
                    //double attenuation = 2*PI*objects[0].transform[0]*objects[0].transform[0] / ((0.25+ixs.t)*(0.25+ixs.t));

                    // Christian's attenuation based on % of hemisphere which is covered by light source.
                    // Note 8 months later: I can't figure out why I'm using that value from the object's transform...
                    // ..it may be a trick to not accidently divide by zero? But what if x is == 0 and t is 0????
                    double attenuation = 1 - ixs.t / sqrt(ixs.t*ixs.t + objects[l].transform[0]*objects[l].transform[0]);

                    // Compute and update the accumulate color pointer passed to the function
                    *accumColor += effectiveColor * lightDotNormal * mask * attenuation;
                }
            }
        }
    }
}

// From OpenGL bi-directional path tracer https://www.shadertoy.com/view/MtfGR4, for inspiration on
// how to build the light path.
void constructLightPath(inout float seed) {
    // start by creating a ray origin in unit space where x,y,z is between -0.5 to 0.5
    vec3 rayOrigin = normalize( hash3(seed)-vec3(0.5) );

    // next, generate a random ray in a hemisphere, probably in the hemisphere as defined by the
    // vector from 0,0,0 -> rayOrigin.xyz
    vec3 rayDirection = randomHemisphereDirection( rayOrigin, seed );

    // Move the ray origin into world coordinates, (where the light bulb is)
    // and modify its location by "half" its local space???
    rayOrigin = lightSphere.xyz + rayOrigin*0.5;
    vec3 color = LIGHTCOLOR;

    // Init first light path node
    lpNodes[0].position = rayOrigin;
    lpNodes[0].color = color;
    lpNodes[0].normal = rayDirection;

    // initialize each node in the expected path with 0,0,0 for position, color (black) and normal.
    for( int i=1; i<LIGHTPATHLENGTH; ++i ) {
        lpNodes[i].position = lpNodes[i].color = lpNodes[i].normal = vec3(0.);
    }

    // Start building the light path
    for( int i=1; i<LIGHTPATHLENGTH; i++ ) {
		vec3 normal;

		// Intersect world objects with the ray. Note that res.x == t (distance from origin to intersection)
		// y is a "material index" that translates to a fixed color used in the scene for each primitive object.
        vec2 res = intersect( rayOrigin, rayDirection, normal );

        // ugly hack to never intersect the light source (it has material index=4)
        if( res.y > -0.5 && res.y < 4. ) {
            // set new rayOrigin using good ol' last origin + distance along direction vector.
            rayOrigin = rayOrigin + rayDirection*res.x;

            // get intersected object color.
            // Note that color starts with LIGHTCOLOR (16.86, 10.76, 8.2)*1.3
            // and is multiplied here, so "color" seems to be a bit like a mask where each path along the light ray will
            // get less light contribution.
            color *= calcColor( res.y );
            lpNodes[i].position = rayOrigin;  // at next (starts at 1) store position, color and normal.
            lpNodes[i].color = color;
            lpNodes[i].normal = normal; // note that the value of normal was populated in the interection code

            // continue by picking a new rayDirection in the hemisphere of the new normal.
            rayDirection = cosWeightedRandomHemisphereDirection( normal, seed );
        } else break;
    }
}

// the sampler is used to "pick" colors from textures using normalized (e.g. floating point) coordinates where
// CLK_ADDRESS_REPEAT makes sure that we don't get "mirrored" textures when crossing the 1.0 or 0.0 boundaries.
__constant sampler_t sampler = CLK_NORMALIZED_COORDS_TRUE | CLK_ADDRESS_REPEAT | CLK_FILTER_LINEAR;

__kernel void trace(__constant object *global_objects, unsigned int numObjects, __global triangle *triangles, __global group *groups, __global double *output,
                    __constant double *seedX, unsigned int samples, __global camera *cam, unsigned int yOffset,
                    image2d_array_t image, image2d_array_t sphereTextures, image2d_array_t cubeMapTextures) {

    // int skipped = 0;
    // int hit = 0;
    double colorWeight = 1.0 / samples;
    int i = get_global_id(0);
    __local float fgi, fgi2;
    fgi = seedX[i] / numObjects;
    fgi2 = seedX[i] / samples;
    double4 originPoint = (double4)(0.0f, 0.0f, 0.0f, 1.0f);
    double4 colors = (double4)(0, 0, 0, 0);

    // experiment: copy objects to local memory. May actually be faster, at least on CPU?
    __local object objects[16];
    for (unsigned int a = 0;a < numObjects;a++) {
        objects[a] = global_objects[a];
    }

    // START BI-DIRECTIONAL EXPERIMENT!
    // 1. Pick a random outgoing ray from the sphere light source. In the future, perhaps we should
    //    limit so we only emit light rays in useful directions... For now, the light is fully inside the cornell box.

    // 2. Follow the ray around the scene, record position, obj color, cosine, normal.

    // 3. Then cast a ray from the camera, and let it (likewise) bounce around.

    // 4. Finally, iterate over the _camera_ rays, and for each bounce (not origin), cast a shadow ray
    //    to each (including point on light light source) vertex on the light ray.

    // 4.1 For each shadow ray that intersects a light vertex, accumulate color and emission. HOW?!?!

    // 4.2 Continue doing this for each path on the camera ray.
    // 4.3 Finally, we should be able to sum together all collected color and emission, weighed by the usual
    //     cosine stuff and get a much better result than naive path tracing.
    //     Compared to next-event estimation, that only cast a shadow ray to a random point on each light source,
    //     BPT should be able to accumulate emission from several verticies on the light path.
    double lightScale = 0.15;
    double4 lightPos = (double4)(objects[0].transform[3], objects[0].transform[7], objects[0].transform[11], 1.0);

    // start by creating a ray origin in world space at random point on spherical light source
    double4 rpos = randomPointOnSphere(1.0, noise3D(fgi, fgi*fgi2, fgi2), noise3D(fg2, fgi, fgi*fgi*1.4324);
    double4 rayOrigin = lightPos + (rpos * lightScale);

    // then compute the resulting random ray direction, with rayOrigin adjusted to be an overpoint to
    // avoid self-intersection
    double4 rayDirection = normalize(lightPos - rayOrigin);
    rayOrigin = rayOrigin + rayDirection*EPSILON; // overpoint

    //    vec3 color = LIGHTCOLOR;

    // allow up to 16 verticies on the light path
    __local bounce lpNodes[16]; // = {};
    lpNodes[0] = {rayOrigin, 1.0, normalize(rayDirection), objects[0].emission, normalize(rayDirection), 1.0, false};

    // Init first light path node
//    lpNodes[0].position = rayOrigin;
//    lpNodes[0].color = color;
//    lpNodes[0].normal = rayDirection;

    // initialize each node in the expected path with 0,0,0 for position, color (black) and normal.
//    for( int i=1; i<LIGHTPATHLENGTH; ++i ) {
//        lpNodes[i].position = lpNodes[i].color = lpNodes[i].normal = vec3(0.);
//    }
    unsigned int LIGHTPATHLENGTH = 6;
    // Start building the light path
    for( unsigned int i=1; i<LIGHTPATHLENGTH; i++ ) {
        vec3 normal;

        // Intersect world objects with the ray. Note that res.x == t (distance from origin to intersection)
        // y is a "material index" that translates to a fixed color used in the scene for each primitive object.
        vec2 res = intersect( rayOrigin, rayDirection, normal );

        // ugly hack to never intersect the light source (it has material index=4)
        if( res.y > -0.5 && res.y < 4. ) {
            // set new rayOrigin using good ol' last origin + distance along direction vector.
            rayOrigin = rayOrigin + rayDirection*res.x;

            // get intersected object color.
            // Note that color starts with LIGHTCOLOR (16.86, 10.76, 8.2)*1.3
            // and is multiplied here, so "color" seems to be a bit like a mask where each path along the light ray will
            // get less light contribution.
            color *= calcColor( res.y );
            lpNodes[i].position = rayOrigin;  // at next (starts at 1) store position, color and normal.
            lpNodes[i].color = color;
            lpNodes[i].normal = normal; // note that the value of normal was populated in the interection code

            // continue by picking a new rayDirection in the hemisphere of the new normal.
            rayDirection = cosWeightedRandomHemisphereDirection( normal, seed );
        } else break;
    }




    // END LIGHT PATH!


    __local intersection ixs;

    // get current x,y coordinate from i given image width
    unsigned int x = i % cam->width;
    unsigned int y = yOffset + i / cam->width;

// Comment in to debug a single pixel
//    if (x != 428 || y != 558) {
//        output[i * 4] = 1.0;
//        output[i * 4 + 1] = 1.0;
//        output[i * 4 + 2] = 1.0;
//        output[i * 4 + 3] = 1.0;
//        return;
//    }

    __local double4 rayOrigin, rayDirection;
    for (unsigned int n = 0; n < samples; n++) {
        // For each sample, compute a new ray cast through the target (x,y) pixel with random offset within the pixel.
        ray r = rayForPixel(x, y, *cam, noise3D(fgi, n, fgi2), noise3D(fgi, fgi2, n), n, samples);
        rayOrigin = r.origin;
        rayDirection = r.direction;

        unsigned int actualBounces = 0;
        unsigned int effectiveBounces = 0;
        // Each ray may bounce up to 16 times
        __local bounce bounces[16]; // = {};
        bool entering = false;
        bool inside = false;
        bool exiting = false;
        bool reflecting = false;

        // For each ray, allow up to MAX_BOUNCES bounces, with a cap of MAX_EFFECTIVE_BOUNCES since refraction
        // does not "consume" a color-contributing "effective" bounce.
        for (unsigned int b = 0; b < MAX_BOUNCES && effectiveBounces < MAX_EFFECTIVE_BOUNCES ; b++) {

            context ctx = {{0},{0},{0},{0},{0}};
            ixs = findClosestIntersection(objects, numObjects, groups, triangles, rayOrigin, rayDirection, &ctx);

            if (ixs.lowestIntersectionIndex > -1) {
                object obj = objects[ixs.lowestIntersectionIndex];

                // Remember that we use the untransformed ray here!

                // Position gives us the intersection position along RAY at T
                double4 position = rayOrigin + rayDirection * ixs.t;

                // The vector to the eye (or last bounce origin) is exactly the opposite
                // of the ray direction
                double4 eyeVector = -rayDirection;

                // object normal at intersection: Transform point from world to object
                // space
                double4 objectNormal;

                // PLANE always have its normal UP in local space (unless we have a normal map)
                if (obj.type == 0) {
                    if (obj.isTexturedNM) {
                        double4 localPoint = mul(obj.inverse, position);
                        float4 rgba = read_imagef(image, sampler, (float4)(fabs(localPoint.x) * obj.textureScaleXNM, fabs(localPoint.z) * obj.textureScaleYNM, obj.textureIndexNM, 0));
                        objectNormal = (double4)(rgba.x, rgba.y, rgba.z, 0.0);
                        objectNormal = normalize(objectNormal);
                    } else {
                        objectNormal = (double4)(0.0, 1.0, 0.0, 0.0);
                    }
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
                    if (dist < 1 && localPoint.y >= obj.maxY - EPSILON) {
                        objectNormal = (double4)(0.0, 1.0, 0.0, 0.0);
                    } else if (dist < 1 && localPoint.y <= obj.minY + EPSILON) {
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
                        objectNormal = (double4)(0.0, localPoint.y, 0.0, 0.0);
                    } else {
                        objectNormal = (double4)(0.0, 0.0, localPoint.z, 0.0);
                    }
                } else if (obj.type == 4) {
                    // GROUP, which in practice means a triangle, whose normal is typically pre-populated in N and stored in xsTriangles
                    objectNormal = ctx.xsTriangle[ixs.normalIndex];
                }
                // Finish the normal vector by multiplying it back into world coord
                // using the inverse transpose matrix and then normalize it
                double4 normalVec = mul(obj.inverseTranspose, objectNormal);
                normalVec.w = 0.0; // set w to 0
                normalVec = normalize(normalVec);

                // negate the normal if the normal if facing
                // away from the "eye"
                if (dot(eyeVector, normalVec) < 0.0) {
                    normalVec = normalVec * -1.0;
                }

                // Compute the over point, with a slight offset along the normal, in
                // order to avoid self-intersection on the next bounce.
                double4 overPoint = position + normalVec * EPSILON;

                // Prepare the outgoing ray (next bounce) by reusing the original ray, just
                // update its origin and direction.

                // Impl here supports either diffuse or reflected, but for obj.reflectivity > 0 a proportionate portion of samples
                // will diffuse instead of reflect. Poor-man's BRDF
                double cosine = 1.0; // experiment: for reflected, always use 1.0
                entering = false;
                exiting = false;
                reflecting = false;
                double sch = 0.0;

                // First, decide to refract or reflect depending on material properties.
                if (obj.reflectivity != 0.0 && noise3D(fgi, n, b) < obj.reflectivity) {
                    // reflect, even if transparent.
                    // Reflected, calculate reflection vector and set as rayDirection
                    double dotScalar = dot(rayDirection, normalVec);
                    double4 norm = (normalVec * 2.0) * dotScalar;
                    rayDirection = rayDirection - norm;
                    reflecting = true;
                }
                 else if (obj.refractiveIndex == -1.0) {
                    // Slightly hacky - a refractive index of -1.0 means we have a super-thin material that should be handled
                    // as a "refraction without refraction", e.g. transparent but won't affect the ray direction.

                      if (schlick(eyeVector, normalVec,  1.0, 1.5) < noise3D(fgi, n*n, b)) {
                          // passing through, set underpoint
                          overPoint = position - normalVec * EPSILON;
                          // do not touch rayDirection
                      } else {
                          // reflected
                          double dotScalar = dot(rayDirection, normalVec);
                          double4 norm = (normalVec * 2.0) * dotScalar;
                          rayDirection = rayDirection - norm;
                          reflecting = true;
                      }
                }
                // Consider removing this HACK for handling glass models without thickness.
                else if (obj.refractiveIndex != 1.0) {
                    // Handle "normal" refraction for solid objects

                    if (!inside) {
                        // if we have hit a refractive object and we're not inside one...

                        // compute schlick to determine chance of reflection
                        sch = schlick(eyeVector, normalVec,  1.0, obj.refractiveIndex);
                        double rnd = noise3D(fgi, n*n, b);
                         if (x == 428 && y == 591) {
                            printf("NOT INSIDE: schlick was %f, chance was %f\n", sch, rnd);
                         }
                        if (sch < rnd) {
                            // refraction
                            rayDirection = computeRefractedRay(eyeVector, normalVec,  1.0, obj.refractiveIndex);
                            overPoint = position - normalVec * EPSILON;
                            inside = true;
                            entering = true;
                            exiting = false;
                        } else {
                            // reflection
                            double dotScalar = dot(rayDirection, normalVec);
                            double4 norm = (normalVec * 2.0) * dotScalar;
                            rayDirection = rayDirection - norm;
                            reflecting = true;
                        }
                    } else {
                        // If already inside, we are passing back into air but we may still reflect internally in the medium??
                         double sch = schlick(eyeVector, normalVec,  obj.refractiveIndex, 1.0);
                         if (x == 378 && y == 558) {
                             printf("IS INSIDE: schlick was %f\n", sch);
                          }
                         if (sch < noise3D(fgi, n*n, b)) {
                            // refract back into air
                            rayDirection = computeRefractedRay(eyeVector, normalVec,  obj.refractiveIndex, 1.0);
                            overPoint = position - normalVec * EPSILON;
                            inside = false;
                            entering = false;
                            exiting = true;
                         } else {
                            // internal reflection??
                            double dotScalar = dot(rayDirection, normalVec);
                            double4 norm = (normalVec * 2.0) * dotScalar;
                            rayDirection = rayDirection - norm;
                            entering = false;
                            exiting = false;
                            reflecting = true;
                         }
                    }
                } else {
                    // Diffuse
                    rayDirection = randomVectorInHemisphere(normalVec, fgi, b, n);
                    // Calculate the cosine of the OUTGOING ray in relation to the surface
                    // normal.
                    cosine = dot(rayDirection, normalVec);
                }
                rayOrigin = overPoint;

                // 378 , 591
                if (x == 428 && y == 558) {
                    printf("iteration: %d === intersected: %s === schlick: %f ===new origin: %f, %f, %f ==== direction: %f %f %f\n", b, obj.label,sch, rayOrigin.x, rayOrigin.y, rayOrigin.z, rayDirection.x, rayDirection.y, rayDirection.z);
                }

                // Finish this iteration by storing the bounce. Objects (with triangles) gets special treatment
                // since a model may have many different materials. See xsTriangleColor
                if (obj.type == 4) {
                    bounce bnce = {position, cosine, ctx.xsTriangleColor[ixs.normalIndex], ctx.xsTriangleEmission[ixs.normalIndex], normalVec, 1.0, entering || exiting};
                    bounces[b] = bnce;
                } else {
                    // texture experiment for PLANE, CUBE and SPHERE
                    double4 color = obj.color;
                    if (obj.isTextured) {
                          if (obj.type == 0) { // PLANE
                              double4 localPoint = mul(obj.inverse, position);
                              float4 rgba = read_imagef(image, sampler, (float4)(localPoint.x * obj.textureScaleX, localPoint.z * obj.textureScaleY, obj.textureIndex, 0));
                              color = (double4)(rgba.x, rgba.y, rgba.z, 1.0);
                          } else if (obj.type == 1) { // SPHERE
                              double4 localPoint = mul(obj.inverse, position);
                              double2 uv = sphericalMap(localPoint);
                              float4 rgba = read_imagef(sphereTextures, sampler, (float4)(uv.x, 1.0-uv.y, obj.textureIndex, 0));
                              color = (double4)(rgba.x, rgba.y, rgba.z, 1.0);
                          } else if (obj.type == 3) { // CUBE
                              double4 localPoint = mul(obj.inverse, position);
                              double2 uv = cubeUV(localPoint);
                              float4 rgba = read_imagef(cubeMapTextures, sampler, (float4)(uv.x, uv.y, obj.textureIndex, 0));
                              color = (double4)(rgba.x, rgba.y, rgba.z, 1.0);
                          }
                    }
                    bounce bnce = {position, cosine, color, obj.emission, normalVec, 1.0, entering || exiting};
                    bounces[b] = bnce;
                }

                // Only increment effective bounces for non-refractive/reflective materials
                if (!entering && !exiting && !reflecting) {
                    effectiveBounces++;
                }

                // increment total bounces.
                actualBounces++;

                // experiment - stop bouncing if intersecting a light source
                if (obj.emission.x > 0.0) {
                    break;
                }
            }
        }

        // ------------------------------------ //
        // Calculate final color using bounces! //
        // ------------------------------------ //
        double4 accumColor = (double4)(0.0, 0.0, 0.0, 0.0);
        double4 mask = (double4)(1.0, 1.0, 1.0, 1.0);
        unsigned int imageX = x;
        for (unsigned int x = 0; x < actualBounces; x++) {

            // first run - just use the material color of the first bounce
            //accumColor = bounces[x].color;
            //break;

            // second run - as above, but use cos for the OUTGOING ray from the object's normal.
            // this gives a slightly noisy result since the outgoing ray is random.
            //accumColor = bounces[x].color * bounces[x].cos;
            //break;

            // third run - we actually use the cos of the _incoming_ ray instead
            //accumColor = bounces[x].color * bounces[x].inCosine;
            //break;

            // fourth run - just add all FLAT colors together and average them
            //accumColor += bounces[x].color / actualBounces;

            // fifth run - just add colors together multiplied by OUT ray and average them
            //accumColor += bounces[x].color*bounces[x].inCosine / actualBounces;

            // sixth run - sample the (point) light source on each bounce

            bounce bnce = bounces[x];
            if (imageX == 428 && y == 558) {
                printf("bounce: %d === refraction: %d \n", x, bnce.isRefraction);
            }

            // when refracting, simply pass updating color, mask etc for this bounce.
            if (bnce.isRefraction) {
                continue;
            }

            // add "strength" multiplied by remaining mask to accumColor.
            accumColor = accumColor + mask * bnce.emission;

            // If sampling a light source, ignore further bounces
            if (bnce.emission.x > 0.0) {

                // direct sampling of a light source
                if (x == 0) {
                    accumColor = bnce.color;// original just used emission here.
                }
                break;
            }


            // Here is the next event estimation experiment:  iterate over all light sources in the scene, accumulate light
            // from all, updating accumColor. Works well for diffuse materials, but not for reflections/refraction.
            // nextEventEstimation(objects, numObjects, groups, triangles, &bnce, fgi, fgi2, n, mask, x, &accumColor);

            // Update the mask by multiplying it with the hit object's color
            mask *= bnce.color;

            // perform cosine-weighted importance sampling by multiplying the mask
            // with the cosine. Note to self: For refracting/reflection, we set cos to 1.0.
            mask *= bnce.cos;
        }

        // Finish this "sample" by adding the accumulated color to the total
        colors += accumColor;
    }

    // Finish the pixel by multiplying each RGB component by its total fraction and
    // store in the output buffer.
    output[i * 4] = colors.x * colorWeight;
    output[i * 4 + 1] = colors.y * colorWeight;
    output[i * 4 + 2] = colors.z * colorWeight;
    output[i * 4 + 3] = 1.0;
}