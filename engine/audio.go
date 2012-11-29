package engine

import (
	"github.com/timshannon/go-openal/openal"
	vmath "github.com/timshannon/vectormath"
	"path"
)

const (
	audioRollOffDefault = 0.5
	//if Audio Node is occluded, sound drops off quicker
	audioRollOffOccluded = 4
)

// All audio files must be Mono wav files at 44100 Hz
const AudioFrequency = 44100

type Listener struct {
	openal.Listener
	node               *Node
	upOrient, atOrient *openal.Vector
	tempVec            *vmath.Vector4
	prevTime           float64
	curVec, prevVec    *vmath.Vector3
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
		curVec:   new(vmath.Vector3),
		prevVec:  new(vmath.Vector3),
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
	s.SetGain(newAudio.gain)
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
	gain        float32
	source      *audioSource
	//TODO: optional velocity
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
		gain:     1.0, //openal default
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

	//TODO: Streaming - Stream based on maxBufferSize.
	// if total size of file > maxBuffer size, then split into buffers the
	// size of maxBufferSize
	a.SetData(openal.FormatMono16, data, AudioFrequency)
	a.loaded = true
	return nil
}

func (a *Audio) Play() {
	if a.source == nil {
		if len(sources) < maxAudioSources {
			newSource := &audioSource{Source: openal.NewSource(),
				audio: a}
			newSource.setAudio(a)
			a.source = newSource

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

func (a *Audio) SetGain(value float32) {
	a.gain = value
	if a.source != nil {
		a.source.SetGain(value)
	}
}

func (a *Audio) Gain() float32 {
	return a.gain
}

func (a *Audio) State() int {
	if a.source != nil {
		return int(a.source.State())
	}
	return int(openal.Stopped)
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

				if sources[i].occluded() {
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

func (s *audioSource) occluded() bool {
	//TODO: Fix occlusion test.  Not visible to camera, but 
	// visible to ears.  A sound behind your head isn't automatically
	// occluded because you can't see it 
	//360 camera attached to listener?
	//ray cast from source to listener
	//return sources[i].audio.node.IsVisible(MainCam, true, false) == -1
	return false
}

func (l *Listener) updatePositionOrientation() {

	l.node.AbsoluteTransMat().Translation(l.curVec)
	l.Set3f(openal.AlPosition, l.curVec.X, l.curVec.Y, l.curVec.Z)

	//forward
	vmath.V4MakeZAxis(l.tempVec)
	l.tempVec.Z = -1 //horde has flipped z
	setOpenAlRelativeVector(l.atOrient, l.tempVec, l.node.AbsoluteTransMat())

	//up
	vmath.V4MakeYAxis(l.tempVec)
	setOpenAlRelativeVector(l.upOrient, l.tempVec, l.node.AbsoluteTransMat())

	l.SetOrientation(listener.atOrient, listener.upOrient)

	vmath.V3Velocity(l.prevVec, l.prevVec, l.curVec, float32(Time()-l.prevTime))

	l.Set3f(openal.AlVelocity, l.prevVec.X, l.prevVec.Y, l.prevVec.Z)

	vmath.V3Copy(l.prevVec, l.curVec)
	l.prevTime = Time()
}

func setOpenAlRelativeVector(alVec *openal.Vector, v4 *vmath.Vector4, matrix *vmath.Matrix4) {
	vmath.M4MulV4(v4, matrix, v4)
	v4.Normalize()

	alVec.X = v4.X
	alVec.Y = v4.Y
	alVec.Z = v4.Z

}
