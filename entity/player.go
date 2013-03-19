package entity

import (
	"code.google.com/p/vmath"
	"excavation/engine"
	"math"
)

const (
	maxSpeed        = 20
	acceleration    = 100
	mouseMultiplier = 0.001 // makes for some saner numbers in the config file
)

var input [3]int
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
	lastUpdate   float64
	speed        *vmath.Vector3
	curVx, curVy int
}

func (p *Player) Add(node *engine.Node, args EntityArgs) {
	p.node = node

	//TODO: Only activate camera if set to active in arg
	engine.SetMainCamera(&engine.Camera{p.node})
	engine.MainCamera().SetOcclusionCulling(true)

	p.translate = new(vmath.Vector3)
	p.rotate = new(vmath.Vector3)
	p.rotationMatrix = new(vmath.Matrix3)
	p.curTranslate = new(vmath.Vector3)
	p.relM3 = new(vmath.Matrix3)
	p.speed = new(vmath.Vector3)

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

func (p *Player) Trigger(value float32) {
	//TODO: Fade view as camera changes?
	if value > 0 {
		engine.SetMainCamera(&engine.Camera{p.node})
		l := engine.AudioListener()
		l.SetNode(p.node)
		engine.SetMousePos(0, 0)
	}
}

func updatePlayer(t *engine.Task) {
	p := t.Data.(*Player)

	elapsedTime := float32(engine.GameTime() - p.lastUpdate)
	p.lastUpdate = engine.GameTime()

	for i := 0; i < 3; i++ {
		if input[i] == 0 {
			p.speed[i] = deccelerate(p.speed[i], elapsedTime)
		} else {
			p.speed[i] = accelerate(p.speed[i], elapsedTime, input[i])
		}
	}

	p.translate.ScalarMul(p.speed, elapsedTime)

	if !p.invert {
		p.rotate[1] = (float32(vY-p.curVy) * p.mouseSensitivity)
	} else {
		p.rotate[1] = (float32(vY-p.curVy) * (p.mouseSensitivity * -1))
	}

	p.rotate[0] = (float32(vX-p.curVx) * p.mouseSensitivity)

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
	p.rotationMatrix.MakeRotationZYX(p.rotate)
	matrix.Upper3x3(p.relM3)

	p.translate.MulM3Self(p.relM3)
	p.rotationMatrix.MulSelf(p.relM3)

	p.translate.AddToSelf(p.curTranslate)

	matrix.MakeFromM3V3(p.rotationMatrix, p.translate)
	p.node.SetRelativeTransMat(matrix)

	p.translate.ScalarMulSelf(0)
	p.rotate.ScalarMulSelf(0)
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
		input[2] += -1 * modifier
	case "Backward":
		input[2] += 1 * modifier
	case "StrafeLeft":
		input[0] += -1 * modifier
	case "StrafeRight":
		input[0] += 1 * modifier
	case "MoveUp":
		input[1] += 1 * modifier
	case "MoveDown":
		input[1] += -1 * modifier
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

//TODO: Split into two different types, one for ship movement, one for human movement,
// use trigger to switch between them
// Share one player type which holds player data.
