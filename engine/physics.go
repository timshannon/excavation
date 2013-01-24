package engine

import (
	"code.google.com/p/gonewton/newton"
)

var newtonWorld *newton.World
var newtonLastUpdate float32

type PhysicsScene struct {
	Node *Node
	Body *newton.Body
}

type PhysicsBody struct {
	Node *Node
	Body *newton.Body
}

func initPhysics() {
	newtonWorld = newton.CreateWorld()
}

func updatePhysics() {
	newtonWorld.Update(float32(GameTime()) - newtonLastUpdate)
}

func clearAllPhysics() {
	newtonWorld.DestroyAllBodies()
}

//AddPhysicsScene adds a scene physics collision type built from
// the passed in node's geometry.  
func AddPhysicsScene(node *Node) *PhysicsScene {

}

//AddPhysicsSceneFromGeometry adds a scene physics collision type built
// from the passed in geometry resource
func AddPhysicsSceneFromGeometry(node *Node, geometry *Geometry) *PhysicsScene {

}

//Adds a physics body using the passed in node to determine the collision hull.
// a convex hull will be built from the node's visible geometry.  If the node has
// children a compound shape will be made from the children's nodes
func AddPhysicsBody(node *Node) *PhysicsBody {

}

//Adds a physics body using the passed in geometry, but updates are transfered back to
// the node for visual updates
func AddPhysicsBodyFromGeometry(node *Node, geometry *Geometry) *PhysicsBody {

}

//AddPhysicsBodyFromNewtonBody allows for creating a specific body in newton directly
// and associating it for engine updates via the passed in node
func AddPhysicsBodyFromNewtonBody(node *Node, body *newton.Body) {

}
