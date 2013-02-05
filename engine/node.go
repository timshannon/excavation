// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"errors"
	vmath "github.com/timshannon/vectormath"
)

const (
	NodeTypeUndefined = horde3d.NodeTypes_Undefined
	NodeTypeGroup     = horde3d.NodeTypes_Group
	NodeTypeModel     = horde3d.NodeTypes_Model
	NodeTypeMesh      = horde3d.NodeTypes_Mesh
	NodeTypeJoint     = horde3d.NodeTypes_Joint
	NodeTypeLight     = horde3d.NodeTypes_Light
	NodeTypeCamera    = horde3d.NodeTypes_Camera
	NodeTypeEmitter   = horde3d.NodeTypes_Emitter
)

//temp variables used to keep the GC from thrashing
// by retaining temporary memory space for the conversion from horde's float arrays
// to vectormath's structs, we should be able to keep gc collection to a minimum
// Multiple threads shouldnt' be an issue, as Horde3d is singlethreaded. 
// This may need to change in the future.
var tempRelMat = make([]float32, 16)
var tempAbsMat = make([]float32, 16)
var tempHordeVector = make([]float32, 3)

type Node struct {
	horde3d.H3DNode
	relMat      *vmath.Matrix4
	absMat      *vmath.Matrix4
	updateFrame int
}

func NewNode(hordeNode horde3d.H3DNode) *Node {
	node := &Node{
		horde3d.H3DNode: hordeNode,
		relMat:          new(vmath.Matrix4),
		absMat:          new(vmath.Matrix4),
	}

	return node
}

//Adds nodes from a SceneGraph resource to the scene.
func AddNodes(parent *Node, sceneResource *Scene) (*Node, error) {
	node := NewNode(horde3d.AddNodes(parent.H3DNode, sceneResource.H3DRes))

	if node.H3DNode == 0 {
		return nil, errors.New("Error adding nodes to the scene")
	}

	return node, nil
}

//This function returns the type of a specified scene node.  If the node handle is invalid, 
//the function returns the node type Unknown.
func (n *Node) Type() int {
	intType := horde3d.GetNodeType(n.H3DNode)
	return intType
}

//Returns the parent of a scene node.
func (n *Node) Parent() *Node {
	parent := NewNode(horde3d.GetNodeParent(n.H3DNode))
	return parent
}

//Relocates a node in the scene graph.
func (n *Node) SetParent(parent *Node) bool { return horde3d.SetNodeParent(n.H3DNode, parent.H3DNode) }

//Returns a slice of the children of the current node
func (n *Node) Children() []*Node {
	var hNode horde3d.H3DNode = -1
	var children []*Node
	for i := 0; hNode != 0; i++ {
		hNode = horde3d.GetNodeChild(n.H3DNode, i)
		if hNode != 0 {
			children = append(children, NewNode(hNode))
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
func (n *Node) Transform(translate, rotate, scale *vmath.Vector3) {
	var tx, ty, tz float32
	var rx, ry, rz float32
	var sx, sy, sz float32
	horde3d.GetNodeTransform(n.H3DNode, &tx, &ty, &tz,
		&rx, &ry, &rz, &sx, &sy, &sz)

	vmath.V3MakeFromElems(translate, tx, ty, tz)
	vmath.V3MakeFromElems(rotate, rx, ry, rz)
	vmath.V3MakeFromElems(scale, sx, sy, sz)
}

func (n *Node) Translate(result *vmath.Vector3) {
	var tx, ty, tz float32

	horde3d.GetNodeTransform(n.H3DNode, &tx, &ty, &tz,
		nil, nil, nil, nil, nil, nil)
	vmath.V3MakeFromElems(result, tx, ty, tz)
}

func (n *Node) Rotate(result *vmath.Vector3) {
	var rx, ry, rz float32

	horde3d.GetNodeTransform(n.H3DNode, nil, nil, nil,
		&rx, &ry, &rz, nil, nil, nil)
	vmath.V3MakeFromElems(result, rx, ry, rz)
}

func (n *Node) Scale(result *vmath.Vector3) {
	var sx, sy, sz float32

	horde3d.GetNodeTransform(n.H3DNode, nil, nil, nil,
		nil, nil, nil, &sx, &sy, &sz)
	vmath.V3MakeFromElems(result, sx, sy, sz)
}

func (n *Node) Occluded() bool {
	//TODO: Cast physics ray from node to camera?
	return false
}

//This function sets the relative translation, rotation and scale of a 
//specified scene node object.  The coordinates are in local space and 
//contain the transformation of the node relative to its parent.
func (n *Node) SetTransform(translate, rotate, scale *vmath.Vector3) {
	horde3d.SetNodeTransform(n.H3DNode, translate.X, translate.Y, translate.Z,
		rotate.X, rotate.Y, rotate.Z,
		scale.X, scale.Y, scale.Z)
}

func (n *Node) updateTransMats() {
	//only jump into cgo code if the matrices haven't
	// been updated for this frame
	if n.updateFrame != frames {
		horde3d.GetNodeTransMats(n.H3DNode, tempRelMat, tempAbsMat)
		vmath.SliceToM4(n.relMat, tempRelMat)
		vmath.SliceToM4(n.absMat, tempAbsMat)
		n.updateFrame = frames
	}
}

//Gets the relative transformation matrix of the node
func (n *Node) RelativeTransMat() *vmath.Matrix4 {
	n.updateTransMats()
	return n.relMat
}

//Gets the absolute transformation matrix of the node
func (n *Node) AbsoluteTransMat() *vmath.Matrix4 {
	n.updateTransMats()
	return n.absMat
}

//Sets the relative transformation matrix of the node
func (n *Node) SetRelativeTransMat(matrix *vmath.Matrix4) {
	//reset update frame so that changes to local matrix
	// will be refreshed from c code
	n.updateFrame = 0
	vmath.M4ToSlice(tempRelMat, matrix)
	horde3d.SetNodeTransMat(n.H3DNode, tempRelMat)
}

func (n *Node) SetLocalTransform(translate, rotate *vmath.Vector3) {
	//set transform relative to itself
	n.SetTransformRelativeTo(n, translate, rotate)
}

//SetTransformRelativeTo sets the transform relative 
// to another node's position and rotation
// Note this function creates a lot of temp variables, and may
// cause GC collection performance issues
func (n *Node) SetTransformRelativeTo(otherNode *Node,
	trans, rotate *vmath.Vector3) {

	var matrix *vmath.Matrix4
	transform := new(vmath.Transform3)
	m3 := new(vmath.Matrix3)
	translate := new(vmath.Vector3)
	newTranslate := new(vmath.Vector3)
	rotM3 := new(vmath.Matrix3)

	vmath.M3MakeRotationZYX(rotM3, rotate)
	vmath.V3Copy(newTranslate, trans)

	matrix = otherNode.RelativeTransMat()
	vmath.M4GetTranslation(translate, matrix)
	vmath.M4GetUpper3x3(m3, matrix)

	transform.SetUpper3x3(m3)
	vmath.T3MulV3(newTranslate, transform, newTranslate)
	vmath.M3Mul(rotM3, m3, rotM3)

	vmath.V3Add(newTranslate, translate, newTranslate)

	vmath.M4MakeFromM3V3(matrix, rotM3, newTranslate)
	n.SetRelativeTransMat(matrix)

}

//Returns the bounds of a box that encompasses the node
func (n *Node) BoundingBox(min, max *vmath.Vector3) {
	var minX, minY, minZ, maxX, maxY, maxZ float32
	horde3d.GetNodeAABB(n.H3DNode, &minX, &minY, &minZ, &maxX, &maxY, &maxZ)

	min.X = minX
	min.Y = minY
	min.Z = minZ
	max.X = maxX
	max.Y = maxY
	max.Z = maxZ
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
		results[i] = NewNode(horde3d.GetNodeFindResult(i))
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
	Intersection *vmath.Vector3
}

//This function checks recursively if the specified ray intersects the specified node or one of its children.  
//The function finds intersections relative to the ray origin and returns the number of intersecting scene nodes.  
//The ray is a line segment and is specified by a starting point (the origin) and a finite direction vector 
//which also defines its length.  Currently this function is limited to returning intersections with Meshes.  
//For Meshes, the base LOD (LOD0) is always used for performing the ray-triangle intersection tests.
func (n *Node) CastRay(results []*CastRayResult, origin, direction *vmath.Vector3) {
	size := horde3d.CastRay(n.H3DNode, origin.X, origin.Y, origin.Z,
		direction.X, direction.Y, direction.Z, len(results))

	results = results[:size]
	for i := range results {
		results[i].ResultNode = NewNode(0)
		_ = horde3d.GetCastRayResult(i, &results[i].ResultNode.H3DNode, results[i].Distance,
			tempHordeVector)

		newVec := &vmath.Vector3{}
		vmath.SliceToV3(newVec, tempHordeVector)
		results[i].Intersection = newVec
	}
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
	group := &Group{NewNode(horde3d.AddGroupNode(parent.H3DNode, name))}
	if group.H3DNode == 0 {
		return nil, errors.New("Error adding group node")
	}
	return group, nil
}

type Model struct{ *Node }

//Adds a new model
func AddModel(parent *Node, name string, geometry *Geometry) (*Model, error) {
	model := &Model{NewNode(horde3d.AddModelNode(parent.H3DNode, name, geometry.H3DRes))}
	if model.H3DNode == 0 {
		return nil, errors.New("Error adding Model")
	}
	return model, nil
}

//Gets the Geometry resource for the given model
func (m *Model) Geometry() *Geometry {
	geom := &Geometry{new(Resource)}
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
	mesh := &Mesh{NewNode(horde3d.AddMeshNode(parent.H3DNode, name, material.H3DRes, batchStart,
		batchCount, vertRStart, vertREnd))}
	if mesh.H3DNode == 0 {
		return nil, errors.New("Error adding Mesh")
	}
	return mesh, nil
}

func (m *Mesh) Material() *Material {
	material := &Material{new(Resource)}
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
	joint := &Joint{NewNode(horde3d.AddJointNode(parent.H3DNode, name, jointIndex))}
	if joint.H3DNode == 0 {
		return nil, errors.New("Error adding Joint")
	}
	return joint, nil
}

func (j *Joint) Index() int { return horde3d.GetNodeParamI(j.H3DNode, horde3d.Joint_JointIndexI) }

type Light struct{ *Node }

func AddLight(parent *Node, name string, material *Material, lightingContext string,
	shadowContext string) *Light {
	light := &Light{NewNode(horde3d.AddLightNode(parent.H3DNode, name, material.H3DRes,
		lightingContext, shadowContext))}
	return light
}

func (l *Light) Material() *Material {
	material := &Material{new(Resource)}
	material.H3DRes = horde3d.H3DRes(horde3d.GetNodeParamI(l.H3DNode, horde3d.Light_MatResI))
	return material
}

func (l *Light) SetMaterial(material *Material) {
	horde3d.SetNodeParamI(l.H3DNode, horde3d.Light_MatResI, int(material.H3DRes))
}

func (l *Light) FOV() float32 { return horde3d.GetNodeParamF(l.H3DNode, horde3d.Light_FovF, 0) }
func (l *Light) SetFOV(newFOV float32) {
	horde3d.SetNodeParamF(l.H3DNode, horde3d.Light_FovF, 0, newFOV)
}

func (l *Light) Color() (r, g, b float32) {
	r = horde3d.GetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 0)
	b = horde3d.GetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 1)
	g = horde3d.GetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 2)
	return
}

func (l *Light) SetColor(r, g, b float32) {
	horde3d.SetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 0, r)
	horde3d.SetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 1, g)
	horde3d.SetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 2, b)
}

func (l *Light) ColorMultiplier() float32 {
	return horde3d.GetNodeParamF(l.H3DNode, horde3d.Light_ColorMultiplierF, 0)
}

func (l *Light) SetColorMultiplier(multiplier float32) {
	horde3d.SetNodeParamF(l.H3DNode, horde3d.Light_ColorMultiplierF, 0, multiplier)
}

func (l *Light) ShadowMapCount() int {
	return horde3d.GetNodeParamI(l.H3DNode, horde3d.Light_ShadowMapCountI)
}

func (l *Light) SetShadowMapCount(count int) {
	horde3d.SetNodeParamI(l.H3DNode, horde3d.Light_ShadowMapCountI, count)
}

func (l *Light) ShadowSplitLambda() float32 {
	return horde3d.GetNodeParamF(l.H3DNode, horde3d.Light_ShadowSplitLambdaF, 0)
}

func (l *Light) SetShadowSplitLambda(lambda float32) {
	horde3d.SetNodeParamF(l.H3DNode, horde3d.Light_ShadowSplitLambdaF, 0, lambda)
}

func (l *Light) ShadowMapBias() float32 {
	return horde3d.GetNodeParamF(l.H3DNode, horde3d.Light_ShadowMapBiasF, 0)
}

func (l *Light) SetShadowMapBias(bias float32) {
	horde3d.SetNodeParamF(l.H3DNode, horde3d.Light_ShadowMapBiasF, 0, bias)
}

func (l *Light) LightingContext() string {
	return horde3d.GetNodeParamStr(l.H3DNode, horde3d.Light_LightingContextStr)
}

func (l *Light) SetLightingContext(context string) {
	horde3d.SetNodeParamStr(l.H3DNode, horde3d.Light_LightingContextStr, context)
}

func (l *Light) ShadowContext() string {
	return horde3d.GetNodeParamStr(l.H3DNode, horde3d.Light_ShadowContextStr)
}

func (l *Light) SetShadowContext(context string) {
	horde3d.SetNodeParamStr(l.H3DNode, horde3d.Light_ShadowContextStr, context)
}

type Camera struct{ *Node }

func AddCamera(parent *Node, name string, pipeline *Pipeline) *Camera {
	camera := &Camera{NewNode(horde3d.AddCameraNode(parent.H3DNode, name, pipeline.H3DRes))}
	return camera
}

func (c *Camera) SetupView(FOV, aspect, nearDist, farDist float32) {
	horde3d.SetupCameraView(c.H3DNode, FOV, aspect, nearDist, farDist)
}

func (c *Camera) ProjectionMatrix(result *vmath.Matrix4) {
	horde3d.GetCameraProjMat(c.H3DNode, tempRelMat)
	vmath.SliceToM4(result, tempRelMat)
}

func (c *Camera) Pipeline() *Pipeline {
	pipeline := &Pipeline{new(Resource)}
	pipeline.H3DRes = horde3d.H3DRes(horde3d.GetNodeParamI(c.H3DNode, horde3d.Camera_PipeResI))
	return pipeline
}

func (c *Camera) SetPipeline(pipeline *Pipeline) {
	horde3d.SetNodeParamI(c.H3DNode, horde3d.Camera_PipeResI, int(pipeline.H3DRes))
}

//2D Texture resource used as output buffer (can be 0 to use main framebuffer) (default: 0)
func (c *Camera) OutTexture() *Texture {
	texture := &Texture{new(Resource)}
	texture.H3DRes = horde3d.H3DRes(horde3d.GetNodeParamI(c.H3DNode, horde3d.Camera_OutTexResI))
	return texture
}

func (c *Camera) SetOutTexture(texture *Texture) {
	horde3d.SetNodeParamI(c.H3DNode, horde3d.Camera_OutTexResI, int(texture.H3DRes))
}

//Index of the output buffer for stereo rendering (values: 0 for left eye, 1 for right eye) (default: 0)
func (c *Camera) OutputBufferIndex() int {
	return horde3d.GetNodeParamI(c.H3DNode, horde3d.Camera_OutBufIndexI)
}

func (c *Camera) SetOutputBufferIndex(index int) {
	horde3d.SetNodeParamI(c.H3DNode, horde3d.Camera_OutBufIndexI, index)
}

func (c *Camera) Viewport() (x, y, width, height int) {
	x = horde3d.GetNodeParamI(c.H3DNode, horde3d.Camera_ViewportXI)
	y = horde3d.GetNodeParamI(c.H3DNode, horde3d.Camera_ViewportYI)
	width = horde3d.GetNodeParamI(c.H3DNode, horde3d.Camera_ViewportWidthI)
	height = horde3d.GetNodeParamI(c.H3DNode, horde3d.Camera_ViewportHeightI)
	return
}

func (c *Camera) SetViewport(x, y, width, height int) {
	horde3d.SetNodeParamI(c.H3DNode, horde3d.Camera_ViewportXI, x)
	horde3d.SetNodeParamI(c.H3DNode, horde3d.Camera_ViewportYI, y)
	horde3d.SetNodeParamI(c.H3DNode, horde3d.Camera_ViewportWidthI, width)
	horde3d.SetNodeParamI(c.H3DNode, horde3d.Camera_ViewportHeightI, height)
}

func (c *Camera) IsOrtho() bool {
	i := horde3d.GetNodeParamI(c.H3DNode, horde3d.Camera_OrthoI)
	if i != 0 {
		return true
	}
	return false
}

func (c *Camera) SetOrtho(value bool) {
	var i int
	if value {
		i = 1
	} else {
		i = 0
	}
	horde3d.SetNodeParamI(c.H3DNode, horde3d.Camera_OrthoI, i)
}

func (c *Camera) OcclusionCulling() bool {
	i := horde3d.GetNodeParamI(c.H3DNode, horde3d.Camera_OccCullingI)

	if i != 0 {
		return true
	}
	return false
}

func (c *Camera) SetOcclusionCulling(value bool) {
	var i int
	if value {
		i = 1
	} else {
		i = 0
	}

	horde3d.SetNodeParamI(c.H3DNode, horde3d.Camera_OccCullingI, i)
}

type Emitter struct{ *Node }

func AddEmitter(parent *Node, name string, material *Material, particleEffect *ParticleEffect,
	maxParticleCount int, respawnCount int) *Emitter {
	emitter := &Emitter{NewNode(horde3d.AddEmitterNode(parent.H3DNode, name, material.H3DRes,
		particleEffect.H3DRes, maxParticleCount, respawnCount))}

	return emitter
}

func (e *Emitter) AdvanceTime(timeDelta float32) {
	horde3d.AdvanceEmitterTime(e.H3DNode, timeDelta)
}

func (e *Emitter) IsFinished() bool {
	return horde3d.HasEmitterFinished(e.H3DNode)
}

func (e *Emitter) Material() *Material {
	material := &Material{new(Resource)}
	material.H3DRes = horde3d.H3DRes(horde3d.GetNodeParamI(e.H3DNode, horde3d.Emitter_MatResI))
	return material
}

func (e *Emitter) SetMaterial(material *Material) {
	horde3d.SetNodeParamI(e.H3DNode, horde3d.Emitter_MatResI, int(material.H3DRes))
}

func (e *Emitter) ParticleEffect() *ParticleEffect {
	partEffect := &ParticleEffect{new(Resource)}
	partEffect.H3DRes = horde3d.H3DRes(horde3d.GetNodeParamI(e.H3DNode, horde3d.Emitter_PartEffResI))
	return partEffect
}

func (e *Emitter) SetParticleEffect(particleEffect *ParticleEffect) {
	horde3d.SetNodeParamI(e.H3DNode, horde3d.Emitter_PartEffResI, int(particleEffect.H3DRes))
}

func (e *Emitter) MaxCount() int {
	return horde3d.GetNodeParamI(e.H3DNode, horde3d.Emitter_MaxCountI)
}

func (e *Emitter) SetMaxCount(count int) {
	horde3d.SetNodeParamI(e.H3DNode, horde3d.Emitter_MaxCountI, count)
}

func (e *Emitter) RespawnCount() int {
	return horde3d.GetNodeParamI(e.H3DNode, horde3d.Emitter_RespawnCountI)
}

func (e *Emitter) SetRespawnCount(count int) {
	horde3d.SetNodeParamI(e.H3DNode, horde3d.Emitter_RespawnCountI, count)
}

func (e *Emitter) Delay() float32 {
	return horde3d.GetNodeParamF(e.H3DNode, horde3d.Emitter_DelayF, 0)
}

func (e *Emitter) SetDelay(delay float32) {
	horde3d.SetNodeParamF(e.H3DNode, horde3d.Emitter_DelayF, 0, delay)
}

func (e *Emitter) EmissionRate() float32 {
	return horde3d.GetNodeParamF(e.H3DNode, horde3d.Emitter_EmissionRateF, 0)
}

func (e *Emitter) SetEmissionRate(rate float32) {
	horde3d.SetNodeParamF(e.H3DNode, horde3d.Emitter_EmissionRateF, 0, rate)
}

func (e *Emitter) SpreadAngle() float32 {
	return horde3d.GetNodeParamF(e.H3DNode, horde3d.Emitter_SpreadAngleF, 0)
}

func (e *Emitter) SetSpreadAngle(angle float32) {
	horde3d.SetNodeParamF(e.H3DNode, horde3d.Emitter_SpreadAngleF, 0, angle)
}

func (e *Emitter) Force(result *vmath.Vector3) {
	result.X = horde3d.GetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 0)
	result.Y = horde3d.GetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 1)
	result.Z = horde3d.GetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 2)
}

func (e *Emitter) SetForce(force *vmath.Vector3) {
	horde3d.SetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 0, force.X)
	horde3d.SetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 1, force.Y)
	horde3d.SetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 2, force.Z)
}
