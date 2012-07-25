// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"errors"
	"excavation/math3d"
)

type Node struct {
	horde3d.H3DNode
}

var Root *Node

func init() {
	Root = new(Node)
	Root.H3DNode = horde3d.RootNode
}

//Adds nodes from a SceneGraph resource to the scene.
func AddNodes(parent Node, sceneResource Resource) (*Node, error) {
	node := new(Node)
	node.H3DNode = horde3d.AddNodes(parent.H3DNode, sceneResource.H3DRes)

	if node.H3DNode == 0 {
		return nil, errors.New("Error adding nodes to the scene")
	}

	return node, nil
}

//This function returns the type of a specified scene node.  If the node handle is invalid, 
//the function returns the node type Unknown.
func (n *Node) Type() int { return horde3d.GetNodeType(n.H3DNode) }

//Returns the parent of a scene node.
func (n *Node) Parent() *Node {
	parent := new(Node)
	parent.H3DNode = horde3d.GetNodeParent(n.H3DNode)
	return parent
}

//Relocates a node in the scene graph.
func (n *Node) SetParent(parent Node) bool { return horde3d.SetNodeParent(n.H3DNode, parent.H3DNode) }

//Returns a slice of the children of the current node
func (n *Node) Children() []*Node {
	var hNode horde3d.H3DNode = -1
	var children []*Node
	for i := 0; hNode != 0; i++ {
		hNode = horde3d.GetNodeChild(n.H3DNode, i)
		if hNode != 0 {
			children = append(children, &Node{hNode})
		}
	}
	return children
}

//removes the node from the scene
func (n *Node) Remove() { horde3d.RemoveNode(n.H3DNode) }

//This function checks if a scene node has been transformed by the engine 
//since the last time the transformation flag was reset.  Therefore, it stores 
//a flag that is set to true when a setTransformation function is called 
//explicitely by the application or when the node transformation has been 
//updated by the animation system.  The function also makes it possible to 
//reset the transformation flag.
func (n *Node) CheckTransFlag(reset bool) bool {
	return horde3d.CheckNodeTransFlag(n.H3DNode, reset)
}

//This function gets the translation, rotation and scale of a specified scene node object. 
// The coordinates are in local space and contain the transformation of the node relative to its parent.
func (n *Node) Transform() (translate, rotate, scale math3d.Vector3) {
	translate = math3d.MakeVector3(0, 0, 0)
	rotate = math3d.MakeVector3(0, 0, 0)
	scale = math3d.MakeVector3(0, 0, 0)
	horde3d.GetNodeTransform(n.H3DNode, &translate[0], &translate[1], &translate[2],
		&rotate[0], &rotate[1], &rotate[2], &scale[0], &scale[1], &scale[2])
	return
}

//This function sets the relative translation, rotation and scale of a 
//specified scene node object.  The coordinates are in local space and 
//contain the transformation of the node relative to its parent.
func (n *Node) SetTransform(translate, rotate, scale math3d.Vector3) {
	horde3d.SetNodeTransform(n.H3DNode, translate[0], translate[1], translate[2],
		rotate[0], rotate[1], rotate[2], scale[0], scale[1], scale[2])
}

//Gets the relative transformation matrix of the node
func (n *Node) RelativeTransMat() math3d.Matrix4 {
	relative := math3d.MakeMatrix4()
	horde3d.GetNodeTransMats(n.H3DNode, relative, nil)
	return relative
}

//Gets the absolute transformation matrix of the node
func (n *Node) AbsoluteTransMat() math3d.Matrix4 {
	absolute := math3d.MakeMatrix4()
	horde3d.GetNodeTransMats(n.H3DNode, nil, absolute)
	return absolute
}

//Sets the relative transformation matrix of the node
func (n *Node) SetRelativeTransMat(matrix math3d.Matrix4) {
	horde3d.SetNodeTransMat(n.H3DNode, matrix)
}

//Returns the bounds of a box that encompasses the node
func (n *Node) BoundingBox() (min, max math3d.Vector3) {
	min = math3d.MakeVector3(0, 0, 0)
	max = math3d.MakeVector3(0, 0, 0)
	horde3d.GetNodeAABB(n.H3DNode, &min[0], &min[1], &min[2],
		&max[0], &max[1], &max[2])
	return
}

//FindChild: This function loops recursively over all children of startNode and adds 
//them to an internal list of results if they match the specified name and type.  
//The result list is cleared each time this function is called.  
//The function returns the number of nodes which were found and added to the list.

//Parameters
//name name of nodes to be searched (empty string for all nodes) 
//nodeType type of nodes to be searched (NodeTypes_Undefined for all types) 
func (n *Node) FindChild(name string, nodeType int) []*Node {
	size := horde3d.FindNodes(n.H3DNode, name, nodeType)
	results := make([]*Node, size)

	for i := range results {
		results[i].H3DNode = horde3d.GetNodeFindResult(i)
	}
	return results
}

//Excludes scene node from all rendering
func (n *Node) NoDraw() bool { return horde3d.NodeFlags_NoDraw == horde3d.GetNodeFlags(n.H3DNode) }

//Excludes scene node from list of shadow casters
func (n *Node) NoCastShadow() bool {
	return horde3d.NodeFlags_NoCastShadow == horde3d.GetNodeFlags(n.H3DNode)
}

//Excludes scene node from ray intersection queries
func (n *Node) NoRayQuery() bool {
	return horde3d.NodeFlags_NoRayQuery == horde3d.GetNodeFlags(n.H3DNode)
}

//Deactivates scene node so that it is completely ignored (combination of all flags above)
func (n *Node) Inactive() bool { return horde3d.NodeFlags_Inactive == horde3d.GetNodeFlags(n.H3DNode) }

func (n *Node) SetNoDraw(recursive bool) {
	horde3d.SetNodeFlags(n.H3DNode, horde3d.NodeFlags_NoDraw, recursive)
}
func (n *Node) SetNoCastShadow(recursive bool) {
	horde3d.SetNodeFlags(n.H3DNode, horde3d.NodeFlags_NoCastShadow, recursive)
}
func (n *Node) SetNoRayQuery(recursive bool) {
	horde3d.SetNodeFlags(n.H3DNode, horde3d.NodeFlags_NoRayQuery, recursive)
}
func (n *Node) SetInactive(recursive bool) {
	horde3d.SetNodeFlags(n.H3DNode, horde3d.NodeFlags_Inactive, recursive)
}

//Gets the name of the node
func (n *Node) Name() string {
	return horde3d.GetNodeParamStr(n.H3DNode, horde3d.NodeParams_NameStr)
}

//SetName sets the name of the node
func (n *Node) SetName(name string) {
	horde3d.SetNodeParamStr(n.H3DNode, horde3d.NodeParams_NameStr, name)
}

//Optional application-specific meta data for a node encapsulated in an Attachment XML string 
func (n *Node) Attachment() string {
	return horde3d.GetNodeParamStr(n.H3DNode, horde3d.NodeParams_AttachmentStr)
}

//Optional application-specific meta data for a node encapsulated in an Attachment XML string 
func (n *Node) SetAttachment(value string) {
	horde3d.SetNodeParamStr(n.H3DNode, horde3d.NodeParams_AttachmentStr, value)
}

//Returns true if both nodes refer to the same internal node
func (n *Node) IsSame(other *Node) bool {
	return n.H3DNode == other.H3DNode
}

type CastRayResult struct {
	ResultNode   *Node
	Distance     *float32
	Intersection math3d.Vector3
}

//This function checks recursively if the specified ray intersects the specified node or one of its children.  
//The function finds intersections relative to the ray origin and returns the number of intersecting scene nodes.  
//The ray is a line segment and is specified by a starting point (the origin) and a finite direction vector 
//which also defines its length.  Currently this function is limited to returning intersections with Meshes.  
//For Meshes, the base LOD (LOD0) is always used for performing the ray-triangle intersection tests.
func (n *Node) CastRay(origin, direction math3d.Vector3, maxNearest int) []*CastRayResult {
	size := horde3d.CastRay(n.H3DNode, origin[0], origin[1], origin[2],
		direction[0], direction[1], direction[2], maxNearest)

	results := make([]*CastRayResult, size)

	for i := range results {
		intersection := math3d.MakeVector3(0, 0, 0)
		_ = horde3d.GetCastRayResult(i, &results[i].ResultNode.H3DNode, results[i].Distance,
			intersection)
		results[i].Intersection = intersection
	}
	return results
}

//This function checks if a specified node is visible from the perspective of a specified camera.  
//The function always checks if the node is in the camera.s frustum.  If checkOcclusion is true, 
//the function will take into account the occlusion culling information from the previous frame 
//(if occlusion culling is disabled the flag is ignored).  The flag calcLod determines whether the 
//detail level for the node should be returned in case it is visible.  The function returns -1 if 
//the node is not visible, otherwise 0 (base LOD level) or the computed LOD level
func (n *Node) IsVisible(camera *Camera, checkOcclusion, calcLOD bool) int {
	return horde3d.CheckNodeVisibility(n.H3DNode, camera.H3DNode, checkOcclusion, calcLOD)
}

type Group struct{ *Node }

//Adds a new group node
func AddGroup(parent *Node, name string) (*Group, error) {
	group := new(Group)
	group.H3DNode = horde3d.AddGroupNode(parent.H3DNode, name)
	if group.H3DNode == 0 {
		return nil, errors.New("Error adding group node")
	}
	return group, nil
}

type Model struct{ *Node }

//Adds a new model
func AddModel(parent *Node, name string, geometry *Geometry) (*Model, error) {
	model := new(Model)
	model.H3DNode = horde3d.AddModelNode(parent.H3DNode, name, geometry.H3DRes)
	if model.H3DNode == 0 {
		return nil, errors.New("Error adding Model")
	}
	return model, nil
}

//Gets the Geometry resource for the given model
func (m *Model) Geometry() *Geometry {
	geom := new(Geometry)
	geom.H3DRes = horde3d.H3DRes(horde3d.GetNodeParamI(m.H3DNode, horde3d.Model_GeoResI))
	return geom
}

//Sets the Geometry resource for the given model
func (m *Model) SetGeometry(newGeom *Geometry) {
	horde3d.SetNodeParamI(m.H3DNode, horde3d.Model_GeoResI, int(newGeom.H3DRes))
}

//Gets state of software skinning
func (m *Model) SWSkinning() int {
	return horde3d.GetNodeParamI(m.H3DNode, horde3d.Model_SWSkinningI)
}

//Sets software skinning
func (m *Model) SetSWSkinning(value int) {
	horde3d.SetNodeParamI(m.H3DNode, horde3d.Model_SWSkinningI, value)
}

//Gets the distances for the LevelOfDetail settings
func (m *Model) LODDist() (LOD1, LOD2, LOD3, LOD4 float32) {
	LOD1 = horde3d.GetNodeParamF(m.H3DNode, horde3d.Model_LodDist1F, 0)
	LOD2 = horde3d.GetNodeParamF(m.H3DNode, horde3d.Model_LodDist2F, 0)
	LOD3 = horde3d.GetNodeParamF(m.H3DNode, horde3d.Model_LodDist3F, 0)
	LOD4 = horde3d.GetNodeParamF(m.H3DNode, horde3d.Model_LodDist4F, 0)
	return
}

//Sets the distances for the LevelOfDetail settings
// subsequent LODs must be greater than the previous i.e. LOD1 < LOD2
func (m *Model) SetLODDist(LOD1, LOD2, LOD3, LOD4 float32) {
	horde3d.SetNodeParamF(m.H3DNode, horde3d.Model_LodDist1F, 0, LOD1)
	horde3d.SetNodeParamF(m.H3DNode, horde3d.Model_LodDist2F, 0, LOD2)
	horde3d.SetNodeParamF(m.H3DNode, horde3d.Model_LodDist3F, 0, LOD3)
	horde3d.SetNodeParamF(m.H3DNode, horde3d.Model_LodDist4F, 0, LOD4)
}

func (m *Model) SetupAnimStage(stage int, animation *Animation, layer int,
	startNode string, additive bool) {
	horde3d.SetupModelAnimStage(m.H3DNode, stage, animation.H3DRes, layer, startNode, additive)
}

func (m *Model) SetAnimParams(stage int, time, weight float32) {
	horde3d.SetModelAnimParams(m.H3DNode, stage, time, weight)
}

func (m *Model) SetMorpher(target string, weight float32) bool {
	return horde3d.SetModelMorpher(m.H3DNode, target, weight)
}

type Mesh struct{ *Node }

func AddMesh(parent *Node, name string, material *Material, batchStart, batchCount,
	vertRStart, vertREnd int) (*Mesh, error) {
	mesh := new(Mesh)
	mesh.H3DNode = horde3d.AddMeshNode(parent.H3DNode, name, material.H3DRes, batchStart,
		batchCount, vertRStart, vertREnd)
	if mesh.H3DNode == 0 {
		return nil, errors.New("Error adding Mesh")
	}
	return mesh, nil
}

func (m *Mesh) Material() *Material {
	material := new(Material)
	material.H3DRes = horde3d.H3DRes(horde3d.GetNodeParamI(m.H3DNode, horde3d.Mesh_MatResI))
	return material
}

func (m *Mesh) SetMaterial(newMaterial *Material) {
	horde3d.SetNodeParamI(m.H3DNode, horde3d.Mesh_MatResI, int(newMaterial.H3DRes))
}

func (m *Mesh) BatchStart() int { return horde3d.GetNodeParamI(m.H3DNode, horde3d.Mesh_BatchStartI) }
func (m *Mesh) BatchCount() int { return horde3d.GetNodeParamI(m.H3DNode, horde3d.Mesh_BatchCountI) }
func (m *Mesh) VertRStart() int { return horde3d.GetNodeParamI(m.H3DNode, horde3d.Mesh_VertRStartI) }
func (m *Mesh) VertREnd() int   { return horde3d.GetNodeParamI(m.H3DNode, horde3d.Mesh_VertREndI) }

func (m *Mesh) LODLevel() int { return horde3d.GetNodeParamI(m.H3DNode, horde3d.Mesh_LodLevelI) }
func (m *Mesh) SetLODLevel(level int) {
	horde3d.SetNodeParamI(m.H3DNode, horde3d.Mesh_LodLevelI, level)
}

type Joint struct{ *Node }

func AddJoint(parent *Node, name string, jointIndex int) (*Joint, error) {
	joint := new(Joint)
	joint.H3DNode = horde3d.AddJointNode(parent.H3DNode, name, jointIndex)
	if joint.H3DNode == 0 {
		return nil, errors.New("Error adding Joint")
	}
	return joint, nil
}

type Light struct{ *Node }
type Camera struct{ *Node }
type Emitter struct{ *Node }
