package ocl

// This file provides a broken-ish implementation using float32 rather than float64. The underlying reason is not fully
// evident, but objects being passed as structs tend to mess up some seemingly arbitrary values from Go-land to OpenCL-land,
// which does not happen with the float64 -> double codebase.
import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
	"unsafe"

	"github.com/jgillich/go-opencl/cl"
	"github.com/sirupsen/logrus"
)

type CLRay32 struct {
	Origin    [4]float32
	Direction [4]float32
}

type CLObject32 struct {
	Transform        [16]float32 // 64 bytes 16x4
	Inverse          [16]float32 // 64 bytes
	InverseTranspose [16]float32 // 64 bytes
	Color            [4]float32  // 16 bytes
	Emission         [4]float32  // 16 bytes
	RefractiveIndex  float32     // 4 bytes
	Type             int32       // 4 bytes
	Padding          [6]int32    // 24 bytes
	// Total: 64x4 == 256 bytes, nice Power of 64
}

// Trace32 is the entry point for transforming input data into their OpenCL representations, setting up boilerplate
// and calling the entry kernel. Should return a slice of float32 RGBA RGBA RGBA once finished.
func Trace32(rays []CLRay32, objects []CLObject32, width, height int) []float32 {
	dumpObjects(objects)

	platforms, err := cl.GetPlatforms()
	if err != nil {
		logrus.Fatalf("Failed to get platforms: %+v", err)
	}
	platform := platforms[0]

	devices, err := platform.GetDevices(cl.DeviceTypeAll)
	if err != nil {
		logrus.Fatalf("Failed to get devices: %+v", err)
	}
	if len(devices) == 0 {
		logrus.Fatalf("GetDevices returned no devices")
	}
	deviceIndex := 0

	if deviceIndex < 0 {
		deviceIndex = 0
	}
	device := devices[deviceIndex] // 0 == CPU 1 == iGPU 2 == GPU
	logrus.Infof("Using device %d %v", deviceIndex, devices[deviceIndex].Name())

	// 1. Select a device to use. On my mac: 0 == CPU, 1 == Iris GPU, 2 == GeForce 750M GPU
	// Use selected device to create an OpenCL context
	context, err := cl.CreateContext([]*cl.Device{device})
	if err != nil {
		logrus.Fatalf("CreateContext failed: %+v", err)
	}

	// 2. Create a "Command Queue" bound to the selected device
	queue, err := context.CreateCommandQueue(device, 0)
	if err != nil {
		logrus.Fatalf("CreateCommandQueue failed: %+v", err)
	}

	// 3.0 Read kernel source from disk
	kernelBytes, err := ioutil.ReadFile("internal/ocl/tracer_float4.cl")
	if err != nil {
		logrus.Fatalf("reading kernel source failed: %+v", err)
	}
	// 3.1 Create an OpenCL "program" from the source code.
	program, err := context.CreateProgramWithSource([]string{string(kernelBytes)})
	if err != nil {
		logrus.Fatalf("CreateProgramWithSource failed: %+v", err)
	}

	// 3.2 Build the OpenCL program (compile it?)
	if err := program.BuildProgram(nil, ""); err != nil {
		logrus.Fatalf("BuildProgram failed: %+v", err)
	}

	// 3.3 Create the actual Kernel with a name, the Kernel is what we call when we want to execute something.
	kernel, err := program.CreateKernel("trace")
	if err != nil {
		logrus.Fatalf("CreateKernel failed: %+v", err)
	}

	// 4. Some kind of error-check where we make sure the parameters passed are supported?
	for i := 0; i < 4; i++ {
		name, err := kernel.ArgName(i)
		if err == cl.ErrUnsupported {
			logrus.Errorf("GetKernelArgInfo for arg: %d ErrUnsupported", i)
			break
		} else if err != nil {
			logrus.Errorf("GetKernelArgInfo for name failed: %+v", err)
			break
		} else {
			logrus.Infof("Kernel arg %d: %s", i, name)
		}
	}

	// split work into batches in order to avoid kernels running for more than 10 seconds
	// otherwise, the GPU driver will kill us.
	results := make([]float32, 0)
	for y := 0; y < height; y += 8 {
		st := time.Now()
		results = append(results, computeBatch32(rays[y*width:y*width+width*8], objects, context, kernel, queue)...)
		logrus.Infof("%d/%d lines done in %v\n", y, height, time.Since(st))
	}

	return results
}

func dumpObjects(objects []CLObject32) {
	for y := 0; y < len(objects); y++ {
		fmt.Printf("obj: transform:\n")
		for x := 0; x < 16; x++ {
			fmt.Printf("object: %d elem: %d: %f ", y, x, objects[y].Transform[x])
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
		fmt.Printf("obj: inverse:\n")
		for x := 0; x < 16; x++ {
			fmt.Printf("object: %d elem: %d: %f ", y, x, objects[y].Inverse[x])
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
		fmt.Printf("obj: inverse transpose:\n")
		for x := 0; x < 16; x++ {
			fmt.Printf("object: %d elem: %d: %f ", y, x, objects[y].InverseTranspose[x])
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
		fmt.Printf("obj: color:\n")
		for x := 0; x < 4; x++ {
			fmt.Printf("object: %d elem: %d: %f ", y, x, objects[y].Color[x])
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
		fmt.Printf("obj: emission:\n")
		for x := 0; x < 4; x++ {
			fmt.Printf("object: %d elem: %d: %f ", y, x, objects[y].Emission[x])
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")
}

func computeBatch32(rays []CLRay32, objects []CLObject32, context *cl.Context, kernel *cl.Kernel, queue *cl.CommandQueue) []float32 {
	// 5. Time to start loading data into GPU memory
	raySizeBytes := int(unsafe.Sizeof(rays[0]))
	objectSizeBytes := int(unsafe.Sizeof(objects[0]))

	// 5.1 create OpenCL buffers (memory) for the input data. Note that we're allocating 9x bytes the size of data.
	//     since each float32 uses 8 bytes.
	inputRays, err := context.CreateEmptyBuffer(cl.MemReadOnly, raySizeBytes*len(rays))
	if err != nil {
		logrus.Fatalf("CreateBuffer failed for vectors input: %+v", err)
	}
	defer inputRays.Release()

	inputObjects, err := context.CreateEmptyBuffer(cl.MemReadOnly, objectSizeBytes*len(objects))
	if err != nil {
		logrus.Fatalf("CreateBuffer failed for objects input: %+v", err)
	}
	defer inputObjects.Release()

	seed := make([]float32, len(rays))
	for i := 0; i < len(rays); i++ {
		seed[i] = rand.Float32()
	}

	seedNumbers, err := context.CreateEmptyBuffer(cl.MemReadOnly, 4*len(seed))
	if err != nil {
		logrus.Fatalf("CreateBuffer failed for vectors input: %+v", err)
	}
	defer seedNumbers.Release()

	// 5.2 create OpenCL buffers (memory) for the output data, we want RGBA per ray, and each ray consists of 8 float32
	// However, this is the buffer size IN BYTES, so to get RGBA from let's say 1000 rays:
	// if len(ray) == 8000 we actually have 1000 "real" rays.
	// To get RGBA for 1000 real rays, we need 4000 float32, which equals 16000 bytes. Thus len(rays) * 4.
	output, err := context.CreateEmptyBuffer(cl.MemReadOnly, len(rays)*16)
	if err != nil {
		logrus.Fatalf("CreateBuffer failed for output: %+v", err)
	}
	defer output.Release()

	// 5.3 This is where we connect our input to the command queue, and upload the actual data into GPU memory
	//     The rayDataPtr:s seems to be a point to the first element of the input,
	//     while rayDataTotalSizeBytes should be the total length of the data, in bytes.

	rayDataPtr := unsafe.Pointer(&rays[0])
	rayDataTotalSizeBytes := raySizeBytes * len(rays)
	if _, err := queue.EnqueueWriteBuffer(inputRays, true, 0, rayDataTotalSizeBytes, rayDataPtr, nil); err != nil {
		logrus.Fatalf("EnqueueWriteBuffer failed: %+v", err)
	}

	dataPtrVec2 := unsafe.Pointer(&objects[0])
	dataSizeVec2 := objectSizeBytes * len(objects)
	if _, err := queue.EnqueueWriteBuffer(inputObjects, true, 0, dataSizeVec2, dataPtrVec2, nil); err != nil {
		logrus.Fatalf("EnqueueWriteBuffer failed: %+v", err)
	}

	dataPtrVec3 := unsafe.Pointer(&seed[0])
	dataSizeVec3 := int(unsafe.Sizeof(seed[0])) * len(seed)
	if _, err := queue.EnqueueWriteBuffer(seedNumbers, true, 0, dataSizeVec3, dataPtrVec3, nil); err != nil {
		logrus.Fatalf("EnqueueWriteBuffer failed: %+v", err)
	}

	// 5.4 Kernel is our program and here we explicitly bind our 4 parameters to it
	if err := kernel.SetArgs(inputRays, inputObjects, uint32(len(objects)), output, seedNumbers); err != nil {
		logrus.Fatalf("SetKernelArgs failed: %+v", err)
	}

	// 6. Determine device's WorkGroup size. This is probably how many items the GPU can process at a time.
	//local, err := kernel.WorkGroupSize(device)
	//if err != nil {
	//	logrus.Fatalf("WorkGroupSize failed: %+v", err)
	//}
	//logrus.Infof("Work group size: %d", local)
	//size, _ := kernel.PreferredWorkGroupSizeMultiple(nil)
	//logrus.Infof("Preferred Work Group Size Multiple: %d", size)

	// 6.1 calc local/global sizes. This stuff is passed on to the "enqueue". I think it's purpose is to handle
	//     cases where the data set size isn't divideable by the WG size
	//global := len(rays) / 8 // number of items to process, e.g. 32768
	//d := global % local     // given the preferred WG size, d is
	//logrus.Infof("Global: %d, D: %d", global, d)
	//if d != 0 {
	//	global += local - d
	//}
	//logrus.Infof("Global after applying D: %d, D: %d", global, d)

	// 7. Finally, start work! Enqueue executes the loaded args on the specified kernel.
	if _, err := queue.EnqueueNDRangeKernel(kernel, nil, []int{len(rays)}, nil, nil); err != nil {
		logrus.Fatalf("EnqueueNDRangeKernel failed: %+v", err)
	}

	// 8. Finish() blocks the main goroutine until the OpenCL queue is empty, i.e. all calculations are done
	if err := queue.Finish(); err != nil {
		logrus.Fatalf("Finish failed: %+v", err)
	}

	// 9. Allocate storage for loading the output from the OpenCL program, 4 float32 per cast ray. RGBA
	results := make([]float32, len(rays)*4)

	// 10. The EnqueueReadBuffer copies the data in the OpenCL "output" buffer into the "results" slice.
	dataPtrOut := unsafe.Pointer(&results[0])
	resSize := unsafe.Sizeof(results[0])
	sizePerEntry := int(resSize)
	dataSizeOut := sizePerEntry * len(results)

	if _, err := queue.EnqueueReadBuffer(output, true, 0, dataSizeOut, dataPtrOut, nil); err != nil {
		logrus.Fatalf("EnqueueReadBuffer failed: %+v", err)
	}

	queue.Flush()

	return results
}
