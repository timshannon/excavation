// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package entity

import (
	"excavation/engine"
)

type PhysicsObject struct {
	body *engine.PhysicsBody
}

func (p *PhysicsObject) Add(node *engine.Node, args EntityArgs) {
	p.body = engine.AddPhysicsBody(node, args.Float("mass"))

}

func (p *PhysicsObject) Trigger(value float32) {
	return
}
