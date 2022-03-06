package cmd

import "github.com/spf13/viper"

type Config struct {
	Width       int
	Height      int
	Workers     int
	Samples     int
	Aperture    float64
	FocalLength float64
	DeviceIndex int
	ListDevices bool
}

var Cfg *Config

func FromConfig() {
	Cfg = &Config{
		Width:       viper.GetInt("width"),
		Height:      viper.GetInt("height"),
		Samples:     viper.GetInt("samples"),
		Aperture:    viper.GetFloat64("aperture"),
		FocalLength: viper.GetFloat64("focal-length"),
		DeviceIndex: viper.GetInt("device-index"),
		ListDevices: viper.GetBool("list-devices"),
	}
}
