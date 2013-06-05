// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package entity

import (
	"excavation/engine"
)

type Audio struct {
	*engine.Audio
}

func (a *Audio) Add(node *engine.Node, args EntityArgs) {
	a.Audio = engine.AddAudioNode(node, args.String("file"), args.Float("minDistance"),
		args.Float("maxDistance"), 10)

	a.Load()
	a.SetLooping(args.Bool("loop"))
	a.Occlude = args.Bool("occlude")
	//TODO: task to check distance and automatically start and stop audio based on
	// distance from listener i.e 2xMaxDistance

	if args.Bool("autoStart") {
		a.Trigger(1)
	}
}

func (a *Audio) Trigger(value float32) {
	if value > 0 {
		if a.State() == engine.AudioStopped {
			a.Play()
		} else {
			a.SetGain(value)
		}
	} else {
		a.Stop()
	}
}
