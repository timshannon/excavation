package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"code.google.com/p/gonewton/newton"
	"fmt"
	vmath "github.com/timshannon/vectormath"
)

const (
	GRAVITY         = -9.8
	CONVEXTOLERANCE = 0.1
)

var phWorld *newton.World
var phLastUpdate float32
var phMatrix []float32

type PhysicsScene struct {
	Node *Node
	Body *newton.Body
}

type PhysicsBody struct {
	Node *Node
	Body *newton.Body
}

func initPhysics() {
	phWorld = newton.CreateWorld()
	phMatrix = make([]float32, 16)
}

func PhysicsWorld() *newton.World {
	return phWorld
}

func updatePhysics() {
	phWorld.Update(float32(GameTime()) - phLastUpdate)
}

func NewtonApplyForceAndTorque(body *newton.Body, timestep float32, threadIndex int) {
	var Ixx, Iyy, Izz, mass float32

	body.MassMatrix(&mass, &Ixx, &Iyy, &Izz)
	body.SetForce([]float32{0.0, mass * GRAVITY, 0.0, 1.0})
}

func NewtonTransformUpdate(body *newton.Body, matrix []float32, threadIndex int) {
	body.Matrix(phMatrix)
	//TODO: Translate abs physics matrix to relative matrix?
	//TODO: interpolate visual position from physics

	pBody := body.UserData().(*PhysicsBody)
	//Can only set relative matrix
	horde3d.SetNodeTransMat(pBody.Node.H3DNode, phMatrix)
}

func clearAllPhysics() {
	phWorld.DestroyAllBodies()
}

//Allows me to share face access code between scene trees and regular meshes
type hordeMeshFaceIterator func(face []float32)

//NewtonMeshListFromNode Returns an array of newton meshes from the passed in Nodes child meshes
// each child being a separate mesh
func NewtonMeshListFromNode(node *Node) []*newton.Mesh {
	hMeshes := node.FindChild("", NodeTypeMesh)
	nMeshes := make([]*newton.Mesh, len(hMeshes))
	geom := horde3d.H3DRes(horde3d.GetNodeParamI(node.H3DNode, horde3d.Model_GeoResI))

	for i := range hMeshes {
		nMeshes[i] = phWorld.CreateMesh()

		nMeshes[i].BeginFace()
		iterateFacesInMesh(func(face []float32) {
			nMeshes[i].AddFace(3, face, 3*4, phWorld.DefaultMaterialGroupID())
		}, hMeshes[i].H3DNode, geom)
		nMeshes[i].EndFace()
	}

	return nMeshes
}

func NewtonTreeFromNode(node *Node) *newton.Collision {
	collision := phWorld.CreateTreeCollision(int(node.H3DNode))

	hMeshes := node.FindChild("", NodeTypeMesh)
	geom := horde3d.H3DRes(horde3d.GetNodeParamI(node.H3DNode, horde3d.Model_GeoResI))

	collision.BeginTreeBuild()
	for i := range hMeshes {
		iterateFacesInMesh(func(face []float32) {
			collision.AddTreeFace(3, face, 3*4, phWorld.DefaultMaterialGroupID())
		}, hMeshes[i].H3DNode, geom)
	}
	collision.EndTreeBuild(true)

	return collision
}

func iterateFacesInMesh(iterator hordeMeshFaceIterator, hMesh horde3d.H3DNode, geom horde3d.H3DRes) {
	//mesh
	batchStart := horde3d.GetNodeParamI(hMesh, horde3d.Mesh_BatchStartI)
	batchCount := horde3d.GetNodeParamI(hMesh, horde3d.Mesh_BatchCountI)

	//geom
	isInt16 := horde3d.GetResParamI(geom, horde3d.GeoRes_GeometryElem, 0,
		horde3d.GeoRes_GeoIndices16I)
	vertCount := horde3d.GetResParamI(geom, horde3d.GeoRes_GeometryElem, 0,
		horde3d.GeoRes_GeoVertexCountI)
	indexCount := horde3d.GetResParamI(geom, horde3d.GeoRes_GeometryElem, 0,
		horde3d.GeoRes_GeoIndexCountI)

	//Indices
	var indices16 []uint16
	var indices32 []uint32
	var err error

	if isInt16 == 1 {
		indices16, err = horde3d.MapUint16ResStream(geom, horde3d.GeoRes_GeometryElem, 0,
			horde3d.GeoRes_GeoIndexStream, true, false, indexCount)
	} else {
		indices32, err = horde3d.MapUint32ResStream(geom, horde3d.GeoRes_GeometryElem, 0,
			horde3d.GeoRes_GeoIndexStream, true, false, indexCount)
	}
	horde3d.UnmapResStream(geom)

	if err != nil {
		RaiseError(err)
		return
	}

	//Vertices
	vertices, err := horde3d.MapFloatResStream(geom, horde3d.GeoRes_GeometryElem, 0,
		horde3d.GeoRes_GeoVertPosStream, true, false, vertCount*3)
	horde3d.UnmapResStream(geom)

	if err != nil {
		RaiseError(err)
		return
	}

	face := make([]float32, 16)
	var vIndex1, vIndex2, vIndex3 uint32

	for index := 0; index < batchCount; index += 3 {
		if isInt16 == 1 {
			vIndex1 = uint32(indices16[index+batchStart])
			vIndex2 = uint32(indices16[index+batchStart+1])
			vIndex3 = uint32(indices16[index+batchStart+2])
		} else {
			vIndex1 = indices32[index+batchStart]
			vIndex2 = indices32[index+batchStart+1]
			vIndex3 = indices32[index+batchStart+2]
		}

		//pos
		face[0] = vertices[vIndex1*3]
		face[1] = vertices[vIndex1*3+1]
		face[2] = vertices[vIndex1*3+2]

		face[3] = vertices[vIndex2*3]
		face[4] = vertices[vIndex2*3+1]
		face[5] = vertices[vIndex2*3+2]

		face[6] = vertices[vIndex3*3]
		face[7] = vertices[vIndex3*3+1]
		face[8] = vertices[vIndex3*3+2]

		iterator(face)
	}

}

//AddPhysicsScene adds a scene physics collision type built from
// the passed in node's geometry.  
func AddPhysicsScene(node *Node) *PhysicsScene {
	newScene := new(PhysicsScene)
	newScene.Node = node

	collision := NewtonTreeFromNode(node)

	fmt.Println("Scene: ", collision)

	vmath.M4ToSlice(phMatrix, node.AbsoluteTransMat())
	newScene.Body = phWorld.CreateDynamicBody(collision, phMatrix)

	return newScene
}

//Adds a physics body using the passed in node to determine the collision hull.
// a convex hull will be built from the node's visible geometry. Include all children 
// if more than one submesh is found in the node, or it's children, then the newton
// body will be built as a compound collision
func AddPhysicsBody(node *Node, mass float32) *PhysicsBody {
	newBody := new(PhysicsBody)
	newBody.Node = node

	meshes := NewtonMeshListFromNode(node)

	if len(meshes) == 1 {
		collision := phWorld.CreateConvexHullFromMesh(meshes[0], CONVEXTOLERANCE,
			int(node.H3DNode))
		return AddPhysicsBodyFromCollision(node, collision, mass)
	}

	collision := phWorld.CreateCompoundCollision(int(node.H3DNode))

	collision.CompoundBeginAddRemove()
	for i := range meshes {
		subCollision := phWorld.CreateConvexHullFromMesh(meshes[i], CONVEXTOLERANCE,
			int(node.H3DNode))
		collision.CompoundAddSubCollision(subCollision)
	}
	collision.CompoundEndAddRemove()

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

	vmath.M4ToSlice(phMatrix, node.AbsoluteTransMat())

	fmt.Println("Node: ", node.Name())
	fmt.Println("NodeMat: ", phMatrix)
	fmt.Println("Collision: ", collision)
	fmt.Println("collisionInfo: ", collision.CalculateVolume())
	body := phWorld.CreateDynamicBody(collision, phMatrix)

	collision.CalculateInertialMatrix(inertia, origin)
	body.SetMassMatrix(mass, mass*inertia[0], mass*inertia[1], mass*inertia[2])

	body.SetCentreOfMass(origin)

	body.SetForceAndTorqueCallback(NewtonApplyForceAndTorque)
	body.SetTransformCallback(NewtonTransformUpdate)
	body.SetUserData(newBody)

	newBody.Body = body

	return newBody
}
