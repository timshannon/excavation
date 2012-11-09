package entity

import (
	"excavation/engine"
	vmath "github.com/timshannon/vectormath"
	"math"
)

const (
	maxSpeed        = 75
	acceleration    = 400
	mouseMultiplier = 0.001 // makes for some saner numbers in the config file
)

var inX, inY, inZ int
var vX, vY int

type Player struct {
	node              *engine.Node
	translate, rotate *vmath.Vector3

	//Temp movement variables
	rotationMatrix *vmath.Matrix3
	curTranslate   *vmath.Vector3
	relM3          *vmath.Matrix3

	//mouse
	invert           bool
	mouseSensitivity float32

	//movement
	lastUpdate             float64
	speedX, speedY, speedZ float32
	curVx, curVy           int
}

func (p *Player) Add(node *engine.Node, args EntityArgs) {
	p.node = node

	engine.SetMainCam(&engine.Camera{p.node})

	p.translate = new(vmath.Vector3)
	p.rotate = new(vmath.Vector3)
	p.rotationMatrix = new(vmath.Matrix3)
	p.curTranslate = new(vmath.Vector3)
	p.relM3 = new(vmath.Matrix3)

	p.invert = engine.Cfg().Bool("InvertMouse")
	p.mouseSensitivity = engine.Cfg().Float("MouseSensitivity") * mouseMultiplier

	//test: does this actually work?
	engine.Cfg().RegisterOnWriteHandler(func(cfg *engine.Config) {
		p.invert = engine.Cfg().Bool("InvertMouse")
		p.mouseSensitivity = engine.Cfg().Float("MouseSensitivity") * mouseMultiplier
	})

	l := engine.AudioListener()
	l.SetNode(p.node)

	engine.BindInput(handlePlayerInput, "Forward", "Backward", "StrafeRight", "StrafeLeft",
		"MoveUp", "MoveDown")
	engine.BindInput(handlePitchYaw, "PitchYaw", "PitchUp", "PitchDown", "YawLeft", "YawRight")
	engine.AddTask("updatePlayer", updatePlayer, p, 0, 0)

	engine.SetMousePos(0, 0)

}

func updatePlayer(t *engine.Task) {
	p := t.Data.(*Player)

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

	p.translate.X = (p.speedX * elapsedTime)
	p.translate.Y = (p.speedY * elapsedTime)
	p.translate.Z = (p.speedZ * elapsedTime)

	if !p.invert {
		p.rotate.Y = (float32(vY-p.curVy) * p.mouseSensitivity)
	} else {
		p.rotate.Y = (float32(vY-p.curVy) * (p.mouseSensitivity * -1))
	}

	p.rotate.X = (float32(vX-p.curVx) * p.mouseSensitivity)

	p.localTransform()

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

func (p *Player) localTransform() {
	matrix := p.node.RelativeTransMat()
	matrix.Translation(p.curTranslate)
	vmath.M3MakeRotationZYX(p.rotationMatrix, p.rotate)
	matrix.Upper3x3(p.relM3)

	vmath.M3MulV3(p.translate, p.relM3, p.translate)
	vmath.M3Mul(p.rotationMatrix, p.relM3, p.rotationMatrix)

	vmath.V3Add(p.translate, p.curTranslate, p.translate)

	vmath.M4MakeFromM3V3(matrix, p.rotationMatrix, p.translate)
	p.node.SetRelativeTransMat(matrix)

	zeroVector(p.translate)
	zeroVector(p.rotate)
}

func zeroVector(vector *vmath.Vector3) {
	vector.X = 0
	vector.Y = 0
	vector.Z = 0
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
