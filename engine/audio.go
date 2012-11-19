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

var audioNodes []*Audio
var sources []*audioSource

//TODO: Audio Node priority, and inactivating audio nodes

func initAudio(deviceName string, maxSources int) {
	listener = &Listener{
		Listener: openal.Listener{},
		upOrient: new(openal.Vector),
		atOrient: new(openal.Vector),
		tempVec:  new(vmath.Vector4),
	}
	maxAudioSources = maxSources

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
	audioNodes = make([]*Audio, 0, maxAudioSources)
	openalContext.Destroy()
	openalContext = openalDevice.CreateContext()
	openalContext.Activate()

}

type audioSource struct {
	openal.Source
	audio *Audio
}

type Audio struct {
	openal.Buffer
	node        *Node
	Priority    int
	file        string
	loaded      bool
	Occlude     bool
	Looping     bool
	minDistance float32
	maxDistance float32
	active      bool
}

//TODO: Add option for static audio
// so we don't have to take time updating position

func playAudioNode(audioNode *Audio) {
	if len(sources) < maxAudioSources {
		newSource := &audioSource{Source: openal.NewSource(),
			audio: audioNode}
		newSource.update()
		newSource.SetRolloffFactor(audioRollOffDefault)

		sources = append(sources, newSource)
	} else {
		//evaluate priorities
	}
}

func (s *audioSource) update() {
	s.SetLooping(s.audio.Looping)
	s.SetMaxDistance(s.audio.maxDistance)
	s.SetReferenceDistance(s.audio.minDistance)
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

func updateAudio() {
	if listener.node == nil {
		return
	}

	//TODO: Option to track occlusion.  If source is
	// occluded, muffle the sound 
	/*source := audioNodes[0]
	if source.Occlude {
		if source.node.IsVisible(MainCam, true, false) {
			source.SetRolloffFactor(audioRollOffOccluded)

		else {
			source.SetRolloffFactor(audioRollOffDefault)
		}
	}
	*/
	listener.updatePositionOrientation()

	//for i := range audioNodes {
	//horde3d.GetNodeTransform(audioNodes[i].node.H3DNode, &x, &y, &z,
	//&rx, &ry, &rz, nil, nil, nil)
	//audioNodes[i].source.Set3f(openal.AlPosition, x, y, z)
	//audioNodes[i].source.Set3f(openal.AlDirection, rx, ry, rz)
}

func (l *Listener) updatePositionOrientation() {
	//TODO: Track velocity

	l.Set3f(openal.AlPosition, l.node.AbsoluteTransMat().GetElem(3, 0),
		l.node.AbsoluteTransMat().GetElem(3, 1),
		l.node.AbsoluteTransMat().GetElem(3, 2))

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
