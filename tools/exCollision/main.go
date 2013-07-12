// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"bitbucket.org/tshannon/gohorde/horde3d"
	"bitbucket.org/tshannon/gonewton/newton"
	"excavation/engine"
	"flag"
	"fmt"
	"github.com/jteeuwen/glfw"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

//cmd line options
var (
	collisionType   string
	outputFile      string
	convexTolerance float64
)

func init() {
	flag.StringVar(&collisionType, "type", "scene", "Type of collision to serialize: scene or compound.")
	flag.StringVar(&outputFile, "file", "", "Output file name. If no name is specified the current node "+
		"name will be used. <type>.ngd will be added.")
	flag.Float64Var(&convexTolerance, "tolerance", 0.01, "Tolerance allowed when converting mesh to convex "+
		"collision. A higher number will simplify the mesh more, and is useful for highly detailed models.")
}

func main() {
	flag.Parse()
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := glfw.OpenWindow(10,
		10, 8, 8, 8, 8,
		24, 8,
		glfw.Windowed); err != nil {
		fmt.Println(err.Error())
		return
	}

	horde3d.Init()
	engine.InitPhysics()

	var scene *engine.Scene
	var collision *newton.Collision

	nodeName := flag.Arg(0)

	scene, err := engine.NewScene(nodeName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	data, err := loadData(nodeName)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	good := scene.H3DRes.Load(data)

	if !good {
		fmt.Println("Horde3D was unable to load the resource " + nodeName + ".")
		return
	}

	node, err := engine.Root.AddScene(scene)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("ResRoot: ", nodeName)

	//TODO: Get node geom, load it directly from the same directory
	//	as the root node

	if outputFile == "" {
		outputFile = node.Name()
	}
	if !strings.HasSuffix(outputFile, ".ngd") {
		outputFile += "." + collisionType + ".ngd"

	}

	switch collisionType {
	case "scene":
		//Scene is a concave, static collision model with no overlapping polygons
		fmt.Println("Processing Scene Collision.")
		collision = engine.NewtonTreeFromNode(engine.Root)
	case "compound":
		//Compound is a convex only, dynamic collision which breaks the individual
		// horde meshes into separate parts of the compound collision
		fmt.Println("Processing Compound Collision.")
		collision = buildCompoundCollision(engine.Root)
	default:
		fmt.Println("Invalid collision type. Must be scene or compound.")
		return
	}

	fmt.Println("Collision: ", collision)

	file, err := os.Create(outputFile)
	if os.IsExist(err) {
		fmt.Println("File " + outputFile + " already exists and is being overwritting")
	} else if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}

	engine.PhysicsWorld().SerializeCollision(collision, writeToCollisionFile, file)
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}

	horde3d.Release()
	glfw.Terminate()
	glfw.CloseWindow()

}

func writeToCollisionFile(file interface{}, buffer []byte) {
	fmt.Println(len(buffer))
	_, err := file.(*os.File).Write(buffer)
	if err != nil {
		file.(*os.File).Close()
		panic(err)
	}
}

func buildCompoundCollision(node *engine.Node) *newton.Collision {

	meshes := engine.NewtonMeshListFromNode(node)

	if len(meshes) == 1 {
		collision := engine.PhysicsWorld().CreateConvexHullFromMesh(meshes[0], float32(convexTolerance),
			int(node.H3DNode))
		return collision
	}

	collision := engine.PhysicsWorld().CreateCompoundCollision(int(node.H3DNode))

	collision.CompoundBeginAddRemove()
	for i := range meshes {
		subCollision := engine.PhysicsWorld().CreateConvexHullFromMesh(meshes[i], float32(convexTolerance),
			int(node.H3DNode))
		collision.CompoundAddSubCollision(subCollision)
	}
	collision.CompoundEndAddRemove()

	return collision
}

func loadData(resourcePath string) ([]byte, error) {

	data, err := ioutil.ReadFile(resourcePath)

	if os.IsNotExist(err) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}
