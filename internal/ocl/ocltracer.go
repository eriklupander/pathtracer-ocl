package ocl

import (
	_ "embed"
	"math/rand"
	"time"
	"unsafe"

	"github.com/jgillich/go-opencl/cl"
	"github.com/sirupsen/logrus"
)

// BuildSceneBuffer maps shapes to a float64 slice:
// Transform:        4x4 float64, offset: 0
// Inverse:          4x4 float64, offset: 16
// InverseTranspose: 4x4 float64, offset: 32
// Color:            4xfloat64, offset: 48
// Emission:         4xfloat64, offset: 52
// RefractiveIndex:  1xfloat64, offset: 56
// Type:             1xInt64, offset: 57

//go:embed tracer.cl
var kernelSource string

type CLRay struct {
	Origin    [4]float64
	Direction [4]float64
}

type CLObject struct {
	Transform        [16]float64
	Inverse          [16]float64
	InverseTranspose [16]float64
	Color            [4]float64
	Emission         [4]float64
	RefractiveIndex  float64
	Type             int64
	MinY             float64
	MaxY             float64
	Reflectivity     float64
	Padding2         int64
	Padding3         int64
	Padding4         int64
}

type CLCamera struct {
	Width       int32   // 4
	Height      int32   // 8
	Fov         float64 // 16
	PixelSize   float64 // 24
	HalfWidth   float64 // 32
	HalfHeight  float64 // 40
	Aperture    float64 // 48
	FocalLength float64 // 56
	//Transform   [16]float64
	Inverse [16]float64 // 128 + 56 == 184
	Padding [72]byte    // 256-72 == 184
}

// Trace is the entry point for transforming input data into their OpenCL representations, setting up boilerplate
// and calling the entry kernel. Should return a slice of float64 RGBA RGBA RGBA once finished.
func Trace(objects []CLObject, width, height, samples int, camera CLCamera) []float64 {
	numPixels := int(camera.Width * camera.Height)
	logrus.Infof("trace with %d objects %dx%d", len(objects), camera.Width, camera.Height)

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
	// Use the "highest" device index, is usually the discrete GPU
	deviceIndex := 0 //len(devices) - 1

	if deviceIndex < 0 {
		deviceIndex = 0
	}
	device := devices[deviceIndex] // 0 == CPU 1 == iGPU 2 == GPU
	logrus.Infof("Using device %d %v", deviceIndex, devices[deviceIndex].Name())

	// 1. Select a device to use.
	//    On my mac           : 0 == CPU, 1 == Iris GPU, 2 == GeForce 750M GPU
	//    On my windows AMD PC: 0 == Gefore RTX2080
	//    Use selected device to create an OpenCL context
	context, err := cl.CreateContext([]*cl.Device{device})
	if err != nil {
		logrus.Fatalf("CreateContext failed: %+v", err)
	}

	// 2. Create a "Command Queue" bound to the selected device
	queue, err := context.CreateCommandQueue(device, 0)
	if err != nil {
		logrus.Fatalf("CreateCommandQueue failed: %+v", err)
	}

	// 3.0 Read kernel source from embedded .cl file and
	//     create an OpenCL "program" from the source code.
	program, err := context.CreateProgramWithSource([]string{kernelSource})
	if err != nil {
		logrus.Fatalf("CreateProgramWithSource failed: %+v", err)
	}

	// 3.2 Build the OpenCL program
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
		_, err := kernel.ArgName(i)
		if err == cl.ErrUnsupported {
			logrus.Errorf("GetKernelArgInfo for arg: %d ErrUnsupported", i)
			continue
		} else if err != nil {
			logrus.Errorf("GetKernelArgInfo for arg: %d failed: %+v", i, err)
			continue
		}
	}

	// 5. Determine device's WorkGroup size. This is probably how many items the GPU can process at a time.
	workGroupSize, err := kernel.WorkGroupSize(device)
	if err != nil {
		logrus.Fatalf("WorkGroupSize failed: %+v", err)
	}
	logrus.Infof("Work group size: %d", workGroupSize)

	// Make sure the WGS is never greater than the total number of items we're going to process
	if workGroupSize > numPixels {
		workGroupSize = numPixels
	}
	if numPixels%workGroupSize != 0 {
		logrus.Fatal("The number of rays must be a power of the WorkGroupSize")
	}

	// split work into batches in order to avoid kernels running for more than 10 seconds
	// otherwise, the GPU driver will kill us.
	results := make([]float64, 0)
	batchSize := 8
	if batchSize > numPixels {
		batchSize = numPixels
	}
	for y := 0; y < height; y += batchSize {
		st := time.Now()
		results = append(results, computeBatch(objects, camera, context, kernel, queue, samples, workGroupSize, y, batchSize)...)
		logrus.Infof("%d/%d lines done in %v", y+batchSize, height, time.Since(st))
	}

	return results
}

func computeBatch(objects []CLObject, camera CLCamera, context *cl.Context, kernel *cl.Kernel, queue *cl.CommandQueue, samples, workGroupSize, rowOffset, rowsPerBatch int) []float64 {
	pixelsInBatch := rowsPerBatch * int(camera.Width)

	// populate seed of random numbers, OpenCL can't do random by itself AFAIK
	seed := make([]float64, pixelsInBatch)
	for i := 0; i < pixelsInBatch; i++ {
		seed[i] = rand.Float64()
	}

	// 5. Time to start loading data into GPU memory

	// 5.1 create OpenCL buffers (memory) for the pre-computed rays and scene objects.
	// Note that we're allocating 64 bytes per ray (8xfloat64) and 512 bytes per scene object.
	// Remember - each float64 uses 8 bytes.

	objectsBuffer, err := context.CreateEmptyBuffer(cl.MemReadOnly, 512*len(objects))
	if err != nil {
		logrus.Fatalf("CreateBuffer failed for objects input: %+v", err)
	}
	defer objectsBuffer.Release()

	seedBuffer, err := context.CreateEmptyBuffer(cl.MemReadOnly, 8*len(seed))
	if err != nil {
		logrus.Fatalf("CreateBuffer failed for seedBuffer input: %+v", err)
	}
	defer seedBuffer.Release()

	cameraBuffer, err := context.CreateEmptyBuffer(cl.MemReadOnly, 256)
	if err != nil {
		logrus.Fatalf("CreateBuffer failed for camera input: %+v", err)
	}
	defer cameraBuffer.Release()

	// 5.2 create OpenCL buffer (memory) for the output data, we want RGBA per ray, i.e. 4 float64 per ray.
	// So, we'll need 32 bytes to store the final computed color for each ray. Remember, we pass 1 ray per pixel.
	output, err := context.CreateEmptyBuffer(cl.MemReadOnly, pixelsInBatch*32)
	if err != nil {
		logrus.Fatalf("CreateBuffer failed for output: %+v", err)
	}
	defer output.Release()

	// 5.3 This is where we connect our input to the command queue and upload the actual data into GPU memory
	//     The rayDataPtr is a pointer to the first element of the rays slice,
	//     while rayDataTotalSizeBytes is the total length of the ray data, in bytes - i.e len(rays) * 64.
	objectsDataPtr := unsafe.Pointer(&objects[0])
	objectsDataSize := int(unsafe.Sizeof(objects[0])) * len(objects)
	if _, err := queue.EnqueueWriteBuffer(objectsBuffer, true, 0, objectsDataSize, objectsDataPtr, nil); err != nil {
		logrus.Fatalf("EnqueueWriteBuffer failed: %+v", err)
	}

	seedDataPtr := unsafe.Pointer(&seed[0])
	seedDataSize := int(unsafe.Sizeof(seed[0])) * len(seed)
	if _, err := queue.EnqueueWriteBuffer(seedBuffer, true, 0, seedDataSize, seedDataPtr, nil); err != nil {
		logrus.Fatalf("EnqueueWriteBuffer failed: %+v", err)
	}

	cameraDataPtr := unsafe.Pointer(&camera)
	cameraDataSize := int(unsafe.Sizeof(camera))
	if _, err := queue.EnqueueWriteBuffer(cameraBuffer, true, 0, cameraDataSize, cameraDataPtr, nil); err != nil {
		logrus.Fatalf("EnqueueWriteBuffer failed: %+v", err)
	}

	// 5.4 Kernel is our program and here we explicitly bind our 4 parameters to it
	if err := kernel.SetArgs(objectsBuffer, uint32(len(objects)), output, seedBuffer, uint32(samples), cameraBuffer, uint32(rowOffset)); err != nil {
		logrus.Fatalf("SetKernelArgs failed: %+v", err)
	}

	// 7. Finally, start work! Enqueue executes the loaded args on the specified kernel.
	if _, err := queue.EnqueueNDRangeKernel(kernel, nil, []int{pixelsInBatch}, []int{workGroupSize}, nil); err != nil {
		logrus.Fatalf("EnqueueNDRangeKernel failed: %+v", err)
	}

	// 8. Finish() blocks the main goroutine until the OpenCL queue is empty, i.e. all calculations are done
	if err := queue.Finish(); err != nil {
		logrus.Fatalf("Finish failed: %+v", err)
	}

	// 9. Allocate storage for loading the output from the OpenCL program, 4 float64 per cast ray. RGBA
	results := make([]float64, pixelsInBatch*4)

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
