package entity

import (
	"excavation/engine"
)

type PhysicsObject struct {
	body *engine.PhysicsBody
}

func (p *PhysicsObject) Add(node *engine.Node, args EntityArgs) {
	matrix := make([]float32, 16)
	p.body = engine.AddPhysicsBody(node, args.Float("mass"))

}

func (p *PhysicsObject) Trigger(value float32) {
	return
}
