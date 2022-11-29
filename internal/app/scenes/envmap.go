package scenes

import (
	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/camera"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
	"image"
	"math"
)

func EnvironmentMap() func() *Scene {
	return func() *Scene {

		//cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0, 0.13, -0.9), geom.NewPoint(0, 0.02, -.1))
		cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0, 0.1, -1.5), geom.NewPoint(0, 0.15, 0))
		//cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0.01, .05, 0.01), geom.NewPoint(0,0,0))

		cam.FocalLength = cmd.Cfg.FocalLength
		cam.Aperture = cmd.Cfg.Aperture

		// cube
		cube := shapes.NewCube()
		cube.SetTransform(geom.Translate(0.1, -0.1, 0.1))
		cube.SetTransform(geom.Scale(0.1, 0.05, 0.04))
		cube.SetTransform(geom.RotateY(math.Pi / 4))
		cube.SetTransform(geom.RotateZ(math.Pi / 2))
		cube.SetMaterial(material.NewDiffuse(0.25, 0.25, 0.75))

		// left sphere
		leftSphere := shapes.NewSphere()
		leftSphere.SetTransform(geom.Translate(-0.35, -0.28, -0.15))
		leftSphere.SetTransform(geom.Scale(0.12, 0.12, 0.12))
		leftSphere.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// middle sphere
		rightSphere := shapes.NewSphere()
		rightSphere.SetTransform(geom.Translate(0, -0.14, -0.30))
		rightSphere.SetTransform(geom.Scale(0.16, 0.16, 0.16))
		rightSphere.SetMaterial(material.NewMirror()) //0.9, 0.8, 0.7))
		//rightSphere.Material.Reflectivity = 0.95

		// lightsource
		lightsource := shapes.NewSphere()
		lightsource.SetTransform(geom.Translate(0, .399, 0))
		lightsource.SetTransform(geom.Scale(0.283, 0.01, 0.283))

		light := material.NewLightBulb()
		light.Emission = geom.NewColor(2.5, 2.5, 2.5)
		lightsource.SetMaterial(light)

		// IMPORTANT! The decoded image gets much darker than original
		envTexture := LoadImage("./assets/alps_field_8k.png")
		skySphere := shapes.NewSphere()
		skySphere.SetTransform(geom.Scale(5, 5, 5))
		skySphere.Material = material.NewDefaultMaterial()
		skySphere.Material.Textured = true
		skySphere.Material.TextureID = 0
		skySphere.Material.TextureScaleX = 1.0
		skySphere.Material.TextureScaleY = 1.0
		skySphere.Material.Emission = geom.NewColor(1, 1, 1)

		shapes := []shapes.Shape{rightSphere, skySphere}

		return &Scene{
			Camera:         cam,
			Objects:        shapes,
			SphereTextures: []image.Image{envTexture},
		}
	}
}
