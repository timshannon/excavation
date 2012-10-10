package entity

import (
	"excavation/engine"
	"github.com/spate/vectormath"
)

var player *Player

type Player struct {
	node                     *engine.Node
	translate, rotate, scale *vectormath.Vector3
}

func (p *Player) load(node *engine.Node, args map[string]string) {
	p.node = node

	engine.SetMainCam(&engine.Camera{p.node})
	//set main player for input
	player = p

	p.translate = new(vectormath.Vector3)
	p.rotate = new(vectormath.Vector3)
	p.scale = new(vectormath.Vector3)

	engine.BindInput(handlePlayerInput, "Forward", "Backward", "Strafe_Right", "Strafe_Left")
	engine.AddTask("updatePlayer", updatePlayer, p, 0, 1)
}

func (p *Player) Trigger(value float32) {
	//not used
}

//returns the current active player
// non input controlled players (i.e. other multiplayer players)
// will be handled with a different entity
func MainPlayer() *Player {
	return player
}

func updatePlayer(t *engine.Task) {
	p := t.Data.(*Player)
	n := p.node

	//p.translate.SetZ(p.translate.Z() + -0.1)
	n.SetTransform(p.translate, p.rotate, p.scale)

}

func handlePlayerInput(i *engine.Input) {
	if i.ControlName() == "forward" && i.State == engine.StatePressed {
		player.translate.SetZ(player.translate.Z() - 0.1)
	}

}
