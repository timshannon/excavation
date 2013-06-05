// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package engine

import (
	"bitbucket.org/tshannon/gohorde/horde3d"
	"bitbucket.org/tshannon/vmath"
	"errors"
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

type Node struct {
	horde3d.H3DNode
	relMat      *vmath.Matrix4
	absMat      *vmath.Matrix4
	updateFrame int
}

func NewNode(hordeNode horde3d.H3DNode) *Node {
	node := &Node{
		H3DNode:     hordeNode,
		relMat:      &vmath.Matrix4{},
		absMat:      &vmath.Matrix4{},
		updateFrame: -1,
	}

	return node
}

//Adds nodes from a SceneGraph resource to the scene.
func (parent *Node) AddScene(sceneResource *Scene) (*Node, error) {
	node := NewNode(parent.H3DNode.AddNodes(sceneResource.H3DRes))

	if node.H3DNode == 0 {
		return nil, errors.New("Error adding nodes to the scene")
	}

	return node, nil
}

//Returns the parent of a scene node.
func (n *Node) Parent() *Node {
	parent := NewNode(n.H3DNode.Parent())
	return parent
}

//Relocates a node in the scene graph.
func (n *Node) SetParent(parent *Node) bool { return n.H3DNode.SetParent(parent.H3DNode) }

//Returns a slice of the children of the current node
func (n *Node) Children() []*Node {
	var hNode horde3d.H3DNode = -1
	var children []*Node
	for i := 0; hNode != 0; i++ {
		hNode = n.H3DNode.Child(i)
		if hNode != 0 {
			children = append(children, NewNode(hNode))
		}
	}
	return children
}

//This function gets the translation, rotation and scale of a specified scene node object.
// The coordinates are in local space and contain the transformation of the node relative to its parent.
func (n *Node) Transform(translate, rotate, scale *vmath.Vector3) {
	n.H3DNode.Transform(&translate[0], &translate[1], &translate[2],
		&rotate[0], &rotate[1], &rotate[2],
		&scale[0], &scale[1], &scale[2])

}

func (n *Node) Translate(result *vmath.Vector3) {
	n.H3DNode.Transform(&result[0], &result[1], &result[2],
		nil, nil, nil, nil, nil, nil)
}

func (n *Node) Rotate(result *vmath.Vector3) {
	n.H3DNode.Transform(nil, nil, nil,
		&result[0], &result[1], &result[2], nil, nil, nil)
}

func (n *Node) Scale(result *vmath.Vector3) {
	n.H3DNode.Transform(nil, nil, nil,
		nil, nil, nil, &result[0], &result[1], &result[2])
}

func (n *Node) Occluded() bool {
	//TODO: Cast physics ray from node to camera?
	return false
}

//This function sets the relative translation, rotation and scale of a
//specified scene node object.  The coordinates are in local space and
//contain the transformation of the node relative to its parent.
func (n *Node) SetTransform(translate, rotate, scale *vmath.Vector3) {
	n.H3DNode.SetTransform(translate[0], translate[1], translate[2],
		rotate[0], rotate[1], rotate[2],
		scale[0], scale[1], scale[2])
}

func (n *Node) updateTransMats() {
	//only jump into cgo code if the matrices haven't
	// been updated for this frame
	//TODO: Time CGO vs frame check is it worth it
	if n.updateFrame != frames {
		n.TransMats(n.relMat.Array(), n.absMat.Array())
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
	n.updateFrame = -1
	n.SetNodeTransMat(n.relMat.Array())
}

func (n *Node) SetLocalTransform(translate, rotate *vmath.Vector3) {
	//set transform relative to itself
	n.SetTransformRelativeTo(n, translate, rotate)
}

//SetTransformRelativeTo sets the transform relative
// to another node's position and rotation
// Note this function creates a lot of temp variables, and may
// cause GC collection performance issues
func (n *Node) SetTransformRelativeTo(otherNode *Node, trans, rotate *vmath.Vector3) {
	transform := &vmath.Transform3{}
	m3 := &vmath.Matrix3{}
	translate := &vmath.Vector3{}
	newTranslate := &vmath.Vector3{}
	rotM3 := &vmath.Matrix3{}

	//vmath.M3MakeRotationZYX(rotM3, rotate)
	rotM3.MakeRotationZYX(rotate)

	//vmath.V3Copy(newTranslate, trans)
	newTranslate.Copy(trans)

	matrix := otherNode.RelativeTransMat()
	//vmath.M4GetTranslation(translate, matrix)
	matrix.Translation(translate)
	//vmath.M4GetUpper3x3(m3, matrix)
	matrix.Upper3x3(m3)

	transform.SetUpper3x3(m3)
	//vmath.T3MulV3(newTranslate, transform, newTranslate)
	newTranslate.MulT3Self(transform)
	//vmath.M3Mul(rotM3, m3, rotM3)
	rotM3.MulSelf(m3)

	//vmath.V3Add(newTranslate, translate, newTranslate)
	newTranslate.AddToSelf(translate)

	//vmath.M4MakeFromM3V3(matrix, rotM3, newTranslate)
	matrix.MakeFromM3V3(rotM3, newTranslate)
	n.SetRelativeTransMat(matrix)

}

//Returns the bounds of a box that encompasses the node
func (n *Node) BoundingBox(min, max *vmath.Vector3) {
	n.H3DNode.AABB(&max[0], &max[1], &max[2], &max[0], &max[1], &max[2])
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
func (n *Node) NoDraw() bool { return horde3d.NodeFlags_NoDraw == n.H3DNode.Flags() }

//Excludes scene node from list of shadow casters
func (n *Node) NoCastShadow() bool {
	return horde3d.NodeFlags_NoCastShadow == n.H3DNode.Flags()
}

//Excludes scene node from ray intersection queries
func (n *Node) NoRayQuery() bool {
	return horde3d.NodeFlags_NoRayQuery == n.H3DNode.Flags()
}

//Deactivates scene node so that it is completely ignored (combination of all flags above)
func (n *Node) Inactive() bool { return horde3d.NodeFlags_Inactive == n.H3DNode.Flags() }

//Gets the name of the node
func (n *Node) Name() string {
	return n.NodeParamStr(horde3d.NodeParams_NameStr)
}

//SetName sets the name of the node
func (n *Node) SetName(name string) {
	n.SetNodeParamStr(horde3d.NodeParams_NameStr, name)
}

//Optional application-specific meta data for a node encapsulated in an Attachment XML string
func (n *Node) Attachment() string {
	return n.NodeParamStr(horde3d.NodeParams_AttachmentStr)
}

//Optional application-specific meta data for a node encapsulated in an Attachment XML string
func (n *Node) SetAttachment(value string) {
	n.SetNodeParamStr(horde3d.NodeParams_AttachmentStr, value)
}

//Returns true if both nodes refer to the same internal node
func (n *Node) IsSame(other *Node) bool {
	return n.H3DNode == other.H3DNode
}

type CastRayResult struct {
	ResultNode   *Node
	Distance     float32
	Intersection *vmath.Vector3
}

//This function checks recursively if the specified ray intersects the specified node or one of its children.
//The function finds intersections relative to the ray origin and returns the number of intersecting scene nodes.
//The ray is a line segment and is specified by a starting point (the origin) and a finite direction vector
//which also defines its length.  Currently this function is limited to returning intersections with Meshes.
//For Meshes, the base LOD (LOD0) is always used for performing the ray-triangle intersection tests.
func (n *Node) CastRay(results []*CastRayResult, origin, direction *vmath.Vector3) {
	size := n.H3DNode.CastRay(origin[0], origin[1], origin[2],
		direction[0], direction[1], direction[2], len(results))

	results = results[:size]
	for i := range results {
		results[i].ResultNode = NewNode(0)
		results[i].Intersection = &vmath.Vector3{}
		_ = horde3d.CastRayResult(i, &results[i].ResultNode.H3DNode, &results[i].Distance,
			results[i].Intersection.Array())
	}
}

//This function checks if a specified node is visible from the perspective of a specified camera.
//The function always checks if the node is in the camera.s frustum.  If checkOcclusion is true,
//the function will take into account the occlusion culling information from the previous frame
//(if occlusion culling is disabled the flag is ignored).  The flag calcLod determines whether the
//detail level for the node should be returned in case it is visible.  The function returns -1 if
//the node is not visible, otherwise 0 (base LOD level) or the computed LOD level
func (n *Node) IsVisible(camera *Camera, checkOcclusion, calcLOD bool) int {
	return n.H3DNode.CheckNodeVisibility(camera.H3DNode, checkOcclusion, calcLOD)
}

type Group struct{ *Node }

//Adds a new group node
func AddGroup(parent *Node, name string) (*Group, error) {
	group := &Group{NewNode(parent.H3DNode.AddGroupNode(name))}
	if group.H3DNode == 0 {
		return nil, errors.New("Error adding group node")
	}
	return group, nil
}

type Model struct{ *Node }

//Adds a new model
func AddModel(parent *Node, name string, geometry *Geometry) (*Model, error) {
	model := &Model{NewNode(parent.H3DNode.AddModelNode(name, geometry.H3DRes))}
	if model.H3DNode == 0 {
		return nil, errors.New("Error adding Model")
	}
	return model, nil
}

//Gets the Geometry resource for the given model
func (m *Model) Geometry() *Geometry {
	geom := &Geometry{new(Resource)}
	geom.H3DRes = horde3d.H3DRes(m.H3DNode.NodeParamI(horde3d.Model_GeoResI))
	return geom
}

//Sets the Geometry resource for the given model
func (m *Model) SetGeometry(newGeom *Geometry) {
	m.SetNodeParamI(horde3d.Model_GeoResI, int(newGeom.H3DRes))
}

//Gets state of software skinning
func (m *Model) SWSkinning() int {
	return m.NodeParamI(horde3d.Model_SWSkinningI)
}

//Sets software skinning
func (m *Model) SetSWSkinning(value int) {
	m.SetNodeParamI(horde3d.Model_SWSkinningI, value)
}

//Gets the distances for the LevelOfDetail settings
func (m *Model) LODDist() (LOD1, LOD2, LOD3, LOD4 float32) {
	LOD1 = m.NodeParamF(horde3d.Model_LodDist1F, 0)
	LOD2 = m.NodeParamF(horde3d.Model_LodDist2F, 0)
	LOD3 = m.NodeParamF(horde3d.Model_LodDist3F, 0)
	LOD4 = m.NodeParamF(horde3d.Model_LodDist4F, 0)
	return
}

//Sets the distances for the LevelOfDetail settings
// subsequent LODs must be greater than the previous i.e. LOD1 < LOD2
func (m *Model) SetLODDist(LOD1, LOD2, LOD3, LOD4 float32) {
	m.H3DNode.SetNodeParamF(horde3d.Model_LodDist1F, 0, LOD1)
	m.H3DNode.SetNodeParamF(horde3d.Model_LodDist2F, 0, LOD2)
	m.H3DNode.SetNodeParamF(horde3d.Model_LodDist3F, 0, LOD3)
	m.H3DNode.SetNodeParamF(horde3d.Model_LodDist4F, 0, LOD4)
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
	mesh := &Mesh{NewNode(parent.H3DNode.AddMeshNode(name, material.H3DRes, batchStart,
		batchCount, vertRStart, vertREnd))}
	if mesh.H3DNode == 0 {
		return nil, errors.New("Error adding Mesh")
	}
	return mesh, nil
}

func (m *Mesh) Material() *Material {
	material := &Material{new(Resource)}
	material.H3DRes = horde3d.H3DRes(m.H3DNode.NodeParamI(horde3d.Mesh_MatResI))
	return material
}

func (m *Mesh) SetMaterial(newMaterial *Material) {
	m.H3DNode.SetNodeParamI(horde3d.Mesh_MatResI, int(newMaterial.H3DRes))
}

func (m *Mesh) BatchStart() int { return m.H3DNode.NodeParamI(horde3d.Mesh_BatchStartI) }
func (m *Mesh) BatchCount() int { return m.H3DNode.NodeParamI(horde3d.Mesh_BatchCountI) }
func (m *Mesh) VertRStart() int { return m.H3DNode.NodeParamI(horde3d.Mesh_VertRStartI) }
func (m *Mesh) VertREnd() int   { return m.H3DNode.NodeParamI(horde3d.Mesh_VertREndI) }

func (m *Mesh) LODLevel() int { return m.H3DNode.NodeParamI(horde3d.Mesh_LodLevelI) }
func (m *Mesh) SetLODLevel(level int) {
	m.H3DNode.SetNodeParamI(horde3d.Mesh_LodLevelI, level)
}

type Joint struct{ *Node }

func AddJoint(parent *Node, name string, jointIndex int) (*Joint, error) {
	joint := &Joint{NewNode(parent.H3DNode.AddJointNode(name, jointIndex))}
	if joint.H3DNode == 0 {
		return nil, errors.New("Error adding Joint")
	}
	return joint, nil
}

func (j *Joint) Index() int { return j.H3DNode.NodeParamI(horde3d.Joint_JointIndexI) }

type Light struct{ *Node }

func AddLight(parent *Node, name string, material *Material, lightingContext string,
	shadowContext string) *Light {
	light := &Light{NewNode(parent.H3DNode.AddLightNode(name, material.H3DRes,
		lightingContext, shadowContext))}
	return light
}

func (l *Light) Material() *Material {
	material := &Material{new(Resource)}
	material.H3DRes = horde3d.H3DRes(l.H3DNode.NodeParamI(horde3d.Light_MatResI))
	return material
}

func (l *Light) SetMaterial(material *Material) {
	l.H3DNode.SetNodeParamI(horde3d.Light_MatResI, int(material.H3DRes))
}

func (l *Light) FOV() float32 { return l.H3DNode.NodeParamF(horde3d.Light_FovF, 0) }
func (l *Light) SetFOV(newFOV float32) {
	l.H3DNode.SetNodeParamF(horde3d.Light_FovF, 0, newFOV)
}

func (l *Light) Color() (r, g, b float32) {
	r = l.H3DNode.NodeParamF(horde3d.Light_ColorF3, 0)
	b = l.H3DNode.NodeParamF(horde3d.Light_ColorF3, 1)
	g = l.H3DNode.NodeParamF(horde3d.Light_ColorF3, 2)
	return
}

func (l *Light) SetColor(r, g, b float32) {
	l.H3DNode.SetNodeParamF(horde3d.Light_ColorF3, 0, r)
	l.H3DNode.SetNodeParamF(horde3d.Light_ColorF3, 1, g)
	l.H3DNode.SetNodeParamF(horde3d.Light_ColorF3, 2, b)
}

func (l *Light) ColorMultiplier() float32 {
	return l.H3DNode.NodeParamF(horde3d.Light_ColorMultiplierF, 0)
}

func (l *Light) SetColorMultiplier(multiplier float32) {
	l.H3DNode.SetNodeParamF(horde3d.Light_ColorMultiplierF, 0, multiplier)
}

func (l *Light) ShadowMapCount() int {
	return l.H3DNode.NodeParamI(horde3d.Light_ShadowMapCountI)
}

func (l *Light) SetShadowMapCount(count int) {
	l.H3DNode.SetNodeParamI(horde3d.Light_ShadowMapCountI, count)
}

func (l *Light) ShadowSplitLambda() float32 {
	return l.H3DNode.NodeParamF(horde3d.Light_ShadowSplitLambdaF, 0)
}

func (l *Light) SetShadowSplitLambda(lambda float32) {
	l.H3DNode.SetNodeParamF(horde3d.Light_ShadowSplitLambdaF, 0, lambda)
}

func (l *Light) ShadowMapBias() float32 {
	return l.H3DNode.NodeParamF(horde3d.Light_ShadowMapBiasF, 0)
}

func (l *Light) SetShadowMapBias(bias float32) {
	l.H3DNode.SetNodeParamF(horde3d.Light_ShadowMapBiasF, 0, bias)
}

func (l *Light) LightingContext() string {
	return l.H3DNode.NodeParamStr(horde3d.Light_LightingContextStr)
}

func (l *Light) SetLightingContext(context string) {
	l.H3DNode.SetNodeParamStr(horde3d.Light_LightingContextStr, context)
}

func (l *Light) ShadowContext() string {
	return l.H3DNode.NodeParamStr(horde3d.Light_ShadowContextStr)
}

func (l *Light) SetShadowContext(context string) {
	l.H3DNode.SetNodeParamStr(horde3d.Light_ShadowContextStr, context)
}

type Camera struct{ *Node }

func AddCamera(parent *Node, name string, pipeline *Pipeline) *Camera {
	camera := &Camera{NewNode(parent.H3DNode.AddCameraNode(name, pipeline.H3DRes))}
	return camera
}

func (c *Camera) SetupView(FOV, aspect, nearDist, farDist float32) {
	horde3d.SetupCameraView(c.H3DNode, FOV, aspect, nearDist, farDist)
}

func (c *Camera) ProjectionMatrix(result *vmath.Matrix4) {
	horde3d.GetCameraProjMat(c.H3DNode, result.Array())
}

func (c *Camera) Pipeline() *Pipeline {
	pipeline := &Pipeline{new(Resource)}
	pipeline.H3DRes = horde3d.H3DRes(c.H3DNode.NodeParamI(horde3d.Camera_PipeResI))
	return pipeline
}

func (c *Camera) SetPipeline(pipeline *Pipeline) {
	c.H3DNode.SetNodeParamI(horde3d.Camera_PipeResI, int(pipeline.H3DRes))
}

//2D Texture resource used as output buffer (can be 0 to use main framebuffer) (default: 0)
func (c *Camera) OutTexture() *Texture {
	texture := &Texture{new(Resource)}
	texture.H3DRes = horde3d.H3DRes(c.H3DNode.NodeParamI(horde3d.Camera_OutTexResI))
	return texture
}

func (c *Camera) SetOutTexture(texture *Texture) {
	c.H3DNode.SetNodeParamI(horde3d.Camera_OutTexResI, int(texture.H3DRes))
}

//Index of the output buffer for stereo rendering (values: 0 for left eye, 1 for right eye) (default: 0)
func (c *Camera) OutputBufferIndex() int {
	return c.H3DNode.NodeParamI(horde3d.Camera_OutBufIndexI)
}

func (c *Camera) SetOutputBufferIndex(index int) {
	c.H3DNode.SetNodeParamI(horde3d.Camera_OutBufIndexI, index)
}

func (c *Camera) Viewport() (x, y, width, height int) {
	x = c.H3DNode.NodeParamI(horde3d.Camera_ViewportXI)
	y = c.H3DNode.NodeParamI(horde3d.Camera_ViewportYI)
	width = c.H3DNode.NodeParamI(horde3d.Camera_ViewportWidthI)
	height = c.H3DNode.NodeParamI(horde3d.Camera_ViewportHeightI)
	return
}

func (c *Camera) SetViewport(x, y, width, height int) {
	c.H3DNode.SetNodeParamI(horde3d.Camera_ViewportXI, x)
	c.H3DNode.SetNodeParamI(horde3d.Camera_ViewportYI, y)
	c.H3DNode.SetNodeParamI(horde3d.Camera_ViewportWidthI, width)
	c.H3DNode.SetNodeParamI(horde3d.Camera_ViewportHeightI, height)
}

func (c *Camera) IsOrtho() bool {
	i := c.H3DNode.NodeParamI(horde3d.Camera_OrthoI)
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
	c.H3DNode.SetNodeParamI(horde3d.Camera_OrthoI, i)
}

func (c *Camera) OcclusionCulling() bool {
	i := c.H3DNode.NodeParamI(horde3d.Camera_OccCullingI)

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

	c.H3DNode.SetNodeParamI(horde3d.Camera_OccCullingI, i)
}

type Emitter struct{ *Node }

func AddEmitter(parent *Node, name string, material *Material, particleEffect *ParticleEffect,
	maxParticleCount int, respawnCount int) *Emitter {
	emitter := &Emitter{NewNode(parent.H3DNode.AddEmitterNode(name, material.H3DRes,
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
	material.H3DRes = horde3d.H3DRes(e.H3DNode.NodeParamI(horde3d.Emitter_MatResI))
	return material
}

func (e *Emitter) SetMaterial(material *Material) {
	e.H3DNode.SetNodeParamI(horde3d.Emitter_MatResI, int(material.H3DRes))
}

func (e *Emitter) ParticleEffect() *ParticleEffect {
	partEffect := &ParticleEffect{new(Resource)}
	partEffect.H3DRes = horde3d.H3DRes(e.H3DNode.NodeParamI(horde3d.Emitter_PartEffResI))
	return partEffect
}

func (e *Emitter) SetParticleEffect(particleEffect *ParticleEffect) {
	e.H3DNode.SetNodeParamI(horde3d.Emitter_PartEffResI, int(particleEffect.H3DRes))
}

func (e *Emitter) MaxCount() int {
	return e.H3DNode.NodeParamI(horde3d.Emitter_MaxCountI)
}

func (e *Emitter) SetMaxCount(count int) {
	e.H3DNode.SetNodeParamI(horde3d.Emitter_MaxCountI, count)
}

func (e *Emitter) RespawnCount() int {
	return e.H3DNode.NodeParamI(horde3d.Emitter_RespawnCountI)
}

func (e *Emitter) SetRespawnCount(count int) {
	e.H3DNode.SetNodeParamI(horde3d.Emitter_RespawnCountI, count)
}

func (e *Emitter) Delay() float32 {
	return e.H3DNode.NodeParamF(horde3d.Emitter_DelayF, 0)
}

func (e *Emitter) SetDelay(delay float32) {
	e.H3DNode.SetNodeParamF(horde3d.Emitter_DelayF, 0, delay)
}

func (e *Emitter) EmissionRate() float32 {
	return e.H3DNode.NodeParamF(horde3d.Emitter_EmissionRateF, 0)
}

func (e *Emitter) SetEmissionRate(rate float32) {
	e.H3DNode.SetNodeParamF(horde3d.Emitter_EmissionRateF, 0, rate)
}

func (e *Emitter) SpreadAngle() float32 {
	return e.H3DNode.NodeParamF(horde3d.Emitter_SpreadAngleF, 0)
}

func (e *Emitter) SetSpreadAngle(angle float32) {
	e.H3DNode.SetNodeParamF(horde3d.Emitter_SpreadAngleF, 0, angle)
}

func (e *Emitter) Force(result *vmath.Vector3) {
	result[0] = e.H3DNode.NodeParamF(horde3d.Emitter_ForceF3, 0)
	result[1] = e.H3DNode.NodeParamF(horde3d.Emitter_ForceF3, 1)
	result[2] = e.H3DNode.NodeParamF(horde3d.Emitter_ForceF3, 2)
}

func (e *Emitter) SetForce(force *vmath.Vector3) {
	e.H3DNode.SetNodeParamF(horde3d.Emitter_ForceF3, 0, force[0])
	e.H3DNode.SetNodeParamF(horde3d.Emitter_ForceF3, 1, force[1])
	e.H3DNode.SetNodeParamF(horde3d.Emitter_ForceF3, 2, force[2])
}
