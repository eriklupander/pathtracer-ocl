package scenes

import (
	"bytes"
	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/camera"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

func ReferenceScene() func() *Scene {
	return func() *Scene {

		//cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0, 0.13, -0.9), geom.NewPoint(0, 0.02, -.1))
		cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0, 0.1, -1.5), geom.NewPoint(0, 0.05, 0))
		//cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0.01, .05, 0.01), geom.NewPoint(0,0,0))

		cam.FocalLength = cmd.Cfg.FocalLength
		cam.Aperture = cmd.Cfg.Aperture
		// left wall
		leftWall := shapes.NewPlane()
		leftWall.SetTransform(geom.Translate(-.6, 0, 0))
		leftWall.SetTransform(geom.RotateZ(math.Pi / 2))
		leftWall.SetTransform(geom.RotateY(math.Pi / 2))
		leftWall.SetMaterial(material.NewDiffuse(0.75, 0.25, 0.25))
		leftWall.Material.Textured = true
		leftWall.Material.TextureID = 0
		leftWall.Material.TextureScaleX = 1.0
		leftWall.Material.TextureScaleY = 1.0
		//
		//// right wall
		rightWall := shapes.NewPlane()
		rightWall.SetTransform(geom.Translate(.6, 0, 0))
		rightWall.SetTransform(geom.RotateZ(math.Pi / 2))
		rightWall.SetTransform(geom.RotateY(math.Pi / 2))
		rightWall.SetMaterial(material.NewDiffuse(0.25, 0.25, 0.75))
		rightWall.Material.Textured = true
		rightWall.Material.TextureID = 0
		rightWall.Material.TextureScaleX = 1.0
		rightWall.Material.TextureScaleY = 1.0

		// floor
		floor := shapes.NewPlane()
		floor.SetTransform(geom.Translate(0, -.4, 0))
		floorMaterial := material.NewDiffuse(0.9, 0.8, 0.7)
		floorMaterial.Textured = true
		floorMaterial.TextureID = 1
		floorMaterial.TextureScaleX = 0.25
		floorMaterial.TextureScaleY = 0.25
		floor.SetMaterial(floorMaterial)

		// ceiling
		ceil := shapes.NewPlane()
		ceil.SetTransform(geom.Translate(0, .4, 0))
		ceil.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))
		ceil.Material.Textured = true
		ceil.Material.TextureID = 2
		ceil.Material.TextureScaleX = 1.0
		ceil.Material.TextureScaleY = 1.0

		// back wall
		backWall := shapes.NewPlane()
		backWall.SetTransform(geom.Translate(0, 0, .4))
		backWall.SetTransform(geom.RotateX(math.Pi / 2))
		backWall.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))
		backWall.Material.Textured = true
		backWall.Material.TextureID = 0
		backWall.Material.TextureScaleX = 1.0
		backWall.Material.TextureScaleY = 1.0

		// front wall
		frontWall := shapes.NewPlane()
		frontWall.SetTransform(geom.Translate(0, 0, -2))
		frontWall.SetTransform(geom.RotateX(math.Pi / 2))
		frontWall.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// left sphere
		leftSphere := shapes.NewSphere()
		leftSphere.SetTransform(geom.Translate(-0.35, -0.28, -0.15))
		leftSphere.SetTransform(geom.Scale(0.12, 0.12, 0.12))
		leftSphere.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// middle sphere
		rightSphere := shapes.NewSphere()
		rightSphere.SetTransform(geom.Translate(0, -0.24, -0.30))
		rightSphere.SetTransform(geom.Scale(0.16, 0.16, 0.16))
		rightSphere.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// lightsource
		lightsource := shapes.NewSphere()
		lightsource.SetTransform(geom.Translate(0, .399, 0))
		lightsource.SetTransform(geom.Scale(0.283, 0.01, 0.283))

		light := material.NewLightBulb()
		//light.Emission = geom.NewColor(2.5, 2.5, 2.5)
		//light.Emission = geom.NewColor(7.5, 7.5, 7.5)
		light.Emission = geom.NewColor(9, 9, 9)
		lightsource.SetMaterial(light)

		shapes := []shapes.Shape{lightsource, floor, ceil, leftWall, rightWall, backWall, leftSphere, rightSphere}

		t0, err := os.ReadFile("./assets/concrete_squares.png")
		if err != nil {
			panic(err.Error())
		}
		img0, err := png.Decode(bytes.NewBuffer(t0))
		if err != nil {
			panic(err.Error())
		}

		t1, err := os.ReadFile("./assets/seamless-cobblestone-texture.jpg")
		if err != nil {
			panic(err.Error())
		}
		tmp, err := jpeg.Decode(bytes.NewBuffer(t1))
		if err != nil {
			panic(err.Error())
		}
		b := tmp.Bounds()
		img1 := image.NewNRGBA(image.Rect(0, 0, tmp.Bounds().Dx(), tmp.Bounds().Dy()))
		draw.Draw(img1, img1.Bounds(), tmp, b.Min, draw.Src)

		t2, err := os.ReadFile("./assets/floor_boards.png")
		if err != nil {
			panic(err.Error())
		}
		img2, err := png.Decode(bytes.NewBuffer(t2))
		if err != nil {
			panic(err.Error())
		}

		return &Scene{
			Camera:   cam,
			Objects:  shapes,
			Textures: []image.Image{img0, img1, img2},
		}
	}
}
