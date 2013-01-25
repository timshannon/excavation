package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"code.google.com/p/gonewton/newton"
	vmath "github.com/timshannon/vectormath"
)

const GRAVITY = -9.8

var newtonWorld *newton.World
var newtonLastUpdate float32
var physicsMatrix []float32

var physicsBodies []*PhysicsBody

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
	physicsMatrix = make([]float32, 16)
}

func updatePhysics() {
	newtonWorld.Update(float32(GameTime()) - newtonLastUpdate)
	for i := range physicsBodies {
		physicsBodies[i].Body.Matrix(physicsMatrix)
		//TODO: Translate abs physics matrix to relative matrix?
		//TODO: interpolate visual position from physics
		horde3d.SetNodeTransMat(physicsBodies[i].Node.H3DNode, physicsMatrix)
	}
}

func applyForceAndTorque(body *newton.Body, timestep float32, threadIndex int) {
	var Ixx, Iyy, Izz, mass float32

	body.MassMatrix(&mass, &Ixx, &Iyy, &Izz)
	body.SetForce([]float32{0.0, mass * GRAVITY, 0.0, 1.0})
}

func clearAllPhysics() {
	newtonWorld.DestroyAllBodies()
}

//Returns a newton mesh from the passed in geometry resource
func NewtonMeshFromGeometry(geom *Geometry) *newton.Mesh {
	mesh := new(newton.Mesh)

	//TODO: Write actual code
	return mesh
}

//AddPhysicsScene adds a scene physics collision type built from
// the passed in node's geometry.  
func AddPhysicsScene(node *Node) *PhysicsScene {
	newScene := new(PhysicsScene)
	newScene.Node = node

	mesh := NewtonMeshFromGeometry(node.Geometry())

	collision := newtonWorld.CreateTreeCollsionFromMesh(mesh, int(node.H3DNode))

	vmath.M4ToSlice(physicsMatrix, node.AbsoluteTransMat())
	newScene.Body = newtonWorld.CreateDynamicBody(collision, physicsMatrix)

	return newScene
}

//AddPhysicsSceneFromGeometry adds a scene physics collision type built
// from the passed in geometry resource
func AddPhysicsSceneFromGeometry(node *Node, geometry *Geometry) *PhysicsScene {
	newScene := new(PhysicsScene)
	newScene.Node = node

	mesh := NewtonMeshFromGeometry(geometry)

	collision := newtonWorld.CreateTreeCollsionFromMesh(mesh, int(node.H3DNode))

	vmath.M4ToSlice(physicsMatrix, node.AbsoluteTransMat())
	newScene.Body = newtonWorld.CreateDynamicBody(collision, physicsMatrix)

	return newScene
}

//Adds a physics body using the passed in node to determine the collision hull.
// a convex hull will be built from the node's visible geometry.  If the node has
// children a compound shape will be made from the children's nodes
func AddPhysicsBody(node *Node) *PhysicsBody {
	newBody := new(PhysicsBody)
	newBody.Node = node

	//TODO: Write actual code
	return newBody

}

//Adds a physics body using the passed in geometry, but updates are transfered back to
// the node for visual updates
func AddPhysicsBodyFromGeometry(node *Node, geometry *Geometry) *PhysicsBody {
	newBody := new(PhysicsBody)
	newBody.Node = node

	//TODO: Write actual code
	return newBody

}

//AddPhysicsBodyFromNewtonBody allows for creating a specific body in newton directly
// with any collision primatives and associating it for engine updates via the passed in node
func AddPhysicsBodyFromNewtonBody(node *Node, body *newton.Body) *PhysicsBody {
	newBody := new(PhysicsBody)
	newBody.Node = node

	body.SetForceAndTorqueCallback(applyForceAndTorque)
	newBody.Body = body

	physicsBodies = append(physicsBodies, newBody)
	return newBody

}
