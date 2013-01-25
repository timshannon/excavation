package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"code.google.com/p/gonewton/newton"
	vmath "github.com/timshannon/vectormath"
)

const (
	gravity         = -9.8
	convexTolerance = 0.1
)

var newtonWorld *newton.World
var newtonLastUpdate float32
var physicsMatrix []float32

var physicsBodies []*PhysicsBody

type PhysicsScene struct {
	Node *Node
	Body *newton.Body
}

type PhysicsBody struct {
	Node      *Node
	Body      *newton.Body
	destroyed bool
}

func initPhysics() {
	newtonWorld = newton.CreateWorld()
	physicsMatrix = make([]float32, 16)
	physicsBodies = make([]*PhysicsBody, 10)
}

func updatePhysics() {
	newtonWorld.Update(float32(GameTime()) - newtonLastUpdate)
	//TODO: use callback 
	for i := 0; i < len(physicsBodies); i++ {
		if physicsBodies[i].destroyed {
			if len(physicsBodies) > 1 {
				physicsBodies = append(physicsBodies[:i], physicsBodies[i+1:]...)
			} else {
				physicsBodies = physicsBodies[0:0]
				break
			}

		}
		physicsBodies[i].Body.Matrix(physicsMatrix)
		//TODO: Translate abs physics matrix to relative matrix?
		//TODO: interpolate visual position from physics
		horde3d.SetNodeTransMat(physicsBodies[i].Node.H3DNode, physicsMatrix)
	}
}

func applyForceAndTorque(body *newton.Body, timestep float32, threadIndex int) {
	var Ixx, Iyy, Izz, mass float32

	body.MassMatrix(&mass, &Ixx, &Iyy, &Izz)
	body.SetForce([]float32{0.0, mass * gravity, 0.0, 1.0})
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
// a convex hull will be built from the node's visible geometry.  If includeChildren
// a compound shape will be made from the children's nodes as well
func AddPhysicsBody(node *Node, includeChildren bool) *PhysicsBody {
	newBody := new(PhysicsBody)
	newBody.Node = node

	//TODO: Write actual code
	return newBody

}

//Adds a physics body using the passed in geometry, but updates are transfered back to
// the node for visual updates
func AddPhysicsBodyFromGeometry(node *Node, geometry *Geometry, mass float32) *PhysicsBody {
	mesh := NewtonMeshFromGeometry(geometry)
	collision := newtonWorld.CreateConvexHullFromMesh(mesh, convexTolerance, int(node.H3DNode))

	return AddPhysicsBodyFromCollision(node, collision, mass)
}

//AddPhysicsBodyFromCollision allows for creating a specific collision in newton directly
//  and associating it for engine updates via the passed in node
// Also sets up body to have standard forces applied and associates user data to PhysicsBody
func AddPhysicsBodyFromCollision(node *Node, collision *newton.Collision, mass float32) *PhysicsBody {
	newBody := new(PhysicsBody)
	inertia := make([]float32, 3)
	origin := make([]float32, 3)

	newBody.Node = node

	vmath.M4ToSlice(physicsMatrix, node.AbsoluteTransMat())
	body := newtonWorld.CreateDynamicBody(collision, physicsMatrix)

	//TODO: mass calculations
	collision.CalculateInertialMatrix(inertia, origin)
	body.SetMassMatrix(mass, mass*inertia[0], mass*inertia[1], mass*inertia[2])

	body.SetCentreOfMass(origin)

	body.SetForceAndTorqueCallback(applyForceAndTorque)
	//TODO: Cleanup array entry on destroy
	body.SetUserData(newBody)

	newBody.Body = body

	physicsBodies = append(physicsBodies, newBody)
	return newBody
}
