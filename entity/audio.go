package entity

import (
	"excavation/engine"
	vmath "github.com/timshannon/vectormath"
)

type AudioStatic struct {
	audio  *engine.Audio
	buffer *engine.AudioBuffer
}

func (a *AudioStatic) Add(node *engine.Node, args EntityArgs) {
	position := new(vmath.Vector3)
	node.Translate(position)
	a.buffer = engine.NewAudioBuffer(args.String("file"))
	a.buffer.Load()

	a.audio = engine.AddStaticAudio(position, a.buffer, args.Float("minDistance"),
		args.Float("maxDistance"))

	a.audio.Play()
	a.audio.SetLooping(args.Bool("loop"))

}
