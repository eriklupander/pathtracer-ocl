package scenes

import (
	"fmt"
	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/camera"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"github.com/eriklupander/pathtracer-ocl/internal/app/obj"
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
	"io/ioutil"
	"math"
)

func ChristianScene() func() *Scene {
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

		// left sphere
		leftSphere := shapes.NewSphere()
		leftSphere.SetTransform(geom.Translate(-0.35, -0.28, -0.15))
		leftSphere.SetTransform(geom.Scale(0.12, 0.12, 0.12))
		leftSphere.SetMaterial(material.NewDiffuse(0.9, 0.9, 0.9))
		leftSphere.Material.Reflectivity = 0.99
		// cylinder
		cyl := shapes.NewCylinderMMC(0, 0.4, true)
		cyl.SetTransform(geom.Translate(0.45, -0.5, 0.2))
		//cyl.SetTransform(geom.RotateY(math.Pi / 4))
		//cyl.SetTransform(geom.RotateZ(math.Pi / 2))
		cyl.SetTransform(geom.Scale(0.075, 1, 0.075))
		cyl.SetMaterial(material.NewDiffuse(0.92, 0.4, 0.8))

		// cube
		cube := shapes.NewCube()
		cube.SetTransform(geom.Translate(0.1, -0.1, 0.1))
		cube.SetTransform(geom.Scale(0.1, 0.05, 0.04))
		cube.SetTransform(geom.RotateY(math.Pi / 4))
		cube.SetTransform(geom.RotateZ(math.Pi / 2))
		cube.SetMaterial(material.NewDiffuse(0.25, 0.25, 0.75))

		data, err := ioutil.ReadFile("assets/teapot.obj")
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
		group.SetTransform(geom.Translate(0, -0.4, 0))
		group.SetTransform(geom.Scale(0.07, 0.07, 0.07))
		silver := material.NewDiffuse(0.75, 0.75, 0.75)
		silver.Reflectivity = 0.2
		group.SetMaterial(silver)
		shapes.Divide(group, 50)
		group.Bounds()

		fmt.Printf("distance from camera to teapot: %f", geom.Magnitude(geom.Sub(geom.NewPoint(0, -0.4, 0.1), geom.NewPoint(0, 0.13, -0.9))))

		// lightsource

		lightMtl := material.NewLightBulb()
		lightMtl.Emission = geom.NewColor(90,80,60)

		lightsource := shapes.NewSphere()
		lightsource.SetTransform(geom.Translate(0, .3, 0))
		lightsource.SetTransform(geom.Scale(0.03, 0.03, 0.03))
		lightsource.SetMaterial(lightMtl)

		lightsource1 := shapes.NewSphere()
		lightsource1.SetTransform(geom.Translate(-0.5, .3, 0))
		lightsource1.SetTransform(geom.Scale(0.03, 0.03, 0.03))
		lightsource1.SetMaterial(lightMtl)

		lightsource2 := shapes.NewSphere()
		lightsource2.SetTransform(geom.Translate(-0.3, .3, 0))
		lightsource2.SetTransform(geom.Scale(0.03, 0.03, 0.03))
		lightsource2.SetMaterial(lightMtl)

		coverMtl := material.NewDiffuse(0.8, 0.8, 0.8)
		coverMtl.Reflectivity = 0.95
		cover2 := shapes.NewCylinderMMC(0, 1, false)
		cover2.SetTransform(geom.Translate(-0.3, .295, 0))
		cover2.SetTransform(geom.Scale(0.06, 0.4, 0.06))
		cover2.SetMaterial(coverMtl)

		lightsource3 := shapes.NewSphere()
		lightsource3.SetTransform(geom.Translate(-0.1, .3, 0))
		lightsource3.SetTransform(geom.Scale(0.03, 0.03, 0.03))
		lightsource3.SetMaterial(lightMtl)
		cover3 := shapes.NewCylinderMMC(0, 1, false)
		cover3.SetTransform(geom.Translate(-0.1, .295, 0))
		cover3.SetTransform(geom.Scale(0.06, 0.4, 0.06))
		cover3.SetMaterial(coverMtl)

		lightsource4 := shapes.NewSphere()
		lightsource4.SetTransform(geom.Translate(0.1, .3, 0))
		lightsource4.SetTransform(geom.Scale(0.03, 0.03, 0.03))
		lightsource4.SetMaterial(lightMtl)
		cover4 := shapes.NewCylinderMMC(0, 1, false)
		cover4.SetTransform(geom.Translate(0.1, .295, 0))
		cover4.SetTransform(geom.Scale(0.06, 0.4, 0.06))
		cover4.SetMaterial(coverMtl)

		lightsource5 := shapes.NewSphere()
		lightsource5.SetTransform(geom.Translate(0.3, .3, 0))
		lightsource5.SetTransform(geom.Scale(0.03, 0.03, 0.03))
		lightsource5.SetMaterial(lightMtl)
		cover5 := shapes.NewCylinderMMC(0, 1, false)
		cover5.SetTransform(geom.Translate(0.3, .295, 0))
		cover5.SetTransform(geom.Scale(0.06, 0.4, 0.06))
		cover5.SetMaterial(coverMtl)

		lightsource6 := shapes.NewSphere()
		lightsource6.SetTransform(geom.Translate(0.5, .3, 0))
		lightsource6.SetTransform(geom.Scale(0.03, 0.03, 0.03))
		lightsource6.SetMaterial(lightMtl)

		shapes := []shapes.Shape{ lightsource2, lightsource3, lightsource4, lightsource5, cover2, cover3, cover4, cover5,
			floor, ceil, leftWall, rightWall, backWall, group, leftSphere}

		return &Scene{
			Camera:  cam,
			Objects: shapes,
		}
	}
}
