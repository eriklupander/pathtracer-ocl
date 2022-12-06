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

func TransparentTeapotScene() func() *Scene {
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
		leftSphere.SetTransform(geom.Translate(-0.25, -0.28, 0.25))
		leftSphere.SetTransform(geom.Scale(0.12, 0.12, 0.12))
		leftSphere.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// right sphere
		rightSphere := shapes.NewSphere()
		rightSphere.Label = "right_spr"
		rightSphere.SetTransform(geom.Translate(0.25, -0.28, 0.25))
		rightSphere.SetTransform(geom.Scale(0.12, 0.12, 0.12))
		rightSphere.SetMaterial(material.NewGlass())

		// left sphere
		teapot := teapot()
		teapot.Label = "teapot  "

		// lightsource
		lightsource := shapes.NewSphere()
		lightsource.Label = "light   "
		lightsource.SetTransform(geom.Translate(0, .399, 0))
		lightsource.SetTransform(geom.Scale(0.283, 0.01, 0.283))

		light := material.NewLightBulb()
		light.Emission = geom.NewColor(9, 9, 9)
		//light.Color = geom.NewColor(0.9, 0.8, 0.8)
		lightsource.SetMaterial(light)

		shapes := []shapes.Shape{lightsource, floor, ceil, leftWall, rightWall, backWall, leftSphere, rightSphere, teapot}

		return &Scene{
			Camera:  cam,
			Objects: shapes,
		}
	}
}

func teapot() *shapes.Group {
	data, err := os.ReadFile("assets/teapot.obj")
	if err != nil {
		panic(err.Error())
	}
	model := obj.ParseObj(string(data))
	group := model.ToGroup()

	// iterate over all triangles _before_ doing BVH divide to compute vertex normals since the teapot
	// model doesn't have pre-computed vertex models stored in the .obj file.
	tris := make([]*shapes.Triangle, 0)
	for i := range group.Children[0].(*shapes.Group).Children {
		tris = append(tris, group.Children[0].(*shapes.Group).Children[i].(*shapes.Triangle))
	}
	obj.ComputeVertexNormals(tris)

	group.Bounds()
	group.SetTransform(geom.Translate(0, -0.38, -0.2))
	group.SetTransform(geom.RotateY(math.Pi / 12))
	group.SetTransform(geom.Scale(0.1, 0.1, 0.1))
	group.SetMaterial(material.NewGlass())
	shapes.Divide(group, 50)
	group.Bounds()
	return group
}
