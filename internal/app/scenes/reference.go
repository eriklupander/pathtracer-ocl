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
		leftWall.SetTransform(geom.RotateX(math.Pi))
		leftWall.SetTransform(geom.RotateZ(math.Pi / 2)) // note rotate by -2 so the "correct" side faces inwards.
		leftWall.SetTransform(geom.RotateY(math.Pi / 2))
		leftWall.SetMaterial(material.NewDiffuse(0.75, 0.25, 0.25))
		leftWall.Material.Textured = true
		leftWall.Material.TextureID = 0
		leftWall.Material.TextureScaleX = 1.0
		leftWall.Material.TextureScaleY = 1.0
		leftWall.Material.TexturedNM = true
		leftWall.Material.TextureIDNM = 3
		leftWall.Material.TextureScaleXNM = 1.0
		leftWall.Material.TextureScaleYNM = 1.0

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
		rightWall.Material.TexturedNM = true
		rightWall.Material.TextureIDNM = 3
		rightWall.Material.TextureScaleXNM = 1.0
		rightWall.Material.TextureScaleYNM = 1.0

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
		backWall.Material.TexturedNM = true
		backWall.Material.TextureIDNM = 3
		backWall.Material.TextureScaleXNM = 1.0
		backWall.Material.TextureScaleYNM = 1.0

		// front wall
		frontWall := shapes.NewPlane()
		frontWall.SetTransform(geom.Translate(0, 0, -2))
		frontWall.SetTransform(geom.RotateX(math.Pi / 2))
		frontWall.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// left sphere
		leftSphere := shapes.NewSphere()
		leftSphere.SetTransform(geom.Translate(-0.3, -0.1, -0.25))
		leftSphere.SetTransform(geom.Scale(0.2, 0.2, 0.2))
		leftSphere.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))
		leftSphere.Material.Textured = true
		leftSphere.Material.TextureID = 1

		// middle sphere
		rightSphere := shapes.NewSphere()
		rightSphere.SetTransform(geom.Translate(0.2, 0, -0.3)) // var -.24
		rightSphere.SetTransform(geom.RotateY(math.Pi))
		rightSphere.SetTransform(geom.Scale(0.25, 0.25, 0.25))
		rightSphere.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))
		rightSphere.Material.Textured = true
		rightSphere.Material.TextureID = 0

		// lightsource
		lightsource := shapes.NewSphere()
		lightsource.SetTransform(geom.Translate(0, .395, -.9))
		lightsource.SetTransform(geom.Scale(0.283, 0.01, 0.283))

		light := material.NewLightBulb()
		light.Emission = geom.NewColor(10, 10, 10)
		lightsource.SetMaterial(light)

		lightsource2 := shapes.NewSphere()
		lightsource2.SetTransform(geom.Translate(0, 0, -1.7))
		lightsource2.SetTransform(geom.Scale(0.283, 0.283, 0.01))
		lightsource2.SetMaterial(light)

		shapes := []shapes.Shape{lightsource, lightsource2, floor, ceil, leftWall, rightWall, backWall, leftSphere, rightSphere}

		squares := LoadImage("./assets/concrete_squares.png")
		squaresNormalMap := LoadImage("./assets/concrete_squares_nm2.png")
		cobbleStones := LoadImage("./assets/seamless-cobblestone-texture.jpg")
		floorBoards := LoadImage("./assets/floor_boards.png")
		planet := LoadImage("./assets/planet.png")
		jupiter := LoadImage("./assets/jupiter2_6k_contrast.png")
		return &Scene{
			Camera:         cam,
			Objects:        shapes,
			Textures:       []image.Image{squares, cobbleStones, floorBoards, squaresNormalMap},
			SphereTextures: []image.Image{planet, jupiter},
		}
	}
}
