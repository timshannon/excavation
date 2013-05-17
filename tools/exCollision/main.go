// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"code.google.com/p/gonewton/newton"
	"excavation/engine"
	"flag"
	"fmt"
	"os"
	"strings"
)

//cmd line options
var (
	collisionType   string
	outputFile      string
	convexTolerance float64
)

func init() {
	flag.StringVar(&collisionType, "type", "scene", "Type of collision to serialize: scene or single, or compound.")
	flag.StringVar(&outputFile, "file", "", "Output file name. If no name is specified the current node name will be used. <type>.ngd will be added.")
	flag.Float64Var(&convexTolerance, "tolerance", 0.01, "Tolerance allowed when converting mesh to convex collision. A higher number will simplify the mesh more, and is useful for highly detailed models.")
	flag.Parse()
}

func main() {
	engine.InitPhysics()
	var node *engine.Scene
	var collision *newton.Collision

	nodeName := flag.Arg(0)

	node, err := engine.NewScene(nodeName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = node.Load()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

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
		collision = engine.NewtonTreeFromNode(node)
	case "compound":
		//Compound is a convex only, dynamic collision which breaks the individual
		// horde meshes into separate parts of the compound collision
		fmt.Println("Processing Compound Collision.")
	case "single":
		//Single is a convex only, dynamic collision which adds all horde meshes into
		// one single collision and generates a convex hull from the entire group
		fmt.Println("Processing Single Collision.")
	default:
		fmt.Println("Invalid collision type. Must be scene, single or compound.")
		return
	}

	file, err := os.Create(outputFile)
	if os.IsExist(err) {
		fmt.Println("File " + outputFile + " already exists and is being overwritting")
	}

	engine.PhysicsWorld().SerializeCollision(collision, writeToCollisionFile, file)
}

func writeToCollisionFile(file interface{}, buffer []byte) {
	_, err := file.(os.File).Write(buffer)
	if err != nil {
		panic(err)
	}
}
