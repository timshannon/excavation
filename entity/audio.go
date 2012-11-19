package entity

import (
	"excavation/engine"
)

type AudioStatic struct {
	audio *engine.Audio
}

func (a *AudioStatic) Add(node *engine.Node, args EntityArgs) {
	a.audio = engine.AddAudioNode(node, args.String("file"), args.Float("minDistance"),
		args.Float("maxDistance"), 10)

	a.audio.Load()
	//a.audio.Play()
	a.audio.Looping = args.Bool("loop")
}
