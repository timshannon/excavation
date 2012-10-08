package entity

import (
	"excavation/engine"
)

type Player struct {
	node *engine.Node
}

func (p *Player) load(node *engine.Node, args map[string]string) {
	p.node = node

	engine.SetMainCam(&engine.Camera{p.node})

}

func (p *Player) Trigger(value float32) {
	//not used
}
