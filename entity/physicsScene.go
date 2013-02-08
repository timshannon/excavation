package entity

import (
	"excavation/engine"
)

type PhysicsScene struct {
	body *engine.PhysicsScene
}

func (p *PhysicsScene) Add(node *engine.Node, args EntityArgs) {
	p.body = engine.AddPhysicsScene(node)
}

func (p *PhysicsScene) Trigger(value float32) {
	return
}
