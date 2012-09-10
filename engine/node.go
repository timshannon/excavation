// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"errors"
	"excavation/math3d"
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
	//non-horde types
	NodeTypeSound = horde3d.NodeTypes_Emitter + 1 + iota
	NodeTypeRigidBody
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
func AddNodes(parent *Node, sceneResource *Scene) (*Node, error) {
	node := new(Node)
	node.H3DNode = horde3d.AddNodes(parent.H3DNode, sceneResource.H3DRes)

	if node.H3DNode == 0 {
		return nil, errors.New("Error adding nodes to the scene")
	}

	return node, nil
}

//This function returns the type of a specified scene node.  If the node handle is invalid, 
//the function returns the node type Unknown.
func (n *Node) Type() int {
	intType := horde3d.GetNodeType(n.H3DNode)
	if intType == NodeTypeGroup {
		//TODO: Look up type from attachment
	}
	return intType
}

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

func (j *Joint) Index() int { return horde3d.GetNodeParamI(j.H3DNode, horde3d.Joint_JointIndexI) }

type Light struct{ *Node }

func AddLight(parent *Node, name string, material *Material, lightingContext string,
	shadowContext string) *Light {
	light := new(Light)
	light.H3DNode = horde3d.AddLightNode(parent.H3DNode, name, material.H3DRes,
		lightingContext, shadowContext)
	return light
}

func (l *Light) Material() *Material {
	material := new(Material)
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

func (l *Light) Color() math3d.Vector3 {
	r := horde3d.GetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 0)
	b := horde3d.GetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 1)
	g := horde3d.GetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 2)
	return math3d.MakeVector3(r, g, b)
}

func (l *Light) SetColor(color math3d.Vector3) {
	horde3d.SetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 0, color[0])
	horde3d.SetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 1, color[1])
	horde3d.SetNodeParamF(l.H3DNode, horde3d.Light_ColorF3, 2, color[2])
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
	camera := new(Camera)
	camera.H3DNode = horde3d.AddCameraNode(parent.H3DNode, name, pipeline.H3DRes)
	return camera
}

func (c *Camera) SetupView(FOV, aspect, nearDist, farDist float32) {
	horde3d.SetupCameraView(c.H3DNode, FOV, aspect, nearDist, farDist)
}

func (c *Camera) ProjectionMatrix() math3d.Matrix4 {
	matrix := math3d.MakeMatrix4()
	horde3d.GetCameraProjMat(c.H3DNode, matrix)
	return matrix
}

func (c *Camera) Pipeline() *Pipeline {
	pipeline := new(Pipeline)
	pipeline.H3DRes = horde3d.H3DRes(horde3d.GetNodeParamI(c.H3DNode, horde3d.Camera_PipeResI))
	return pipeline
}

func (c *Camera) SetPipeline(pipeline *Pipeline) {
	horde3d.SetNodeParamI(c.H3DNode, horde3d.Camera_PipeResI, int(pipeline.H3DRes))
}

//2D Texture resource used as output buffer (can be 0 to use main framebuffer) (default: 0)
func (c *Camera) OutTexture() *Texture {
	texture := new(Texture)
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

func (c *Camera) ViewPlanes() (left, right, bottom, top float32) {
	left = horde3d.GetNodeParamF(c.H3DNode, horde3d.Camera_LeftPlaneF, 0)
	right = horde3d.GetNodeParamF(c.H3DNode, horde3d.Camera_RightPlaneF, 0)
	bottom = horde3d.GetNodeParamF(c.H3DNode, horde3d.Camera_BottomPlaneF, 0)
	top = horde3d.GetNodeParamF(c.H3DNode, horde3d.Camera_TopPlaneF, 0)
	return
}

func (c *Camera) SetViewPlanes(left, right, bottom, top float32) {
	horde3d.SetNodeParamF(c.H3DNode, horde3d.Camera_LeftPlaneF, 0, left)
	horde3d.SetNodeParamF(c.H3DNode, horde3d.Camera_RightPlaneF, 0, right)
	horde3d.SetNodeParamF(c.H3DNode, horde3d.Camera_BottomPlaneF, 0, bottom)
	horde3d.SetNodeParamF(c.H3DNode, horde3d.Camera_TopPlaneF, 0, top)
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
	emitter := new(Emitter)
	emitter.H3DNode = horde3d.AddEmitterNode(parent.H3DNode, name, material.H3DRes,
		particleEffect.H3DRes, maxParticleCount, respawnCount)

	return emitter
}

func (e *Emitter) AdvanceTime(timeDelta float32) {
	horde3d.AdvanceEmitterTime(e.H3DNode, timeDelta)
}

func (e *Emitter) IsFinished() bool {
	return horde3d.HasEmitterFinished(e.H3DNode)
}

func (e *Emitter) Material() *Material {
	material := new(Material)
	material.H3DRes = horde3d.H3DRes(horde3d.GetNodeParamI(e.H3DNode, horde3d.Emitter_MatResI))
	return material
}

func (e *Emitter) SetMaterial(material *Material) {
	horde3d.SetNodeParamI(e.H3DNode, horde3d.Emitter_MatResI, int(material.H3DRes))
}

func (e *Emitter) ParticleEffect() *ParticleEffect {
	partEffect := new(ParticleEffect)
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

func (e *Emitter) Force() math3d.Vector3 {
	x := horde3d.GetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 0)
	y := horde3d.GetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 1)
	z := horde3d.GetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 2)

	return math3d.MakeVector3(x, y, z)
}

func (e *Emitter) SetForce(force math3d.Vector3) {
	horde3d.SetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 0, force[0])
	horde3d.SetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 1, force[1])
	horde3d.SetNodeParamF(e.H3DNode, horde3d.Emitter_ForceF3, 2, force[2])
}
