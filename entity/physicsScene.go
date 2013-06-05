// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

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
