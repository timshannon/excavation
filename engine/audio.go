package engine

import (
	"github.com/timshannon/go-openal/openal"
	vmath "github.com/timshannon/vectormath"
	"path"
)

const (
	audioRollOffDefault = 0.5
	//If Audio Node is occulded, sound drops off quicker
	audioRollOffOccluded = 4
)

type Listener struct {
	openal.Listener
	node               *Node
	upOrient, atOrient *openal.Vector
	tempVector         *vmath.Vector4
}

var listener *Listener
var openalDevice *openal.Device
var openalContext *openal.Context
var maxAudioSources int

var audioNodes []*Audio

//TODO: Audio Node priority, and inactivating audio nodes

func initAudio(deviceName string, maxSources int) {
	listener = &Listener{
		Listener:   openal.Listener{},
		upOrient:   new(openal.Vector),
		atOrient:   new(openal.Vector),
		tempVector: new(vmath.Vector4),
	}
	maxAudioSources = maxSources

	openalDevice = openal.OpenDevice(deviceName)
	openalContext = openalDevice.CreateContext()
	openal.SetDistanceModel(openal.InverseDistanceClamped)
	openalContext.Activate()
	audioNodes = make([]*Audio, 0, maxAudioSources)
}

func AudioListener() *Listener {
	return listener
}

func (l *Listener) SetNode(node *Node) {
	l.node = node
}

func ClearAllAudio() {
	//TODO: destroy clears sources, but what
	// about buffers?
	audioNodes = make([]*Audio, 0, maxAudioSources)
	openalContext.Destroy()
	openalContext = openalDevice.CreateContext()
	openalContext.Activate()

}

type Audio struct {
	openal.Source
	node    *Node
	Occlude bool
}

//AddAudioNode adds an audio source who's position gets
// updated based on the passed in node's position
func AddAudioNode(node *Node, buffer *AudioBuffer, minDistance,
	maxDistance float32) *Audio {
	aNode := &Audio{Source: openal.NewSource(),
		node: node,
	}
	aNode.node = node
	aNode.SetBuffer(buffer.Buffer)

	aNode.SetReferenceDistance(minDistance)
	aNode.SetMaxDistance(maxDistance)
	aNode.SetRolloffFactor(audioRollOffDefault)

	audioNodes = append(audioNodes, aNode)
	return aNode
}

//AddStaticAudio Adds an audio source that doesn't move
func AddStaticAudio(position *vmath.Vector3, buffer *AudioBuffer,
	minDistance, maxDistance float32) *Audio {
	aNode := &Audio{Source: openal.NewSource()}
	aNode.SetBuffer(buffer.Buffer)

	aNode.SetReferenceDistance(minDistance)
	aNode.SetMaxDistance(maxDistance)
	aNode.SetRolloffFactor(audioRollOffDefault)

	aNode.Set3f(openal.AlPosition,
		position.X, position.Y, position.Z)

	return aNode

}

type AudioBuffer struct {
	openal.Buffer
	file   string
	loaded bool
}

func NewAudioBuffer(file string) *AudioBuffer {
	abuffer := &AudioBuffer{Buffer: openal.NewBuffer(),
		file:   file,
		loaded: false,
	}

	return abuffer
}

func (b *AudioBuffer) Load() error {
	data, err := loadEngineData(path.Join(path.Join(dataDir, "sounds"), b.file))

	if err != nil {
		RaiseError(err)
		return err
	}

	//TODO: Get wave file info
	//Mono only?  rate from config?
	//TODO: Streaming - Stream based on an arbitrary size
	// or let the user decide? Config option?
	b.SetData(openal.FormatMono16, data, 44100)
	b.loaded = true
	return nil
}

func updateAudio() {
	//TODO: Track velocity
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
	l.Set3f(openal.AlPosition, l.node.AbsoluteTransMat().GetElem(3, 0),
		l.node.AbsoluteTransMat().GetElem(3, 1),
		l.node.AbsoluteTransMat().GetElem(3, 2))

	//forward
	vmath.V4MakeZAxis(l.tempVector)
	l.tempVector.Z = -1 //horde has flipped z
	setOpenAlRelativeVector(l.atOrient, l.tempVector, l.node.AbsoluteTransMat())

	//up
	vmath.V4MakeYAxis(l.tempVector)
	setOpenAlRelativeVector(l.upOrient, l.tempVector, l.node.AbsoluteTransMat())

	l.SetOrientation(listener.atOrient, listener.upOrient)
}

func setOpenAlRelativeVector(alVec *openal.Vector, v4 *vmath.Vector4, matrix *vmath.Matrix4) {
	vmath.M4MulV4(v4, matrix, v4)
	v4.Normalize()

	alVec.X = v4.X
	alVec.Y = v4.Y
	alVec.Z = v4.Z

}
