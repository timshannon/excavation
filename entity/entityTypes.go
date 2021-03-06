// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package entity

import (
	"errors"
	"strings"
)

//NewEntity creates a new entity of the type passed in via a string
// this is so entities can be loaded from the xml scene file
func NewEntity(typeName string) (Entity, error) {

	//Big constantly changing switch for now, not sure of a
	// better way to handle this for now.
	switch strings.ToLower(typeName) {
	case "player":
		return new(Player), nil
	case "audio":
		return new(Audio), nil
	case "timer":
		return new(Timer), nil
	case "physicsobject":
		return new(PhysicsObject), nil
	case "physicsscene":
		return new(PhysicsScene), nil
	case "physicsbox":
		return new(PhysicsBox), nil

	}
	return nil, errors.New("Entity of type " + typeName + " not found.")

}
