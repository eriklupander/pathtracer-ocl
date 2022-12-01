package scenes

import (
	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/camera"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"github.com/eriklupander/pathtracer-ocl/internal/app/obj"
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
	"image"
	"math"
	"os"
)

func EnvironmentCubeMap() func() *Scene {
	return func() *Scene {

		//cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0, 0.13, -0.9), geom.NewPoint(0, 0.02, -.1))
		cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0, 0.3, -2.7), geom.NewPoint(0, 0.45, 0))
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
		rightSphere.SetTransform(geom.Translate(.2, 1, 2))
		rightSphere.SetTransform(geom.Scale(0.26, 0.26, 0.26))
		rightSphere.SetMaterial(material.NewMirror()) //0.9, 0.8, 0.7))
		//rightSphere.Material.Reflectivity = 0.95

		// lightsource
		lightsource := shapes.NewSphere()
		lightsource.SetTransform(geom.Translate(1.1, 1, -4))
		lightsource.SetTransform(geom.Scale(0.7, 0.7, 0.7))
		light := material.NewLightBulb()
		light.Emission = geom.NewColor(19.5, 19.5, 19.5)
		lightsource.SetMaterial(light)

		// IMPORTANT! The decoded image gets much darker than original
		envTexture := LoadImage("./assets/shrine_cubemap.jpeg")
		skySphere := shapes.NewCube()
		skySphere.SetTransform(geom.Translate(0, 0, 0))
		//skySphere.SetTransform(geom.RotateY(-math.Pi))
		skySphere.SetTransform(geom.Scale(5, 5, 5))
		skySphere.Material = material.NewDefaultMaterial()
		skySphere.Material.Textured = true
		skySphere.Material.TextureID = 0
		skySphere.Material.TextureScaleX = 1.0
		skySphere.Material.TextureScaleY = 1.0
		skySphere.Material.Emission = geom.NewColor(1, 1, 1)
		skySphere.Material.IsEnvMap = true

		data, err := os.ReadFile("assets/gopher.obj")
		if err != nil {
			panic(err.Error())
		}
		model := obj.ParseObj(string(data))
		group := model.ToGroup()
		group.Bounds()

		group.SetTransform(geom.Translate(-.7, -0.15, 0.2))
		group.SetTransform(geom.RotateZ(-math.Pi / 2))
		group.SetTransform(geom.RotateX(-math.Pi / 4))
		group.SetTransform(geom.Scale(0.4, 0.4, 0.4))
		silver := material.NewDiffuse(0.75, 0.75, 0.75)
		silver.Reflectivity = 0.0
		group.SetMaterial(silver)
		shapes.Divide(group, 60)
		group.Bounds()

		shapes := []shapes.Shape{lightsource, rightSphere, skySphere, group}

		return &Scene{
			Camera:       cam,
			Objects:      shapes,
			CubeTextures: []image.Image{envTexture},
		}
	}
}
