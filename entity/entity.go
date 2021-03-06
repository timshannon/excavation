// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package entity

import (
	"encoding/xml"
	"errors"
	"excavation/engine"
	"strconv"
	"strings"
)

type Entity interface {
	Add(node *engine.Node, args EntityArgs)
	Trigger(float32)
}

type EntityArgs map[string]string

var entities = make(map[string]Entity)

func LoadEntity(node *engine.Node, attachmentData string) error {

	var newEnt Entity
	reader := strings.NewReader(attachmentData)
	decoder := xml.NewDecoder(reader)

	element, err := decoder.Token()
	if err != nil {
		return err
	}

	attr := element.(xml.StartElement).Attr
	args := make(EntityArgs)

	for i := range attr {
		if strings.ToLower(attr[i].Name.Local) == "type" {
			newEnt, err = NewEntity(attr[i].Value)
			if err != nil {
				return err
			}
		} else {
			args[attr[i].Name.Local] = attr[i].Value
		}
	}

	newEnt.Add(node, args)

	entities[node.Name()] = newEnt

	return nil

}

func EntityFromName(name string) (Entity, bool) {
	entity, ok := entities[name]
	return entity, ok
}

func (e EntityArgs) Bool(argName string) bool {
	value, ok := e[argName]
	if !ok {
		engine.RaiseError(errors.New("Entity has no argument of name " + argName))
		return false
	}

	switch strings.ToLower(value) {
	case "true":
		return true
	case "1":
		return true
	}
	return false
}

func (e EntityArgs) Float(argName string) float32 {
	value, ok := e[argName]
	if !ok {
		engine.RaiseError(errors.New("Entity has no argument of name " + argName))
		return 0
	}

	fValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0
	}
	return float32(fValue)
}

func (e EntityArgs) String(argName string) string {
	value, ok := e[argName]
	if !ok {
		engine.RaiseError(errors.New("Entity has no argument of name " + argName))
		return ""
	}

	return value
}
