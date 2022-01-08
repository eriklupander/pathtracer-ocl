# Performance stuff

Reference image with cornell box with two diffuse spheres, 1280x960
### MacBook Pro mid-2014
* Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz:  10m45.813990173s
* GeForce GT 750M GPU:                        14m12.049519483s

### MacBook Pro 2019
* Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz:     4m20.568928052s
* AMD Radeon Pro 560X Compute Engine:           5m14.439406471s

### Windows 10 - Ryzen 2600X
* NVIDIA GeForce RTX 2080:                      45.4309853s

_(AMD has dropped Ryzen OpenCL support on Windows)_



Branch compare:
with mapping array, RTX2080 1024 samples: 6.5 seconds
NO mapping array, RTX2080 1024 samples: 6.7 seconds

with mapping array, GeForce GT 750M 32 samples: 7 seconds
NO mapping array, GeForce GT 750M 32 samples: 4 seconds

with mapping array, Intel(R) Core(TM) i7-4870HQ 32 samples: 3 seconds
NO mapping array, Intel(R) Core(TM) i7-4870HQ 32 samples: 3 seconds