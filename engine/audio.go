package engine

import (
	"github.com/spate/vectormath"
	"github.com/timshannon/go-openal/openal"
	"path"
)

type Listener struct {
	node                         *Node
	listener                     openal.Listener
	position, upOrient, atOrient *openal.Vector
	matrix                       *vectormath.Matrix4
	tempVector                   *vectormath.Vector4
}

var listener *Listener
var openalDevice *openal.Device
var openalContext *openal.Context

var audioNodes []*Audio

func initAudio(deviceName string) {
	listener = new(Listener)
	listener.listener = openal.Listener{}
	listener.position = new(openal.Vector)
	listener.upOrient = new(openal.Vector)
	listener.atOrient = new(openal.Vector)
	listener.tempVector = new(vectormath.Vector4)
	listener.matrix = new(vectormath.Matrix4)

	openalDevice = openal.OpenDevice(deviceName)
	openalContext = openalDevice.CreateContext()
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
	node   *Node
	source openal.Source
}

//AddAudioNode adds an audio source who's position gets
// updated based on the passed in node's position
func AddAudioNode(node *Node, buffer *AudioBuffer) *Audio {
	aNode := new(Audio)
	aNode.node = node
	aNode.source = openal.NewSource()
	aNode.source.SetBuffer(buffer.buffer)

	audioNodes = append(audioNodes, aNode)
	return aNode
}

//AddStaticAudio Adds an audio source that doesn't move
func AddStaticAudio(position *vectormath.Vector3, buffer *AudioBuffer) *Audio {
	aNode := new(Audio)
	aNode.source = openal.NewSource()
	aNode.source.SetBuffer(buffer.buffer)

	aNode.source.Set3f(openal.AlPosition,
		position.X(), position.Y(), position.Z())

	return aNode

}

//convience methods
// Anything more complicated and they can use the
// openal source directly
func (a *Audio) Play()   { a.source.Play() }
func (a *Audio) Stop()   { a.source.Stop() }
func (a *Audio) Pause()  { a.source.Pause() }
func (a *Audio) Rewind() { a.source.Rewind() }
func (a *Audio) SetLooping(looping bool) {
	a.source.SetLooping(looping)
}
func (a *Audio) Looping() bool { return a.source.GetLooping() }

func (a *Audio) OpenAlSource() openal.Source {
	return a.source
}

type AudioBuffer struct {
	buffer openal.Buffer
	file   string
	loaded bool
}

func NewAudioBuffer(file string) *AudioBuffer {
	abuffer := new(AudioBuffer)
	abuffer.buffer = openal.NewBuffer()
	abuffer.file = file

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
	b.buffer.SetData(openal.FormatMono16, data, 44100)
	return nil
}

func updateAudio() {
	//TODO: Track velocity
	if listener.node == nil {
		return
	}

	listener.updatePositionOrientation()
	//for i := range audioNodes {
	//horde3d.GetNodeTransform(audioNodes[i].node.H3DNode, &x, &y, &z,
	//&rx, &ry, &rz, nil, nil, nil)
	//audioNodes[i].source.Set3f(openal.AlPosition, x, y, z)
	//audioNodes[i].source.Set3f(openal.AlDirection, rx, ry, rz)
}

func (l *Listener) updatePositionOrientation() {
	l.node.AbsoluteTransMat(listener.matrix)

	l.position.X = l.matrix.GetElem(3, 0)
	l.position.Y = l.matrix.GetElem(3, 1)
	l.position.Z = l.matrix.GetElem(3, 2)
	l.listener.SetPosition(listener.position)

	//forward
	l.tempVector.SetX(0)
	l.tempVector.SetY(0)
	l.tempVector.SetZ(-1)
	setOpenAlRelativeVector(l.atOrient, l.tempVector, l.matrix)

	//up
	l.tempVector.SetX(0)
	l.tempVector.SetY(1)
	l.tempVector.SetZ(0)
	setOpenAlRelativeVector(l.upOrient, l.tempVector, l.matrix)

	l.listener.SetOrientation(listener.atOrient, listener.upOrient)
}

func setOpenAlRelativeVector(alVec *openal.Vector, v4 *vectormath.Vector4, matrix *vectormath.Matrix4) {
	vectormath.M4MulV4(v4, matrix, v4)
	vectormath.V4Normalize(v4, v4)

	alVec.X = v4.X()
	alVec.Y = v4.Y()
	alVec.Z = v4.Z()

}
