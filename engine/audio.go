package engine

import (
	"github.com/timshannon/go-openal/openal"
	"github.com/timshannon/vectormath"
	"path"
)

type Listener struct {
	openal.Listener
	node               *Node
	upOrient, atOrient *openal.Vector
	matrix             *vectormath.Matrix4
	tempVector         *vectormath.Vector4
}

var listener *Listener
var openalDevice *openal.Device
var openalContext *openal.Context

//TODO: Hard limit on # of sources 32? 64? config
var audioNodes []*Audio

func initAudio(deviceName string) {
	listener = &Listener{
		Listener:   openal.Listener{},
		upOrient:   new(openal.Vector),
		atOrient:   new(openal.Vector),
		matrix:     new(vectormath.Matrix4),
		tempVector: new(vectormath.Vector4),
	}

	openalDevice = openal.OpenDevice(deviceName)
	openalContext = openalDevice.CreateContext()
	openal.SetDistanceModel(openal.LinearDistanceClamped)
	openalContext.Activate()
	audioNodes = make([]*Audio, 0, 10)
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
	openalContext.Destroy()
	openalContext = openalDevice.CreateContext()
	openalContext.Activate()

}

type Audio struct {
	openal.Source
	node *Node
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

	audioNodes = append(audioNodes, aNode)
	return aNode
}

//AddStaticAudio Adds an audio source that doesn't move
func AddStaticAudio(position *vectormath.Vector3, buffer *AudioBuffer,
	minDistance, maxDistance float32) *Audio {
	aNode := &Audio{Source: openal.NewSource()}
	aNode.SetBuffer(buffer.Buffer)

	aNode.SetReferenceDistance(minDistance)
	aNode.SetMaxDistance(maxDistance)

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
	return nil
}

func updateAudio() {
	//TODO: Track velocity
	if listener.node == nil {
		return
	}

	//TODO: Option to track occlusion.  If source is
	// occluded, muffle the sound
	listener.updatePositionOrientation()
	//for i := range audioNodes {
	//horde3d.GetNodeTransform(audioNodes[i].node.H3DNode, &x, &y, &z,
	//&rx, &ry, &rz, nil, nil, nil)
	//audioNodes[i].source.Set3f(openal.AlPosition, x, y, z)
	//audioNodes[i].source.Set3f(openal.AlDirection, rx, ry, rz)
}

func (l *Listener) updatePositionOrientation() {
	//TODO: Move to Player controller
	l.node.AbsoluteTransMat(listener.matrix)

	l.Set3f(openal.AlPosition, l.matrix.GetElem(3, 0),
		l.matrix.GetElem(3, 1), l.matrix.GetElem(3, 2))

	//forward
	l.tempVector.X = 0
	l.tempVector.Y = 0
	l.tempVector.Z = -1
	setOpenAlRelativeVector(l.atOrient, l.tempVector, l.matrix)

	//up
	l.tempVector.X = 0
	l.tempVector.Y = 1
	l.tempVector.Z = 0
	setOpenAlRelativeVector(l.upOrient, l.tempVector, l.matrix)

	l.SetOrientation(listener.atOrient, listener.upOrient)
}

func setOpenAlRelativeVector(alVec *openal.Vector, v4 *vectormath.Vector4, matrix *vectormath.Matrix4) {
	vectormath.M4MulV4(v4, matrix, v4)
	vectormath.V4Normalize(v4, v4)

	alVec.X = v4.X
	alVec.Y = v4.Y
	alVec.Z = v4.Z

}
