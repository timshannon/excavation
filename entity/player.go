package entity

import (
	"excavation/engine"
	"github.com/spate/vectormath"
	"math"
)

const (
	maxSpeed     = 75
	acceleration = 400
)

var inX, inY, inZ int
var vX, vY int

type Player struct {
	node              *engine.Node
	translate, rotate *vectormath.Vector3

	//mouse
	invert           bool
	mouseSensitivity float32

	//movement
	lastUpdate             float64
	speedX, speedY, speedZ float32
	curVx, curVy           int
}

func (p *Player) Add(node *engine.Node, args map[string]string) {
	p.node = node

	engine.SetMainCam(&engine.Camera{p.node})

	p.translate = new(vectormath.Vector3)
	p.rotate = new(vectormath.Vector3)

	p.invert = engine.Cfg().Bool("InvertMouse")
	p.mouseSensitivity = engine.Cfg().Float("MouseSensitivity")

	//test: does this actually work?
	engine.Cfg().RegisterOnWriteHandler(func(cfg *engine.Config) {
		p.invert = engine.Cfg().Bool("InvertMouse")
		p.mouseSensitivity = engine.Cfg().Float("MouseSensitivity")
	})

	engine.BindInput(handlePlayerInput, "Forward", "Backward", "StrafeRight", "StrafeLeft", "MoveUp", "MoveDown")
	engine.BindInput(handlePitchYaw, "PitchYaw", "PitchUp", "PitchDown", "YawLeft", "YawRight")
	engine.AddTask("updatePlayer", updatePlayer, p, 0, 0)
}

func updatePlayer(t *engine.Task) {
	p := t.Data.(*Player)
	n := p.node

	elapsedTime := float32(engine.Time() - p.lastUpdate)
	p.lastUpdate = engine.Time()

	if inX == 0 {
		p.speedX = deccelerate(p.speedX, elapsedTime)
	} else {
		p.speedX = accelerate(p.speedX, elapsedTime, inX)
	}

	if inY == 0 {
		p.speedY = deccelerate(p.speedY, elapsedTime)
	} else {
		p.speedY = accelerate(p.speedY, elapsedTime, inY)
	}

	if inZ == 0 {
		p.speedZ = deccelerate(p.speedZ, elapsedTime)
	} else {
		p.speedZ = accelerate(p.speedZ, elapsedTime, inZ)
	}

	p.translate.SetX(p.speedX * elapsedTime)
	p.translate.SetY(p.speedY * elapsedTime)
	p.translate.SetZ(p.speedZ * elapsedTime)

	if p.invert {
		p.mouseSensitivity *= -1
	}
	//fmt.Println(float32(vX-p.curVx) * p.mouseSpeed)
	//p.rotate.SetX(float32(vX-p.curVx) * p.mouseSpeed)
	//p.rotate.SetY(float32(vY-p.curVy) * p.mouseSpeed)

	//p.rotate.SetX(0.01)
	//p.rotate.SetY(0.01)
	if p.translate.X() != 0 || p.translate.Y() != 0 || p.translate.Z() != 0 {
		n.SetLocalTransform(p.translate, p.rotate)
	}
	p.curVx = vX
	p.curVy = vY
}

func accelerate(speed, time float32, modifier int) float32 {
	speed += float32(modifier) * (acceleration * time)

	if math.Abs(float64(speed)) > maxSpeed {
		speed = (maxSpeed * float32(modifier))
	}
	return speed
}

func deccelerate(speed, time float32) float32 {
	var modifier float32
	if speed == 0 {
		return 0
	} else if speed > 0 {
		modifier = 1
	} else if speed < 0 {
		modifier = -1
	}

	speed -= modifier * (acceleration * time)
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
	x, y, ok := i.MousePos()

	if ok {
		vY = x
		vX = y
		return
	}

	state, ok := i.ButtonState()
	var modifier int
	if ok {
		if state == engine.StatePressed {
			modifier = 1
		} else {
			modifier = -1
		}

		switch i.ControlName() {
		case "PitchDown":
			vX += -1 * modifier
		case "PitchUp":
			vX += 1 * modifier
		case "YawLeft":
			vY += -1 * modifier
		case "YawRight":
			vY += 1 * modifier
		}
	}
}
