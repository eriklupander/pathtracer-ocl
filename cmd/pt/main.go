package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"

	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/scenes"
	"github.com/eriklupander/pathtracer-ocl/internal/app/tracer"
	"github.com/jgillich/go-opencl/cl"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	var configFlags = pflag.NewFlagSet("config", pflag.ExitOnError)
	configFlags.Int("width", 640, "Image width")
	configFlags.Int("height", 480, "Image height")
	configFlags.Int("samples", 1, "Number of samples per pixel")
	configFlags.Float64("aperture", 0.0, "Aperture. If 0, no DoF will be used")
	configFlags.Float64("focal-length", 0.0, "Focal length.")
	configFlags.String("scene", "gopher", "scene from /scenes")
	configFlags.Int("device-index", 0, "Use device with index (use --list-devices to list available devices)")
	configFlags.Bool("list-devices", false, "List available devices")

	if err := configFlags.Parse(os.Args[1:]); err != nil {
		panic(err.Error())
	}
	if err := viper.BindPFlags(configFlags); err != nil {
		panic(err.Error())
	}
	viper.AutomaticEnv()

	cmd.FromConfig()

	if cmd.Cfg.ListDevices {
		listDevices()
		return
	}

	var scene func() *scenes.Scene
	switch cmd.Cfg.Scene {
	case "reference":
		scene = scenes.ReferenceScene()
	case "teapot":
		scene = scenes.ModelScene()
	case "gopher":
		scene = scenes.GopherScene()
	case "gopher-window":
		scene = scenes.GopherWindowScene()
	case "christian":
		scene = scenes.ChristianScene()
	case "textures":
		scene = scenes.TexturedPlanetsScene()
	case "envmap":
		scene = scenes.EnvironmentMap()
	case "cubemap":
		scene = scenes.EnvironmentCubeMap()
	case "reflection":
		scene = scenes.ReflectionsScene()
	case "transparency":
		scene = scenes.TransparencyScene()
	case "transparent_teapot":
		scene = scenes.TransparentTeapotScene()
	default:
		scene = scenes.OCLScene()
	}

	tracer.Render(scene)
}

func listDevices() {
	platforms, err := cl.GetPlatforms()
	if err != nil {
		logrus.Fatalf("Failed to get platforms: %+v", err)
	}
	platform := platforms[0]

	devices, err := platform.GetDevices(cl.DeviceTypeAll)
	if err != nil {
		logrus.Fatalf("Failed to get devices: %+v", err)
	}
	for idx, device := range devices {
		fmt.Printf("Index: %d Type: %s Name: %s\n", idx, device.Type(), device.Name())
	}
}
