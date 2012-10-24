package entity

import (
	"excavation/engine"
	"github.com/spate/vectormath"
	"math"
)

const (
	maxSpeed     = 2.0
	acceleration = 0.25
)

var pX, pY, pZ int

type Player struct {
	node              *engine.Node
	translate, rotate *vectormath.Vector3

	lastUpdate                float64
	curAccX, curAccY, curAccZ float64
}

func (p *Player) Add(node *engine.Node, args map[string]string) {
	p.node = node

	engine.SetMainCam(&engine.Camera{p.node})

	p.translate = new(vectormath.Vector3)
	p.rotate = new(vectormath.Vector3)

	engine.BindInput(handlePlayerInput, "Forward", "Backward", "StrafeRight", "StrafeLeft", "MoveUp", "MoveDown")
	engine.AddTask("updatePlayer", updatePlayer, p, 0, 0)
}

func updatePlayer(t *engine.Task) {
	p := t.Data.(*Player)
	n := p.node

	p.curAccX += math.Min((float64(pX) * acceleration), maxSpeed)
	p.curAccY += math.Min((float64(pY) * acceleration), maxSpeed)
	p.curAccZ += math.Min((float64(pZ) * acceleration), maxSpeed)

	p.translate.SetX(float32(p.curAccX * (engine.Time() - p.lastUpdate)))
	p.translate.SetY(float32(p.curAccY * (engine.Time() - p.lastUpdate)))
	p.translate.SetZ(float32(p.curAccZ * (engine.Time() - p.lastUpdate)))

	n.SetLocalTransform(p.translate, p.rotate)

	p.lastUpdate = engine.Time()
}

func handlePlayerInput(i *engine.Input) {
	var modifier int

	if i.State == engine.StatePressed {
		modifier = 1
	} else {
		modifier = -1
	}

	switch i.ControlName() {
	case "Forward":
		pZ += -1 * modifier
	case "Backward":
		pZ += 1 * modifier
	case "StrafeLeft":
		pX += -1 * modifier
	case "StrafeRight":
		pX += 1 * modifier
	case "MoveUp":
		pY += 1 * modifier
	case "MoveDown":
		pY += -1 * modifier
	}

}
