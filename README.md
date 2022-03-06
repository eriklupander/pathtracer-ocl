# pathtracer-ocl
Pathtracer written in Go and OpenCL
![example](images/aa-with-box-and-cyl.png)
_(2048 samples)_

## Description
Simple unidirectional pathtracer written just for fun using Go as frontend and OpenCL as computation backend.

Supports:
* Spheres, Planes, Boxes, Cylinders
* Diffuse and reflective materials
* Movable camera
* Anti-aliasing
* Depth of Field with simple focal length and camera aperture.

Based on or inspired by:

* My implementation of "The Ray Tracer Challenge" at https://github.com/eriklupander/rt
* Mask/accumulated color shading by Sam Lapere at https://raytracey.blogspot.com/2016/11/opencl-path-tracing-tutorial-2-path.html
* Ray in hemisphere code by Hunter Loftis at https://github.com/hunterloftis/pbr/blob/1ce8b1c067eea7cf7298745d6976ba72ff12dd50/pkg/geom/dir.go
* And my own mashup of the three above, a simple and Go-native path-tracer https://github.com/eriklupander/pathtracer

Next steps:
* Groups of primitives, with bounding boxes
* Triangle primitives
* .obj model loading into bounding volume hierarchies 
* Rendering models

All of the above is present in my Go-only ray-tracer, but given that recursion is forbidden in OpenCL, as well as variable-length arrays cannot be passed to OpenCL without careful management of struct sizes, "numberOfNN" fields etc, incorporating 3D model rendering with acceptable performance is non-trivial.

The overall solution is _probably_ to:
* Pass ALL triangles (for all models) in a long continous array, each triangle will have pre-computed vertex/surface normals etc and consume 256 or 512 bytes each.
* Organize models into a BVH hierarchy of groups. 
  * All groups goes into a contingous array of groups.
  * Each group has a bounding box
  * Each group has a int[16] array for indexes to subgroups
  * Each group has two fields: trianglesOffset and trianglesSize, these are used to reference into the `triangles` array. I.e. in "go terms", `triangles[trianglesOffset:trianglesOffset+trianglesSize]`.
* A challenge here is that traversing the BVH must be done using for-loops and a local "stack" rather than recursion.
  * Also, since each "subgroup" may have its own "subtranslation" in relation to its parent, we may need to keep a stack of translations around as well.
    * It may be possible to let each group within a group (forming a mesh) to have an absolute world-based transform rather than one in relation to its parent.

## Usage
A few command-line args have been added to simplify testing things.

```
      --width int            Image width (default 640)
      --height int           Image height (default 480)
      --samples int          Number of samples per pixel (default 1)
      --aperture float       Aperture. If 0, no DoF will be used. Default: 0
      --focal-length float   Focal length. Default: 0
      --device-index int     Use device with index (use --list-devices to list available devices)
      --list-devices         List available devices
```
Suggested values for focal length and aperture for the standard cornell box: 1.6 and 0.1

Example:
```shell
go run cmd/pt/main.go --samples 2048 --aperture 0.15 --focal-length 1.6 --width 1280 --height 960
```

### Listing and selecting a device
Not all OpenCL devices are created equal. On the author's semi-ancient MacBook Pro 2014, running `go run cmd/pt/main.go --list-devices` yields:
```shell
Index: 0 Type: CPU Name: Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
Index: 1 Type: GPU Name: Iris Pro
Index: 2 Type: GPU Name: GeForce GT 750M
```
However, the Iris Pro iGPU does not support double-precision floating point numbers. Also, there are subtle differences between CPU and GPU device, which in certain situations may result in panics or segmentation faults. In other words: Your milage may vary. CPU-based devices seems to be the most stable and on MacBooks, CPU has significantly better performance than the discrete GPUs.

## Performance
For this _reference image_ at 1280x960:
![example](images/reference.png)

### MacBook Pro mid-2014
* Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz:  10m45.813990173s
* GeForce GT 750M GPU:                        14m12.049519483s

### MacBook Pro 2019
* Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz:     4m20.568928052s
* AMD Radeon Pro 560X Compute Engine:           5m14.439406471s

### Desktop PC with Windows 10 - Ryzen 2600X
* NVIDIA GeForce RTX 2080:                      45.4309853s

_(AMD has dropped Ryzen CPU OpenCL support on Windows)_

In this scenario, the 8-core Core i9 CPU is more than twice as fast as the 4-cire Core i7 CPU on the older MacBook. Both mGPUs are slower than their respective CPUs.

The king is unsurprisingly enough the GeForce RTX 2080 on my Desktop PC, which is almost 6x faster than the 8-core Intel CPU.

## Issues
The current DoF has some issues producing slight artifacts, probably due to how random numbers are seeded for the aperture-based ray origin.

## Gallery
### Depth of Field
Depth-of-field effect is accomplished through casting a standard camera->pixel ray into the scene, and then creating a new "focal point" by using focal distance (distance to a point along camera ray). A new random origin point is then randomly picked around the camera origin with r==aperture and a _new_ ray is cast from the new camera through the focal point, resulting in objects not near the focal point to appear increasingly out-of focus.

1280x960, 2048 samples, focal length 1.6, aperture 0.15.
![DoF](images/DoF-2048.png)

### Anti-aliasing
Anti-aliasing is accomplished through the age-old trick of casting each ray through a random location in the pixel. Given enough samples, an anti-aliased effect will occurr.

Examples rendered in 640x480 with 512 samples:
#### Without anti-aliasing:
![noaa](images/no-aa.png)

#### With anti-aliasing:
![noaa](images/aa.png)