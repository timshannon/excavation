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
	//TriggerIn(float32)
	//TriggerOut(*Entity)
	//TriggerSource?
}

type EntityArgs map[string]string

var entities = make(map[int]Entity)

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

	entities[int(node.H3DNode)] = newEnt

	return nil

}

//TODO: Match entity triggers
// if entity has an arg of trigger, then
// add it to the list of trigger entities
// if trigger == true, then auto trigger the entity

func EntityFromNode(node engine.Node) (Entity, bool) {
	entity, ok := entities[int(node.H3DNode)]
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
