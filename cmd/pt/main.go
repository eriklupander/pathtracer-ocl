package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/scenes"
	"github.com/eriklupander/pathtracer-ocl/internal/app/tracer"
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
	configFlags.String("scene", "reference", "scene from /scenes")

	if err := configFlags.Parse(os.Args[1:]); err != nil {
		panic(err.Error())
	}
	if err := viper.BindPFlags(configFlags); err != nil {
		panic(err.Error())
	}
	viper.AutomaticEnv()

	cmd.FromConfig()

	var scene = scenes.OCLScene()

	tracer.Render(scene)
}
