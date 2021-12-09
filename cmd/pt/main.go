package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/scenes"
	"github.com/eriklupander/pathtracer-ocl/internal/app/tracer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	var configFlags = pflag.NewFlagSet("config", pflag.ExitOnError)
	configFlags.Int("workers", runtime.NumCPU(), "number of workers")
	configFlags.Int("width", 1280, "Image width")
	configFlags.Int("height", 960, "Image height")
	configFlags.Int("samples", 1, "Number of samples per pixel")
	configFlags.String("scene", "reference", "scene from /scenes")

	if err := configFlags.Parse(os.Args[1:]); err != nil {
		panic(err.Error())
	}
	if err := viper.BindPFlags(configFlags); err != nil {
		panic(err.Error())
	}
	viper.AutomaticEnv()

	cmd.FromConfig()
	logrus.Printf("Running with %d CPUs\n", viper.GetInt("workers"))

	var scene func() *scenes.Scene
	switch viper.GetString("scene") {
	case "reference":
		scene = scenes.OCLScene()
	case "cornell":
		scene = scenes.OCLScene()
	default:
		scene = scenes.OCLScene()
	}

	tracer.Render(scene)
}
