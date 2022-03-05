# pathtracer-ocl
Pathtracer written in Go and OpenCL
![example](images/aa-with-box-and-cyl.png)
## Description
Simple unidirectional pathtracer written just for fun using Go as frontend and OpenCL as computation backend.

Supports:
* Spheres, Planes, Boxes, Cylinders
* Diffuse and reflective materials
* Movable camera
* Anti-aliasing

Based on or inspired by:

* My implementation of "The Ray Tracer Challenge" at https://github.com/eriklupander/rt
* Mask/accumulated color shading by Sam Lapere at https://raytracey.blogspot.com/2016/11/opencl-path-tracing-tutorial-2-path.html
* Ray in hemisphere code by Hunter Loftis at https://github.com/hunterloftis/pbr/blob/1ce8b1c067eea7cf7298745d6976ba72ff12dd50/pkg/geom/dir.go
* And my own mashup of the three above, a simple and Go-native path-tracer https://github.com/eriklupander/pathtracer

## Performance
For this _reference image_ at 1280x960:
![example](images/reference.png)

MacBook Pro mid-2014:
Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz:  10m45.813990173s
GeForce GT 750M GPU:                        14m12.049519483s