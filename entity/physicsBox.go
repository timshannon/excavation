package entity

import (
	"excavation/engine"
)

type PhysicsBox struct {
	body *engine.PhysicsBody
}

func (p *PhysicsBox) Add(node *engine.Node, args EntityArgs) {
	offset := make([]float32, 16)
	collision := engine.PhysicsWorld().CreateBox(args.Float("x"), args.Float("y"), args.Float("z"),
		int(node.H3DNode), offset)
	p.body = engine.AddPhysicsBodyFromCollision(node, collision, args.Float("mass"))
}

func (p *PhysicsBox) Trigger(value float32) {
	return
}
