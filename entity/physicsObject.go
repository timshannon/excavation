package entity

import (
	"excavation/engine"
)

type PhysicsObject struct {
	body *engine.PhysicsBody
}

func (p *PhysicsObject) Add(node *engine.Node, args EntityArgs) {
	p.body = engine.AddPhysicsBody(node, args.Float("mass"))
	p.body.ShowDebugGeometry()

}

func (p *PhysicsObject) Trigger(value float32) {
	return
}
