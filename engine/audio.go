package engine

import (
	"github.com/timshannon/go-openal/openal"
	vmath "github.com/timshannon/vectormath"
	"path"
)

const (
	audioRollOffDefault = 0.5
	//if Audio Node is occulded, sound drops off quicker
	audioRollOffOccluded = 4
)

type Listener struct {
	openal.Listener
	node               *Node
	upOrient, atOrient *openal.Vector
	tempVec            *vmath.Vector4
	prevTime           float32
}

var listener *Listener
var openalDevice *openal.Device
var openalContext *openal.Context
var maxAudioSources int
var maxAudioBufferSize int

var audioNodes []*Audio
var sources []*audioSource

func initAudio(deviceName string, maxSources, maxBufferSize int) {
	listener = &Listener{
		Listener: openal.Listener{},
		upOrient: new(openal.Vector),
		atOrient: new(openal.Vector),
		tempVec:  new(vmath.Vector4),
	}
	maxAudioSources = maxSources
	maxAudioBufferSize = maxBufferSize

	openalDevice = openal.OpenDevice(deviceName)
	openalContext = openalDevice.CreateContext()
	openal.SetDistanceModel(openal.InverseDistanceClamped)
	openalContext.Activate()
	sources = make([]*audioSource, 0, maxAudioSources)
	audioNodes = make([]*Audio, 0, maxAudioSources)
}

func AudioListener() *Listener {
	return listener
}

func (l *Listener) SetNode(node *Node) {
	l.node = node
}

func ClearAllAudio() {
	//TODO: clear buffers
	for i := range audioNodes {
		audioNodes[i].Remove()
	}
	audioNodes = make([]*Audio, 0, maxAudioSources)
	openalContext.Destroy()
	openalContext = openalDevice.CreateContext()
	openalContext.Activate()
}

type audioSource struct {
	openal.Source
	audio *Audio
	free  bool
}

func (s *audioSource) listenerRelative() bool {
	return s.audio.node.H3DNode == listener.node.H3DNode
}

func (s *audioSource) setAudio(newAudio *Audio) {
	newAudio.source = s
	s.audio.source = nil
	s.audio = newAudio
	s.free = false
	s.SetLooping(newAudio.looping)
	s.SetMaxDistance(newAudio.maxDistance)
	s.SetReferenceDistance(newAudio.minDistance)
	if !newAudio.Occlude {
		s.SetRolloffFactor(audioRollOffDefault)
	}

	if s.listenerRelative() {
		//if the source of the sound is the same as the listener
		// don't update position
		s.SetSourceRelative(true)
		s.Set3f(openal.AlPosition, 0, 0, 0)
	} else {
		s.SetSourceRelative(false)
	}

	s.SetBuffer(newAudio.Buffer)

}

type Audio struct {
	openal.Buffer
	node        *Node
	Priority    int
	file        string
	loaded      bool
	Occlude     bool
	looping     bool
	minDistance float32
	maxDistance float32
	source      *audioSource
}

//AddAudioNode adds an audio source who's position gets
// updated based on the passed in node's position
func AddAudioNode(node *Node, audioFile string, minDistance,
	maxDistance float32, priority int) *Audio {
	aNode := &Audio{Buffer: openal.NewBuffer(),
		file:     audioFile,
		node:     node,
		Priority: priority,
		Occlude:  false,
		loaded:   false,
	}

	aNode.minDistance = minDistance
	aNode.maxDistance = maxDistance

	audioNodes = append(audioNodes, aNode)
	return aNode
}

func (a *Audio) Load() error {
	data, err := loadEngineData(path.Join(path.Join(dataDir, "sounds"), a.file))

	if err != nil {
		RaiseError(err)
		return err
	}

	//TODO: Get wave file info
	//Mono only?  rate from config?
	//TODO: Streaming - Stream based on an arbitrary size
	// or let the user decide? Config option?
	a.SetData(openal.FormatMono16, data, 44100)
	a.loaded = true
	return nil
}

func (a *Audio) Play() {
	if a.source == nil {
		if len(sources) < maxAudioSources {
			newSource := &audioSource{Source: openal.NewSource(),
				audio: a}
			newSource.setAudio(a)

			sources = append(sources, newSource)

			a.source.Play()
			return
		} else {
			//check for free'd source or lower priority audio
			var freeSource *audioSource
			var lowest *audioSource
			for i := range sources {
				if sources[i].free {
					freeSource = sources[i]
					break
				}
				if lowest == nil {
					lowest = sources[i]
				} else {
					if sources[i].audio.Priority > lowest.audio.Priority {
						lowest = sources[i]
					}
				}
			}
			if freeSource != nil {
				freeSource.setAudio(a)
			} else if lowest != nil {
				lowest.setAudio(a)
			} else {
				//can't play audio
				return
			}
			a.source.Play()
			return
		}
	}
	a.source.Play()
}

func (a *Audio) Pause() {
	if a.source != nil {
		a.source.Pause()
	}
}

func (a *Audio) SetLooping(value bool) {
	a.looping = value
	if a.source != nil {
		a.source.SetLooping(value)
	}
}

func (a *Audio) Stop() {
	if a.source != nil {
		a.source.Stop()

		if len(audioNodes) > maxAudioSources {
			//free up source
			a.freeSource()
		}
	}
}

func (a *Audio) Remove() {
	a.Stop()
	openal.DeleteBuffer(a.Buffer)
}

func (a *Audio) freeSource() {
	a.source.free = true
	a.source = nil

	if len(audioNodes) > maxAudioSources {
		for i := range audioNodes {
			if audioNodes[i].source == nil &&
				audioNodes[i].looping {
				audioNodes[i].Play()
			}
		}
	}
}

func updateAudio() {
	if listener.node == nil {
		return
	}

	for i := range sources {
		//TODO Option: Dont check every frame
		if sources[i].State() == openal.Stopped {
			sources[i].audio.freeSource()
		}

		if !sources[i].listenerRelative() {
			if sources[i].audio.Occlude {
				if sources[i].audio.node.IsVisible(MainCam, true, false) == -1 {
					sources[i].SetRolloffFactor(audioRollOffOccluded)
				} else {
					sources[i].SetRolloffFactor(audioRollOffDefault)
				}
			}
			//position
			sources[i].Set3f(openal.AlPosition, sources[i].audio.node.AbsoluteTransMat().Col3.X,
				sources[i].audio.node.AbsoluteTransMat().Col3.Y,
				sources[i].audio.node.AbsoluteTransMat().Col3.Z)

			//direction
			//Only needed for sound cones, may not implement
		}
	}

	listener.updatePositionOrientation()

}

func (l *Listener) updatePositionOrientation() {
	//TODO: Track velocity

	l.Set3f(openal.AlPosition, l.node.AbsoluteTransMat().Col3.X,
		l.node.AbsoluteTransMat().Col3.Y,
		l.node.AbsoluteTransMat().Col3.Z)

	//forward
	vmath.V4MakeZAxis(l.tempVec)
	l.tempVec.Z = -1 //horde has flipped z
	setOpenAlRelativeVector(l.atOrient, l.tempVec, l.node.AbsoluteTransMat())

	//up
	vmath.V4MakeYAxis(l.tempVec)
	setOpenAlRelativeVector(l.upOrient, l.tempVec, l.node.AbsoluteTransMat())

	l.SetOrientation(listener.atOrient, listener.upOrient)

}

func setOpenAlRelativeVector(alVec *openal.Vector, v4 *vmath.Vector4, matrix *vmath.Matrix4) {
	vmath.M4MulV4(v4, matrix, v4)
	v4.Normalize()

	alVec.X = v4.X
	alVec.Y = v4.Y
	alVec.Z = v4.Z

}
