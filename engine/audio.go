package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"github.com/spate/vectormath"
	"github.com/timshannon/go-openal/openal"
	"io/ioutil"
)

type Listener struct {
	node     *Node
	listener openal.Listener
}

var listener *Listener
var openalDevice *openal.Device
var openalContext *openal.Context

var audioNodes []*Audio

func initAudio(deviceName string) {
	listener = new(Listener)
	listener.listener = openal.Listener{}
	openalDevice = openal.OpenDevice(deviceName)
	openalContext = openalDevice.CreateContext()
	audioNodes = make([]*Audio, 0, 10)
}

func (l Listener) SetNode(node *Node) {
	l.node = node
}

func ClearAllAudio() {
	for i := range audioNodes {
		openal.DeleteSource(audioNodes[i].source)
	}
	openalContext.Destroy()
	openalContext = openalDevice.CreateContext()

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
	data, err := ioutil.ReadFile(b.file)

	if err != nil {
		addError(err)
		return err
	}

	//TODO: Get wave file info
	//Mono only?  rate from config?
	//TODO: Streaming - Stream based on an arbitrary size
	// or let the user decide?
	b.buffer.SetData(openal.FormatMono16, data, 44100)
	return nil
}

func updateAudio() {
	var x, y, z, rx, ry, rz float32
	//TODO: Track velocity

	horde3d.GetNodeTransform(listener.node.H3DNode, &x, &y, &z,
		&rx, &ry, &rz, nil, nil, nil)
	listener.listener.Set3f(openal.AlPosition, x, y, z)
	listener.listener.Set3f(openal.AlDirection, rx, ry, rz)

	for i := range audioNodes {
		horde3d.GetNodeTransform(audioNodes[i].node.H3DNode, &x, &y, &z,
			&rx, &ry, &rz, nil, nil, nil)
		audioNodes[i].source.Set3f(openal.AlPosition, x, y, z)
		audioNodes[i].source.Set3f(openal.AlDirection, rx, ry, rz)
	}
}
