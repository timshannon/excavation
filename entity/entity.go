package entity

import (
	"encoding/xml"
	"excavation/engine"
	"strings"
)

type Entity interface {
	load(node *engine.Node, args map[string]string) //Called entity load
	Trigger(value float32)                          //Value may not be necessary
}

var entities []Entity

func LoadEntity(node *engine.Node, attachmentData string) error {

	var newEnt Entity
	reader := strings.NewReader(attachmentData)
	decoder := xml.NewDecoder(reader)

	element, err := decoder.Token()
	if err != nil {
		return err
	}

	attr := element.(xml.StartElement).Attr
	args := make(map[string]string)

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

	newEnt.load(node, args)

	entities = append(entities, newEnt)

	return nil

}
