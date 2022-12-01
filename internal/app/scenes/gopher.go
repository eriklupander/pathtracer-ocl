package scenes

import (
	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/camera"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"github.com/eriklupander/pathtracer-ocl/internal/app/obj"
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
	"math"
	"os"
)

func GopherScene() func() *Scene {
	return func() *Scene {

		cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0, 0.1, -1.5), geom.NewPoint(0, 0.05, 0))
		cam.FocalLength = cmd.Cfg.FocalLength
		cam.Aperture = cmd.Cfg.Aperture
		// left wall
		leftWall := shapes.NewPlane()
		leftWall.SetTransform(geom.Translate(-.6, 0, 0))
		leftWall.SetTransform(geom.RotateZ(math.Pi / 2))
		leftWall.SetMaterial(material.NewDiffuse(0.75, 0.25, 0.25))
		//
		//// right wall
		rightWall := shapes.NewPlane()
		rightWall.SetTransform(geom.Translate(.6, 0, 0))
		rightWall.SetTransform(geom.RotateZ(math.Pi / 2))
		rightWall.SetMaterial(material.NewDiffuse(0.25, 0.25, 0.75))

		// floor
		floor := shapes.NewPlane()
		floor.SetTransform(geom.Translate(0, -.4, 0))
		floor.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// ceiling
		ceil := shapes.NewPlane()
		ceil.SetTransform(geom.Translate(0, .4, 0))
		ceil.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// back wall
		backWall := shapes.NewPlane()
		backWall.SetTransform(geom.Translate(0, 0, 1.4))
		backWall.SetTransform(geom.RotateX(math.Pi / 2))
		backWall.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// front wall
		frontWall := shapes.NewPlane()
		frontWall.SetTransform(geom.Translate(0, 0, -2))
		frontWall.SetTransform(geom.RotateX(math.Pi / 2))
		frontWall.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// right sphere
		rightSphere := shapes.NewSphere()
		rightSphere.SetTransform(geom.Translate(0.28, -0.24, 0.15))
		rightSphere.SetTransform(geom.Scale(0.16, 0.16, 0.16))
		//rightSphere.SetMaterial(material.NewDiffuse(0.57, 0.86, 1))
		halfMirror := material.NewMirror()
		halfMirror.Reflectivity = 0.8
		halfMirror.Color = geom.NewColor(0.97, 0.97, 0.843)
		rightSphere.SetMaterial(halfMirror)

		objects := []shapes.Shape{floor, ceil, leftWall, rightWall, backWall, frontWall, rightSphere}

		data, err := os.ReadFile("assets/gopher.obj")
		if err != nil {
			panic(err.Error())
		}
		model := obj.ParseObj(string(data))
		group := model.ToGroup()
		group.Bounds()

		group.SetTransform(geom.Translate(-.4, -0.15, 0.2))
		group.SetTransform(geom.RotateZ(-math.Pi / 2))
		group.SetTransform(geom.RotateX(-math.Pi / 4))
		group.SetTransform(geom.Scale(0.2, 0.2, 0.2))
		silver := material.NewDiffuse(0.75, 0.75, 0.75)
		silver.Reflectivity = 0.2
		group.SetMaterial(silver)
		shapes.Divide(group, 60)
		group.Bounds()
		objects = append(objects, group)

		// lightsource
		lightsource := shapes.NewSphere()
		lightsource.SetTransform(geom.Translate(0, 1.36, 0))
		light := material.NewLightBulb()
		light.Emission = geom.NewColor(9, 8, 6)
		lightsource.SetMaterial(light)
		objects = append(objects, lightsource)

		return &Scene{
			Camera:  cam,
			Objects: objects,
		}
	}
}
