__constant float PI = 3.14159265359f;
__constant unsigned int SAMPLE_COUNT = 256;
__constant unsigned int MAX_BOUNCES = 5;

typedef struct tag_ray{
	float4 origin;
	float4 direction;
} ray;

typedef struct __attribute__ ((packed)) tag_object{
	float16 transform;        // 64 bytes 16x4
    float16 inverse;          // 64 bytes
    float16	inverseTranspose; // 64 bytes
    float4	color;            // 16 bytes
    float4	emission;         // 16 bytes
    float	refractiveIndex;  // 4 bytes
    int 	type;             // 4 bytes
    float4	padding1;          // 16 bytes
    float padding2;           // 4 bytes
    float padding3;           // 4 bytes
} object;

typedef struct tag_intersection{
  float objType;
  float t;
} intersection;

typedef struct tag_bounce{
	float4 point;
	float cos;
	float4 color;
	float4 emission;
	//diffuse         bool
	//refractiveIndex float64
} bounce;

static float noise3D(float x, float y, float z) {
    float ptr = 0.0f;
    return fract(sin(x*112.9898f + y*179.233f + z*237.212f) * 43758.5453f, &ptr);
}

// randomVectorInHemisphere translated into Go from https://raytracey.blogspot.com/2016/11/opencl-path-tracing-tutorial-2-path.html
// The thing is that using this func for diffuse surfaces produces a good and balanced result in the final image,
// while using the randomConeInHemisphere func translated from Hunter Loftis PathTracer produces overexposed highlights.
//
// Need to figure out why.
inline float4 randomVectorInHemisphere(float4 normalVec, float x, float y, float z) {
	
    float rand1 = 2.0 * PI * noise3D(x,y,z);
	float rand2 = noise3D(y,z,x);
	float rand2s = sqrt(rand2);
    
	/* create a local orthogonal coordinate frame centered at the hitpoint */
    float4 axis;
	if (fabs(normalVec.x) > 0.1) {
		axis = (float4)(0.0, 1.0, 0.0, 0.0);
	} else {
		axis = (float4)(1.0, 0.0, 0.0, 0.0);
	}
	float4 u = normalize(cross(axis, normalVec));
	float4 v = cross(normalVec, u);

	/* use the coordinate frame and random numbers to compute the next ray direction */
    return u * cos(rand1)*rand2s + v * sin(rand1)*rand2s + normalVec * sqrt(1.0f-rand2);
}

static float get_random(unsigned int *seed0, unsigned int *seed1) {

	/* hash the seeds using bitwise AND operations and bitshifts */
	*seed0 = 36969 * ((*seed0) & 65535) + ((*seed0) >> 16);  
	*seed1 = 18000 * ((*seed1) & 65535) + ((*seed1) >> 16);

	unsigned int ires = ((*seed0) << 16) + (*seed1);

	/* use union struct to convert int to float */
	union {
		float f;
		unsigned int ui;
	} res;

	res.ui = (ires & 0x007fffff) | 0x40000000;  /* bitwise AND, bitwise OR */
	return (res.f - 2.0f) / 2.0f;
}

inline float3 rndVec(float4 normalVec, unsigned int* seed0, unsigned int* seed1) {
    /* compute two random numbers to pick a random point on the hemisphere above the hitpoint*/
    float rand1 = 2.0f * PI * get_random(seed0, seed1);
    float rand2 = get_random(seed0, seed1);
    float rand2s = sqrt(rand2);
    

    /* create a local orthogonal coordinate frame centered at the hitpoint */
    float3 w = (float3)(normalVec.x, normalVec.y, normalVec.z);
    float3 axis = fabs(w.x) > 0.1f ? (float3)(0.0f, 1.0f, 0.0f) : (float3)(1.0f, 0.0f, 0.0f);
    float3 u = normalize(cross(axis, w));
    float3 v = cross(w, u);

    /* use the coordinte frame and random numbers to compute the next ray direction */
    return normalize(u * cos(rand1)*rand2s + v*sin(rand1)*rand2s + w*sqrt(1.0f - rand2));
}


// inline float4 mulMatrix(__global float16 mat, unsigned int index, float4 vec) {
// 	float4 result;
// 	for (unsigned int row = 0; row < 4; row++) {
// 		float a = mat[index + (row*4)+0] * vec.x;
// 		float b = mat[index + (row*4)+1] * vec.y;
// 		float c = mat[index + (row*4)+2] * vec.z;
// 		float d = mat[index + (row*4)+3] * vec.w;
// 		result[row] = a + b + c + d;
// 	}
// 	return result;
// }

inline float4 mul(float16 mat, float4 vec) {
    float4 result;
    for (unsigned int row = 0; row < 4; row++) {
        float a = mat[(row*4)+0] * vec.x;
        float b = mat[(row*4)+1] * vec.y;
        float c = mat[(row*4)+2] * vec.z;
        float d = mat[(row*4)+3] * vec.w;
        result[row] = a + b + c + d;
    }
    return result;
}


__kernel void trace(
   __global ray* rays2, 
   __constant object* objects,
   const unsigned int objectCount,
   __global float* output,
   __global float* seedX)
{
    int i = get_global_id(0);
    unsigned int seed0 = i;
    unsigned int seed1 = i % 32;
    float fgi = float(seedX[i])/objectCount;
    float4 originPoint = (float4)(0.0f, 0.0f, 0.0f, 1.0f);
	    
	int vOffset = i * 4;
    float4 colors = (float4)(0.0f,0.0f,0.0f,0.0f);

    if (i == 0) {
    
        for (unsigned int y = 0;y<objectCount;y++) {
            printf("obj: transform:\n");
            for (unsigned int x = 0;x<16;x++) {
                printf("object: %d elem: %d: %f ", y, x, objects[y].transform[x]);
                printf("\n");
            }
            printf("\n");
            printf("obj: inverse:\n");
            for (unsigned int x = 0;x<16;x++) {
                printf("object: %d elem: %d: %f ", y, x, objects[y].inverse[x]);
                printf("\n");
            }
            printf("\n");
            printf("obj: inverse transpose:\n");
            for (unsigned int x = 0;x<16;x++) {
                printf("object: %d elem: %d: %f ", y, x, objects[y].inverseTranspose[x]);
                printf("\n");
            }
            printf("\n");
            printf("obj: color:\n");
            for (unsigned int x = 0;x<4;x++) {
                printf("object: %d elem: %d: %f ", y, x, objects[y].color[x]);
                printf("\n");
            }
            printf("\n");
            printf("obj: emission:\n");
            for (unsigned int x = 0;x<4;x++) {
                printf("object: %d elem: %d: %f ", y, x, objects[y].emission[x]);
                printf("\n");
            }
        }
        printf("\n");
    }
    	
    for (unsigned int samples = 0; samples < SAMPLE_COUNT;samples++) {
        // Each new sample needs to reset to original ray
        float4 rayOrigin =  (float4)(rays2[i].origin.x, rays2[i].origin.y, rays2[i].origin.z, rays2[i].origin.w);
	    float4 rayDirection = (float4)(rays2[i].direction.x, rays2[i].direction.y, rays2[i].direction.z, rays2[i].direction.w);

        // for each bounce...
        unsigned int actualBounces = 0;
        // Each ray may bounce up to 5 times
        bounce bounces[5] = {};
        for (unsigned int b = 0; b < MAX_BOUNCES; b++) {
            float intersections[16] = { 0 };
            float intersectedObjectType = -1.0;
            intersection xs[16] = { };

            // ----------------------------------------------------------
            // Loop through scene objects in order to find intersections
            // ----------------------------------------------------------
            for (unsigned int j = 0; j < objectCount; j++) {
                int objType = objects[j].type;
                
                // translate our ray into object space by multiplying ray pos and dir with inverse object matrix
                float4 tRayOrigin = mul(objects[j].transform, rayOrigin);
                float4 tRayDirection = mul(objects[j].transform, rayDirection);
                
                // Intersection code
                if (objType == 0) {  // intersect transformed ray with plane
                    if (fabs(tRayDirection.y) < 0.0001f) {
                        // did not intersect
                        intersections[j] = 0;
                    } else {
                        float t = -tRayOrigin.y / tRayDirection.y;
                        
                        intersections[j] = t;
                        intersectedObjectType = 0.0;
                        intersection ixs = {0.0, t};
                        xs[j] = ixs;
                    }
                }

                if (objType == 1) {    // SPHERE
                    // this is a vector from the origin of the ray to the center of the sphere at 0,0,0
                    float4 vecToCenter = tRayOrigin - originPoint;

                    // This dot product is
                    float a = dot(tRayDirection, tRayDirection);
                    
                    // Take the dot of the direction and the vector from ray origin to sphere center times 2
                    float b = 2.0 * dot(tRayDirection, vecToCenter);

                    // Take the dot of the two sphereToRay vectors and decrease by 1 (is that because the sphere is unit length 1?
                    float c = dot(vecToCenter, vecToCenter) - 1.0;
                
                    // calculate the discriminant
                    float discriminant = (b*b) - 4*a*c;
                    if (discriminant < 0.0) {
                        
                    } else {
                        // finally, find the intersection distances on our ray.
                        float t1 = (-b - sqrt(discriminant)) / (2*a);
                        //float t2 = (-b + sqrt(discriminant)) / (2*a);
                        intersections[j] = t1;
                        intersectedObjectType = 1.0;
                        intersection ixs = {1.0, t1};
                        xs[j] = ixs;
                    }
                }
            }
            
            // find lowest positive intersection index
            float lowestIntersectionT = 999.0;
            int lowestIntersectionIndex = -1;
            for (unsigned int x = 0;x < 16;x++) {
                if (intersections[x] > 0.0001f) {
                    if (intersections[x] < lowestIntersectionT) {
                        lowestIntersectionT = intersections[x];
                        lowestIntersectionIndex = x;
                    }
                }
            }
            if (lowestIntersectionIndex > -1) {
                //printf("Ray %d intersected object %d a %d\n", i, lowestIntersectionIndex, objects[lowestIntersectionIndex].type);
                    
                // START COMPUTATIONS.
                // Remember that we use the untransformed ray here!
                
                // Position gives us the intersection position along RAY at T
                float4 position = rayOrigin + rayDirection * lowestIntersectionT;

                // The vector to the eye is exactly the opposite of the ray direction
                float4 eyeVector = -rayDirection;

                // object normal at intersection: 
                // Transform point from world to object space, including recursively traversing any parent object
                // transforms.
                float4 localPoint = mul(objects[lowestIntersectionIndex].inverse, position);
                float4 objectNormal;
                if (intersectedObjectType == 0) {
                    // PLANE always have its normal UP in local space
                    objectNormal = (float4)(0.0f, 1.0f, 0.0f, 0.0f);
                } else if (intersectedObjectType == 1) {
                    // SPHERE always has its normal from sphere center outwards to the world position.
                    objectNormal = localPoint - originPoint;
                }
                // Finish the normal vector by multiplying it back into world coord using the inverse transpose matrix
                float4 normalVec = mul(objects[lowestIntersectionIndex].inverseTranspose, objectNormal);
                normalVec.w = 0.0f; // set w to 0
                normalVec = normalize(normalVec);
                
                // reflection vector
                // float dotScalar = dot(rayDirection, normalVec);
                // float4 norm = (normalVec * 2.0) * dotScalar;
                // float4 reflectVec = rayDirection - norm;
                
                //comps.Inside = false
                // negate the normal if the normal if facing away from the "eye"
                if (dot(eyeVector, normalVec) < 0) {
                    normalVec = -normalVec;
                }

                // Perhaps only compute these if we're going to cast a new ray?
                float4 offset = normalVec * 0.0001f;
                float4 overPoint = position + offset;
                
                
                // Prepare the outgoing ray (next bounce), reuse the original ray, just update
                // its origin and direction
                //float3 tmp = rndVec(normalVec, fgi, b*3, samples*3);
                float3 rDir = rndVec(normalVec, &seed0, &seed1);//randomVectorInHemisphere(normalVec, fgi, b*3, samples*3);
                rayDirection.x = rDir.x;
                rayDirection.y = rDir.y;
                rayDirection.z = rDir.z;
                rayDirection.w = 0.0f;
                rayOrigin = overPoint;

                // Calculate the cosine of the OUTGOING ray in relation to the surface normal.
                float cosine = dot(rayDirection, normalVec);

                // Finish this iteration by storing the bounce. Alternatively, we could probably just calc
                // the color right away.
                float4 color = objects[lowestIntersectionIndex].color;
                float4 emission = objects[lowestIntersectionIndex].emission;
                bounce bnce = {position, cosine, color, emission};
                bounces[b] = bnce;
                actualBounces++;
                //printf("Ray %d intersected object %d a %d with color %f %f %f\n", i, lowestIntersectionIndex, objects[lowestIntersectionIndex].type, color.x, color.y, color.z);
                
            }
        }

        // ------------------------------------ //
        // Calculate final color using bounces! //
        // ------------------------------------ //
        float4 accumColor = (float4)(0.0f, 0.0f, 0.0f, 0.0f);
        float4 mask = (float4)(1.0f, 1.0f, 1.0f, 1.0f);
        for (unsigned int x = 0; x < actualBounces; x++) {
            
            // Start by dealing with diffuse surfaces. Note: Have no idea if my random code works at all!
            // First, ADD current color with the hadamard of the current mask and the emission properties of the hit object.
            // ctx.accumColor = geom.Add(ctx.accumColor, geom.Hadamard(ctx.mask, b.emission))
            accumColor += mask * bounces[x].emission;

            // The updated mask is used on _the next_ bounce
            // the mask colour picks up surface colours at each bounce
            //geom.HadamardPtr(&ctx.mask, &b.color, &ctx.mask)
            mask *= bounces[x].color;

            // perform cosine-weighted importance sampling for diffuse surfaces
            mask *= bounces[x].cos;
        }
        colors.x += accumColor.x;
        colors.y += accumColor.y;
        colors.z += accumColor.z;
    }

    output[vOffset] = colors.x / SAMPLE_COUNT;
    output[vOffset+1] = colors.y / SAMPLE_COUNT;
    output[vOffset+2] = colors.z / SAMPLE_COUNT;
    output[vOffset+3] = 1.0;
}