package entity

import (
	"excavation/engine"
)

type Audio struct {
	audio *engine.Audio
}

func (a *Audio) Add(node *engine.Node, args EntityArgs) {
	a.audio = engine.AddAudioNode(node, args.String("file"), args.Float("minDistance"),
		args.Float("maxDistance"), 10)

	a.audio.Load()
	//a.audio.Play() //TODO: entity triggers
	a.audio.SetLooping(args.Bool("loop"))
	//TODO: task to check distance and automatically start and stop audio based on
	// distance from listener i.e 2xMaxDistance
}
