package entity

import (
	"excavation/engine"
	"github.com/spate/vectormath"
)

type Player struct {
	node                     *engine.Node
	translate, rotate, scale *vectormath.Vector3
}

func (p *Player) load(node *engine.Node, args map[string]string) {
	p.node = node

	engine.SetMainCam(&engine.Camera{p.node})

	p.translate = new(vectormath.Vector3)
	p.rotate = new(vectormath.Vector3)
	p.scale = new(vectormath.Vector3)

	engine.AddTask("updatePlayer", updatePlayer, p, 0, 1)
}

func (p *Player) Trigger(value float32) {
	//not used
}

func updatePlayer(t *engine.Task) {
	p := t.Data.(*Player)
	n := p.node

	n.Transform(p.translate, p.rotate, p.scale)

	//p.translate.SetZ(p.translate.Z() + -0.1)
	n.SetTransform(p.translate, p.rotate, p.scale)

}
