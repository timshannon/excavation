// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"excavation/engine"
	"flag"
	"fmt"
)

//cmd line options
var (
	collisionType string
)

func init() {
	flag.StringVar(&collisionType, "type", "scene", "Type of collision to serialize: scene or compound.")

	flag.Parse()
}

func main() {
	switch collisionType {
	case "scene":
		fmt.Println("Processing Scene Collision.")
	case "compound":
		fmt.Println("Processing Compound Collision.")
	default:
		fmt.Println("Invalid collision type. Must be scene or compound.")

	}
}
