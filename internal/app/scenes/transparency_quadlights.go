package scenes

import (
	"fmt"
	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/camera"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
	"math"
)

func TransparencyQuadLightsScene() func() *Scene {
	return func() *Scene {

		//cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0, 0.13, -0.9), geom.NewPoint(0, 0.02, -.1))
		cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0, 0.1, -1.5), geom.NewPoint(0, 0.05, 0))
		//cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0.01, .05, 0.01), geom.NewPoint(0,0,0))

		cam.FocalLength = cmd.Cfg.FocalLength
		cam.Aperture = cmd.Cfg.Aperture
		// left wall
		leftWall := shapes.NewPlane()
		leftWall.Label = "leftwall"
		leftWall.SetTransform(geom.Translate(-.6, 0, 0))
		leftWall.SetTransform(geom.RotateZ(math.Pi / 2))
		leftWall.SetMaterial(material.NewDiffuse(0.75, 0.25, 0.25))
		//
		//// right wall
		rightWall := shapes.NewPlane()
		rightWall.Label = "rghtwall"
		rightWall.SetTransform(geom.Translate(.6, 0, 0))
		rightWall.SetTransform(geom.RotateZ(math.Pi / 2))
		rightWall.SetMaterial(material.NewDiffuse(0.25, 0.25, 0.75))

		// floor
		floor := shapes.NewPlane()
		floor.Label = "floor   "
		floor.SetTransform(geom.Translate(0, -.4, 0))
		floor.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// ceiling
		ceil := shapes.NewPlane()
		ceil.Label = "ceiling "
		ceil.SetTransform(geom.Translate(0, .4, 0))
		ceil.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// back wall
		backWall := shapes.NewPlane()
		backWall.Label = "backwall"
		backWall.SetTransform(geom.Translate(0, 0, .6))
		backWall.SetTransform(geom.RotateX(math.Pi / 2))
		backWall.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// front wall
		frontWall := shapes.NewPlane()
		frontWall.Label = "frntwall"
		frontWall.SetTransform(geom.Translate(0, 0, -2))
		frontWall.SetTransform(geom.RotateX(math.Pi / 2))
		frontWall.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// left sphere
		leftSphere := shapes.NewSphere()
		leftSphere.Label = "left_spr"
		leftSphere.SetTransform(geom.Translate(-0.25, -0.18, 0.25))
		leftSphere.SetTransform(geom.Scale(0.14, 0.14, 0.14))
		leftSphere.SetMaterial(material.NewGlass()) //material.NewDiffuse(0.9, 0.8, 0.7))

		// middle sphere
		middleSphere := shapes.NewSphere()
		middleSphere.Label = "mddl_spr"
		middleSphere.SetTransform(geom.Translate(0, -0.24, -0.30))
		middleSphere.SetTransform(geom.Scale(0.16, 0.16, 0.16))
		middleSphere.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))
		middleSphere.Material.RefractiveIndex = 1.57

		// right sphere
		rightSphere := shapes.NewSphere()
		rightSphere.Label = "right_spr"
		rightSphere.SetTransform(geom.Translate(0.35, -0.23, 0.2))
		rightSphere.SetTransform(geom.Scale(0.17, 0.17, 0.17))
		rightSphere.SetMaterial(material.NewMirror())

		// lightsources
		lights := make([]shapes.Shape, 0)
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				light := shapes.NewCube()
				light.Label = fmt.Sprintf("light %d-%d", i, j)
				light.SetTransform(geom.Translate(-0.25+float64(i)*0.5, .399, -0.25+float64(j)*0.5))
				light.SetTransform(geom.Scale(0.15, 0.01, 0.15))
				light.SetMaterial(material.NewLightBulb())
				light.Material.Emission = geom.NewColor(9, 9, 9)
				light.Material.Color = geom.NewColor(1, 1, 1)
				lights = append(lights, light)
			}
		}

		shapes := []shapes.Shape{floor, ceil, leftWall, rightWall, backWall, leftSphere, middleSphere, rightSphere}
		shapes = append(shapes, lights...)
		return &Scene{
			Camera:  cam,
			Objects: shapes,
		}
	}
}
