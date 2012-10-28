package entity

import (
	"excavation/engine"
	"github.com/spate/vectormath"
	"math"
)

const (
	maxSpeed     = 75.0
	acceleration = 1.5
)

var inX, inY, inZ int
var vX, vY int

type Player struct {
	node              *engine.Node
	translate, rotate *vectormath.Vector3

	invert                    bool
	lastUpdate                float64
	curAccX, curAccY, curAccZ float32
	curVx, curVy              int
}

func (p *Player) Add(node *engine.Node, args map[string]string) {
	p.node = node

	engine.SetMainCam(&engine.Camera{p.node})

	p.translate = new(vectormath.Vector3)
	p.rotate = new(vectormath.Vector3)

	//TODO: Handle config changes without having to look up the
	// setting constantly
	p.invert = engine.Cfg().Bool("InvertMouse")

	engine.BindInput(handlePlayerInput, "Forward", "Backward", "StrafeRight", "StrafeLeft", "MoveUp", "MoveDown")
	engine.BindInput(handlePitchYaw, "PitchYaw")
	engine.AddTask("updatePlayer", updatePlayer, p, 0, 0)
}

func updatePlayer(t *engine.Task) {
	p := t.Data.(*Player)
	n := p.node

	if inX == 0 {
		p.curAccX = deccelerate(p.curAccX)
	} else {
		p.curAccX = accelerate(p.curAccX, inX)
	}

	if inY == 0 {
		p.curAccY = deccelerate(p.curAccY)
	} else {
		p.curAccY = accelerate(p.curAccY, inY)
	}

	if inZ == 0 {
		p.curAccZ = deccelerate(p.curAccZ)
	} else {
		p.curAccZ = accelerate(p.curAccZ, inZ)
	}

	p.translate.SetX(p.curAccX * float32(engine.Time()-p.lastUpdate))
	p.translate.SetY(p.curAccY * float32(engine.Time()-p.lastUpdate))
	p.translate.SetZ(p.curAccZ * float32(engine.Time()-p.lastUpdate))

	//p.rotate.SetX(float32(vX))
	//p.rotate.SetY(float32(vY))

	p.rotate.SetX(0.01)

	n.SetLocalTransform(p.translate, p.rotate)

	p.lastUpdate = engine.Time()
}

func accelerate(speed float32, modifier int) float32 {
	speed += (float32(modifier) * acceleration)

	if math.Abs(float64(speed)) > maxSpeed {
		speed = (maxSpeed * float32(modifier))
	}
	return speed
}

func deccelerate(speed float32) float32 {
	var modifier float32
	if speed == 0 {
		return 0
	} else if speed > 0 {
		modifier = 1
	} else if speed < 0 {
		modifier = -1
	}

	speed -= (modifier * acceleration)
	if (speed * modifier) < 0 {
		speed = 0
	}

	return speed
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
		inZ += -1 * modifier
	case "Backward":
		inZ += 1 * modifier
	case "StrafeLeft":
		inX += -1 * modifier
	case "StrafeRight":
		inX += 1 * modifier
	case "MoveUp":
		inY += 1 * modifier
	case "MoveDown":
		inY += -1 * modifier
	}

}

func handlePitchYaw(i *engine.Input) {
	//TODO: handle joy and key input
	vX, vY, _ = i.MousePos()
}
