package scenes

import (
	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/camera"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"github.com/eriklupander/pathtracer-ocl/internal/app/obj"
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
	"io/ioutil"
	"math"
)

func ModelScene() func() *Scene {
	return func() *Scene {

		cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0, 0.13, -1), geom.NewPoint(0, 0.02, 0))
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
		backWall.SetTransform(geom.Translate(0, 0, .4))
		backWall.SetTransform(geom.RotateX(math.Pi / 2))
		backWall.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// front wall
		frontWall := shapes.NewPlane()
		frontWall.SetTransform(geom.Translate(0, 0, -2))
		frontWall.SetTransform(geom.RotateX(math.Pi / 2))
		frontWall.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		data, err := ioutil.ReadFile("assets/teapot.obj")
		if err != nil {
			panic(err.Error())
		}
		model := obj.ParseObj(string(data))
		group := model.ToGroup().Children[0].(*shapes.Group)
		tris := make([]*shapes.Triangle, 0)
		for i := range group.Children {
			tris = append(tris, group.Children[i].(*shapes.Triangle))
		}
		obj.ComputeVertexNormals(tris)
		group.Bounds()
		group.SetTransform(geom.Translate(0, -0.4, 0.1))
		group.SetTransform(geom.Scale(0.07, 0.07, 0.07))
		silver := material.NewDiffuse(0.75, 0.75, 0.75)
		//silver.Reflectivity = 0.2
		group.SetMaterial(silver)

		// lightsource
		lightsource := shapes.NewSphere()
		lightsource.SetTransform(geom.Translate(0, 1.36, 0))
		light := material.NewLightBulb()
		light.Emission = geom.NewColor(9, 8, 6)
		lightsource.SetMaterial(light)

		return &Scene{
			Camera:  cam,
			Objects: []shapes.Shape{floor, ceil, leftWall, rightWall, backWall, lightsource, group},
		}
	}
}
