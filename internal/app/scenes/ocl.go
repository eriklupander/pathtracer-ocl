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

func OCLScene() func() *Scene {
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
		leftSphere.SetTransform(geom.Translate(-0.25, -0.24, 0.1))
		leftSphere.SetTransform(geom.Scale(0.16, 0.16, 0.16))
		leftSphere.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))
		//leftSphere.SetMaterial(material.NewMirror())

		// middle sphere
		middleSphere := shapes.NewSphere()
		middleSphere.SetTransform(geom.Translate(0, -0.24, -0.30))
		middleSphere.SetTransform(geom.Scale(0.16, 0.16, 0.16))
		//		middleSphere.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))
		middleSphere.SetMaterial(material.NewGlass())

		// right sphere
		rightSphere := shapes.NewSphere()
		rightSphere.SetTransform(geom.Translate(0.25, -0.24, 0.1))
		rightSphere.SetTransform(geom.Scale(0.16, 0.16, 0.16))
		//rightSphere.SetMaterial(material.NewDiffuse(0.57, 0.86, 1))
		halfMirror := material.NewMirror()
		halfMirror.Reflectivity = 0.8
		halfMirror.Color = geom.NewColor(0.97, 0.97, 0.843)
		rightSphere.SetMaterial(halfMirror)

		// cylinder
		cyl := shapes.NewCylinderMMC(0, 0.4, true)
		cyl.SetTransform(geom.Translate(0.45, -0.5, -0.2))
		//cyl.SetTransform(geom.RotateY(math.Pi / 4))
		//cyl.SetTransform(geom.RotateZ(math.Pi / 2))
		cyl.SetTransform(geom.Scale(0.075, 1, 0.075))
		cyl.SetMaterial(material.NewDiffuse(0.92, 0.4, 0.8))

		// cube
		cube := shapes.NewCube()
		cube.SetTransform(geom.Translate(-0.3, -0.375, -0.3))
		cube.SetTransform(geom.Scale(0.1, 0.05, 0.04))
		cube.SetTransform(geom.RotateY(math.Pi / 4))
		cube.SetTransform(geom.RotateZ(math.Pi / 2))
		cube.SetMaterial(material.NewDiffuse(0.25, 0.25, 0.75))

		// lightsource
		lightsource := shapes.NewSphere()
		lightsource.SetTransform(geom.Translate(0, 1.36, 0))
		light := material.NewLightBulb()
		light.Emission = geom.NewColor(9, 8, 6)
		lightsource.SetMaterial(light)

		// add two triangles to a group
		tri1 := shapes.NewTriangleN(geom.NewPoint(-0.2, -.4, 0), geom.NewPoint(0.0, -.4, 0), geom.NewPoint(0, -0.1, 0))
		tri2 := shapes.NewTriangleN(geom.NewPoint(0, -.4, 0), geom.NewPoint(0.2, -.4, 0), geom.NewPoint(0, -0.1, 0))
		tri3 := shapes.NewTriangleN(geom.NewPoint(0.1, -.4, -0.4), geom.NewPoint(0, -0.1, 0), geom.NewPoint(0, -.4, 0))

		group := shapes.NewGroup()
		group.SetMaterial(material.NewDiffuse(0.7, 0.4, 0.9))
		group.SetTransform(geom.Translate(0.15, 0, -0.25))
		//group.SetTransform(geom.RotateY(math.Pi / 4))
		group.AddChildren(tri1, tri2, tri3)
		group.Bounds()
		fmt.Printf("%+v\n", group.BoundingBox)
		return &Scene{
			Camera:  cam,
			Objects: []shapes.Shape{floor, ceil, leftWall, rightWall, backWall, leftSphere, rightSphere, cyl, cube, group, lightsource},
		}
	}
}
