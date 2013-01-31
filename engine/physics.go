package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"code.google.com/p/gonewton/newton"
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

//NewtonMeshListFromNode Returns an array of newton meshes from the passed in Nodes child meshes
// each child being a separate mesh
func NewtonMeshListFromNode(node *Node) []*newton.Mesh {
	hMeshes := node.FindChild("", NodeTypeMesh)
	nMeshes := make([]*newton.Mesh, len(hMeshes))

	for i := range hMeshes {
		nMeshes[i] = phWorld.CreateMesh()
		AddMeshNodeToNewtonMesh(nMeshes[i], hMeshes[i].H3DNode)
	}
	return nMeshes
}

//NewtonmeshFromNode returns a single mesh built from all of the child mesh
// nodes in the passed in node
func NewtonMeshFromNode(node *Node) *newton.Mesh {
	mesh := phWorld.CreateMesh()
	hMeshes := node.FindChild("", NodeTypeMesh)

	for i := range hMeshes {
		AddMeshNodeToNewtonMesh(mesh, hMeshes[i].H3DNode)
	}

	return mesh
}

//AddNewtonMeshToMesh adds the passed in Mesh resource to the passed in
// Newton Mesh
func AddMeshNodeToNewtonMesh(newtonMesh *newton.Mesh, hNode horde3d.H3DNode) {
	//mesh
	//vertRStart := horde3d.GetNodeParamI(hNode, horde3d.Mesh_VertRStartI)
	//vertREnd := horde3d.GetNodeParamI(hNode, horde3d.Mesh_VertREndI)
	//batchStart := horde3d.GetNodeParamI(hNode, horde3d.Mesh_BatchStartI)
	//batchCount := horde3d.GetNodeParamI(hNode, horde3d.Mesh_BatchCountI)
	//geom
	geom := horde3d.H3DRes(horde3d.GetNodeParamI(hNode, horde3d.Model_GeoResI))
	isInt16 := horde3d.GetResParamI(geom, horde3d.GeoRes_GeometryElem, 0,
		horde3d.GeoRes_GeoIndices16I)
	vertCount := horde3d.GetResParamI(geom, horde3d.GeoRes_GeometryElem, 0,
		horde3d.GeoRes_GeoVertexCountI)
	indexCount := horde3d.GetResParamI(geom, horde3d.GeoRes_GeometryElem, 0,
		horde3d.GeoRes_GeoIndexCountI)

	var byteSize int
	if isInt16 == 0 {
		byteSize = 4
	} else {
		byteSize = 2
	}

	//Indices
	var indices = make([]byte, indexCount*byteSize)
	indexStream, err := horde3d.MapByteResStream(geom, horde3d.GeoRes_GeometryElem, 0, horde3d.GeoRes_GeoIndexStream,
		true, false, indexCount*byteSize)

	copy(indices, indexStream)
	horde3d.UnmapResStream(geom)

	if err != nil {
		RaiseError(err)
		return
	}
	//Vertices
	vertices := make([]float32, vertCount*3)
	vertStream, err := horde3d.MapFloatResStream(geom, horde3d.GeoRes_GeometryElem, 0, horde3d.GeoRes_GeoVertPosStream,
		true, false, vertCount*3)
	copy(vertices, vertStream)
	horde3d.UnmapResStream(geom)

	if err != nil {
		RaiseError(err)
		return
	}

	//Tangents and Normals
	tangent := make([]float32, vertCount*7)
	tanStream, err := horde3d.MapFloatResStream(geom, horde3d.GeoRes_GeometryElem, 0, horde3d.GeoRes_GeoVertTanStream,
		true, false, vertCount*7)
	copy(tangent, tanStream)
	horde3d.UnmapResStream(geom)

	if err != nil {
		RaiseError(err)
		return
	}

	//face := make([]float32, 

	//TODO: Write actual code
	//for (int i = 0; i < vertexCount; i ++) {
	//dVector p1 (faceVertec[i * 3 + 0], faceVertec[i * 3 + 1], faceVertec[i * 3 + 2]);
	//p1 += displacemnet;

	//face[i][0] = p1.m_x; 
	//face[i][1] = p1.m_y;  
	//face[i][2] = p1.m_z;   

	//face[i][3] = normal.m_x; 
	//face[i][4] = normal.m_y;  
	//face[i][5] = normal.m_z;  

	//face[i][6] = 0.0f; 
	//face[i][7] = 0.0f;  
	//face[i][8] = 0.0f;  
	//face[i][9] = 0.0f;  
	//}

	//// add the face
	//NewtonMeshAddFace (mesh, vertexCount, &face[0][0], 10 * sizeof (float), id);

	//import "encoding/binary"
	// binary.Size(0.0)?

}

//AddPhysicsScene adds a scene physics collision type built from
// the passed in node's geometry.  
func AddPhysicsScene(node *Node) *PhysicsScene {
	newScene := new(PhysicsScene)
	newScene.Node = node

	mesh := NewtonMeshFromNode(node)

	collision := phWorld.CreateTreeCollsionFromMesh(mesh, int(node.H3DNode))

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

	if len(meshes) == 0 {
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
