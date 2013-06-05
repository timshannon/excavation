// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package entity

import (
	"excavation/engine"
)

type PhysicsBox struct {
	body *engine.PhysicsBody
}

func (p *PhysicsBox) Add(node *engine.Node, args EntityArgs) {
	collision := engine.PhysicsWorld().CreateBox(args.Float("x"), args.Float("y"), args.Float("z"),
		int(node.H3DNode), &[16]float32{})
	p.body = engine.AddPhysicsBodyFromCollision(node, collision, args.Float("mass"))
}

func (p *PhysicsBox) Trigger(value float32) {
	return
}
