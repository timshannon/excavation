package horde3d

/*
#cgo LDFLAGS: -lHorde3D
#include "goHorde3D.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "reflect"

//typedef int H3DRes;
//typedef int H3DNode;
type H3DRes int
type H3DNode int

//const H3DNode H3DRootNode = 1;
const H3DRootNode H3DNode = 1

/* Group: Enumerations */

/* Enum: H3DOptions
	The available engine option parameters.

MaxLogLevel         - Defines the maximum log level; only messages which are smaller or equal to this value
                      (hence more important) are published in the message queue. (Default: 4)
MaxNumMessages      - Defines the maximum number of messages that can be stored in the message queue (Default: 512)
TrilinearFiltering  - Enables or disables trilinear filtering for textures. (Values: 0, 1; Default: 1)
MaxAnisotropy       - Sets the maximum quality for anisotropic filtering. (Values: 1, 2, 4, 8; Default: 1)
TexCompression      - Enables or disables texture compression; only affects textures that are
                      loaded after setting the option. (Values: 0, 1; Default: 0)
SRGBLinearization   - Eanbles or disables gamma-to-linear-space conversion of input textures that are tagged as sRGB (Values: 0, 1; Default: 0)
LoadTextures        - Enables or disables loading of textures referenced by materials; this can be useful to reduce
                      loading times for testing. (Values: 0, 1; Default: 1)
FastAnimation       - Disables or enables inter-frame interpolation for animations. (Values: 0, 1; Default: 1)
ShadowMapSize       - Sets the size of the shadow map buffer (Values: 128, 256, 512, 1024, 2048; Default: 1024)
SampleCount         - Maximum number of samples used for multisampled render targets; only affects pipelines
                      that are loaded after setting the option. (Values: 0, 2, 4, 8, 16; Default: 0)
WireframeMode       - Enables or disables wireframe rendering
DebugViewMode       - Enables or disables debug view where geometry is rendered in wireframe without shaders and
                      lights are visualized using their screen space bounding box. (Values: 0, 1; Default: 0)
DumpFailedShaders   - Enables or disables storing of shader code that failed to compile in a text file; this can be
                      useful in combination with the line numbers given back by the shader compiler. (Values: 0, 1; Default: 0)
GatherTimeStats     - Enables or disables gathering of time stats that are useful for profiling (Values: 0, 1; Default: 1)
*/
const (
	_ = iota
	H3DOptions_MaxLogLevel
	H3DOptions_MaxNumMessages
	H3DOptions_TrilinearFiltering
	H3DOptions_MaxAnisotropy
	H3DOptions_TexCompression
	H3DOptions_SRGBLinearization
	H3DOptions_LoadTextures
	H3DOptions_FastAnimation
	H3DOptions_ShadowMapSize
	H3DOptions_SampleCount
	H3DOptions_WireframeMode
	H3DOptions_DebugViewMode
	H3DOptions_DumpFailedShaders
	H3DOptions_GatherTimeStats
)

/* Enum: H3DStats
	The available engine statistic parameters.

TriCount          - Number of triangles that were pushed to the renderer
BatchCount        - Number of batches (draw calls)
LightPassCount    - Number of lighting passes
FrameTime         - Time in ms between two h3dFinalizeFrame calls
AnimationTime     - CPU time in ms spent for animation
GeoUpdateTime     - CPU time in ms spent for software skinning and morphing
ParticleSimTime   - CPU time in ms spent for particle simulation and updates
FwdLightsGPUTime  - GPU time in ms spent for forward lighting passes
DefLightsGPUTime  - GPU time in ms spent for drawing deferred light volumes
ShadowsGPUTime    - GPU time in ms spent for generating shadow maps
ParticleGPUTime   - GPU time in ms spent for drawing particles
TextureVMem       - Estimated amount of video memory used by textures (in Mb)
GeometryVMem      - Estimated amount of video memory used by geometry (in Mb)
*/
const (
	_ = iota + 100
	H3DStats_TriCount
	H3DStats_BatchCount
	H3DStats_LightPassCount
	H3DStats_FrameTime
	H3DStats_AnimationTime
	H3DStats_GeoUpdateTime
	H3DStats_ParticleSimTime
	H3DStats_FwdLightsGPUTime
	H3DStats_DefLightsGPUTime
	H3DStats_ShadowsGPUTime
	H3DStats_ParticleGPUTime
	H3DStats_TextureVMem
	H3DStats_GeometryVMem
)

/* Enum: H3DResTypes
	The available resource types.

Undefined       - An undefined resource, returned by getResourceType in case of error
SceneGraph      - Scene graph subtree stored in XML format
Geometry        - Geometrical data containing bones, vertices and triangles
Animation       - Animation data
Material        - Material script
Code            - Text block containing shader source code
Shader          - Shader program
Texture         - Texture map
ParticleEffect  - Particle configuration
Pipeline        - Rendering pipeline
*/
const (
	H3DResTypes_Undefined = iota
	H3DResTypes_SceneGraph
	H3DResTypes_Geometry
	H3DResTypes_Animation
	H3DResTypes_Material
	H3DResTypes_Code
	H3DResTypes_Shader
	H3DResTypes_Texture
	H3DResTypes_ParticleEffect
	H3DResTypes_Pipeline
)

/* Enum: H3DResFlags
	The available flags used when adding a resource.

NoQuery           - Excludes resource from being listed by queryUnloadedResource function.
NoTexCompression  - Disables texture compression for Texture resource.
NoTexMipmaps      - Disables generation of mipmaps for Texture resource.
TexCubemap        - Sets Texture resource to be a cubemap.
TexDynamic        - Enables more efficient updates of Texture resource streams.
TexRenderable     - Makes Texture resource usable as render target.
TexSRGB           - Indicates that Texture resource is in sRGB color space and should be converted
                    to linear space when being sampled.
*/
const (
	H3DResFlags_NoQuery          = 1
	H3DResFlags_NoTexCompression = 2
	H3DResFlags_NoTexMipmaps     = 4
	H3DResFlags_TexCubemap       = 8
	H3DResFlags_TexDynamic       = 16
	H3DResFlags_TexRenderable    = 32
	H3DResFlags_TexSRGB          = 64
)

/* Enum: H3DFormats
	The available resource stream formats.

Unknown      - Unknown format
TEX_BGRA8    - 8-bit BGRA texture
TEX_DXT1     - DXT1 compressed texture
TEX_DXT3     - DXT3 compressed texture
TEX_DXT5     - DXT5 compressed texture
TEX_RGBA16F  - Half float RGBA texture
TEX_RGBA32F  - Float RGBA texture
*/
const (
	H3DFormats_Unknown = iota
	H3DFormats_TEX_BGRA8
	H3DFormats_TEX_DXT1
	H3DFormats_TEX_DXT3
	H3DFormats_TEX_DXT5
	H3DFormats_TEX_RGBA16F
	H3DFormats_TEX_RGBA32F
)

/* Enum: H3DGeoRes
	The available Geometry resource accessors.

GeometryElem         - Base element
GeoIndexCountI       - Number of indices [read-only]
GeoVertexCountI      - Number of vertices [read-only]
GeoIndices16I        - Flag indicating whether index data is 16 or 32 bit [read-only]
GeoIndexStream       - Triangle index data (uint16 or uint32, depending on flag)
GeoVertPosStream     - Vertex position data (float x, y, z)
GeoVertTanStream     - Vertex tangent frame data (float nx, ny, nz, tx, ty, tz, tw)
GeoVertStaticStream  - Vertex static attribute data (float u0, v0,
                         float4 jointIndices, float4 jointWeights, float u1, v1)
*/
const (
	_ = iota + 200
	H3DGeoRes_GeometryElem
	H3DGeoRes_GeoIndexCountI
	H3DGeoRes_GeoVertexCountI
	H3DGeoRes_GeoIndices16I
	H3DGeoRes_GeoIndexStream
	H3DGeoRes_GeoVertPosStream
	H3DGeoRes_GeoVertTanStream
	H3DGeoRes_GeoVertStaticStream
)

/* Enum: H3DAnimRes
	The available Animation resource accessors.	  

EntityElem      - Stored animation entities (joints and meshes)
EntFrameCountI  - Number of frames stored for a specific entity [read-only]
*/
const (
	_ = iota + 300
	H3DAnimRes_EntityElem
	H3DAnimRes_EntFrameCountI
)

/* Enum: H3DMatRes
	The available Material resource accessors.

MaterialElem  - Base element
SamplerElem   - Sampler element
UniformElem   - Uniform element
MatClassStr   - Material class
MatLinkI      - Material resource that is linked to this material
MatShaderI    - Shader resource
SampNameStr   - Name of sampler [read-only]
SampTexResI   - Texture resource bound to sampler
UnifNameStr   - Name of uniform [read-only]
UnifValueF4   - Value of uniform (a, b, c, d)
*/
const (
	_ = iota + 400
	H3DMatRes_MaterialElem
	H3DMatRes_SamplerElem
	H3DMatRes_UniformElem
	H3DMatRes_MatClassStr
	H3DMatRes_MatLinkI
	H3DMatRes_MatShaderI
	H3DMatRes_SampNameStr
	H3DMatRes_SampTexResI
	H3DMatRes_UnifNameStr
	H3DMatRes_UnifValueF4
)

/* Enum: H3DShaderRes
	The available Shader resource accessors.

ContextElem     - Context element 
SamplerElem     - Sampler element
UniformElem     - Uniform element
ContNameStr     - Name of context [read-only]
SampNameStr     - Name of sampler [read-only]
UnifNameStr     - Name of uniform [read-only]
UnifSizeI       - Size (number of components) of uniform [read-only]
UnifDefValueF4  - Default value of uniform (a, b, c, d)
*/
const (
	_ = iota + 600
	H3DShaderRes_ContextElem
	H3DShaderRes_SamplerElem
	H3DShaderRes_UniformElem
	H3DShaderRes_ContNameStr
	H3DShaderRes_SampNameStr
	H3DShaderRes_UnifNameStr
	H3DShaderRes_UnifSizeI
	H3DShaderRes_UnifDefValueF4
)

/* Enum: H3DTexRes
	The available Texture resource accessors.

TextureElem     - Base element
ImageElem       - Subresources of the texture. A texture consists, depending on the type,
                  of a number of equally sized slices which again can have a fixed number
                  of mipmaps. Each image element represents the base image of a slice or
                  a single mipmap level of the corresponding slice.
TexFormatI      - Texture format [read-only]
TexSliceCountI  - Number of slices (1 for 2D texture and 6 for cubemap) [read-only]
ImgWidthI       - Image width [read-only]
ImgHeightI      - Image height [read-only]
ImgPixelStream  - Pixel data of an image. The data layout matches the layout specified
                  by the texture format with the exception that half-float is converted
                  to float. The first element in the data array corresponds to the lower
                  left corner.
*/
const (
	_ = iota + 700
	H3DTexRes_TextureElem
	H3DTexRes_ImageElem
	H3DTexRes_TexFormatI
	H3DTexRes_TexSliceCountI
	H3DTexRes_ImgWidthI
	H3DTexRes_ImgHeightI
	H3DTexRes_ImgPixelStream
)

/* Enum: H3DPartEffRes
	The available ParticleEffect resource accessors.

ParticleElem     - General particle configuration
ChanMoveVelElem  - Velocity channel
ChanRotVelElem   - Angular velocity channel
ChanSizeElem     - Size channel
ChanColRElem     - Red color component channel
ChanColGElem     - Green color component channel
ChanColBElem     - Blue color component channel
ChanColAElem     - Alpha channel
PartLifeMinF     - Minimum value of random life time (in seconds)
PartLifeMaxF     - Maximum value of random life time (in seconds)
ChanStartMinF    - Minimum for selecting initial random value of channel
ChanStartMaxF    - Maximum for selecting initial random value of channel
ChanEndRateF     - Remaining percentage of initial value when particle is dying
*/
const (
	_ = iota + 800
	H3DPartEffRes_ParticleElem
	H3DPartEffRes_ChanMoveVelElem
	H3DPartEffRes_ChanRotVelElem
	H3DPartEffRes_ChanSizeElem
	H3DPartEffRes_ChanColRElem
	H3DPartEffRes_ChanColGElem
	H3DPartEffRes_ChanColBElem
	H3DPartEffRes_ChanColAElem
	H3DPartEffRes_PartLifeMinF
	H3DPartEffRes_PartLifeMaxF
	H3DPartEffRes_ChanStartMinF
	H3DPartEffRes_ChanStartMaxF
	H3DPartEffRes_ChanEndRateF
	H3DPartEffRes_ChanDragElem
)

/* Enum: H3DPipeRes
	The available Pipeline resource accessors.

StageElem         - Pipeline stage
StageNameStr      - Name of stage [read-only]
StageActivationI  - Flag indicating whether stage is active
*/
const (
	_ = iota + 900
	H3DPipeRes_StageElem
	H3DPipeRes_StageNameStr
	H3DPipeRes_StageActivationI
)

/*	Enum: H3DNodeTypes
		The available scene node types.

	Undefined  - An undefined node type, returned by getNodeType in case of error
	Group      - Group of different scene nodes
	Model      - 3D model with optional skeleton
	Mesh       - Subgroup of a model with triangles of one material
	Joint      - Joint for skeletal animation
	Light      - Light source
	Camera     - Camera giving view on scene
	Emitter    - Particle system emitter
*/
const (
	H3DNodeTypes_Undefined = iota
	H3DNodeTypes_Group
	H3DNodeTypes_Model
	H3DNodeTypes_Mesh
	H3DNodeTypes_Joint
	H3DNodeTypes_Light
	H3DNodeTypes_Camera
	H3DNodeTypes_Emitter
)

/*	Enum: H3DNodeFlags
		The available scene node flags.

	NoDraw         - Excludes scene node from all rendering
	NoCastShadow   - Excludes scene node from list of shadow casters
	NoRayQuery     - Excludes scene node from ray intersection queries
	Inactive       - Deactivates scene node so that it is completely ignored
	                 (combination of all flags above)
*/
const (
	H3DNodeFlags_NoDraw       = 1
	H3DNodeFlags_NoCastShadow = 2
	H3DNodeFlags_NoRayQuery   = 4
	H3DNodeFlags_Inactive     = 7 // NoDraw | NoCastShadow | NoRayQuery
)

/*	Enum: H3DNodeParams
		The available scene node parameters.

	NameStr        - Name of the scene node
	AttachmentStr  - Optional application-specific meta data for a node encapsulated
	                 in an 'Attachment' XML string
*/
const (
	H3DNodeParams_NameStr = 1
	H3DNodeParams_AttachmentStr
)

/*	Enum: H3DModel
		The available Model node parameters

	GeoResI      - Geometry resource used for the model
	SWSkinningI  - Enables or disables software skinning (default: 0)
	LodDist1F    - Distance to camera from which on LOD1 is used (default: infinite)
	               (must be a positive value larger than 0.0)
	LodDist2F    - Distance to camera from which on LOD2 is used
	               (may not be smaller than LodDist1) (default: infinite)
	LodDist3F    - Distance to camera from which on LOD3 is used
	               (may not be smaller than LodDist2) (default: infinite)
	LodDist4F    - Distance to camera from which on LOD4 is used
	               (may not be smaller than LodDist3) (default: infinite)
*/
const (
	_                         = iota + 200
	H3DModel_H3DModel_GeoResI = 200
	H3DModel_SWSkinningI
	H3DModel_LodDist1F
	H3DModel_LodDist2F
	H3DModel_LodDist3F
	H3DModel_LodDist4F
)

/*	Enum: H3DMesh
		The available Mesh node parameters.

	MatResI      - Material resource used for the mesh
	BatchStartI  - First triangle index of mesh in Geometry resource of parent Model node [read-only]
	BatchCountI  - Number of triangle indices used for drawing mesh [read-only]
	VertRStartI  - First vertex in Geometry resource of parent Model node [read-only]
	VertREndI    - Last vertex in Geometry resource of parent Model node [read-only]
	LodLevelI    - LOD level of Mesh; the mesh is only rendered if its LOD level corresponds to
	               the model's current LOD level which is calculated based on the LOD distances (default: 0)
*/
const (
	_ = iota + 300
	H3DMesh_MatResI
	H3DMesh_BatchStartI
	H3DMesh_BatchCountI
	H3DMesh_VertRStartI
	H3DMesh_VertREndI
	H3DMesh_LodLevelI
)

/*	Enum: H3DJoint
		The available Joint node parameters.

	JointIndexI  - Index of joint in Geometry resource of parent Model node [read-only]
*/
const (
	H3DJoint_JointIndexI = 400
)

/*	Enum: H3DLight
		The available Light node parameters.

	MatResI             - Material resource used for the light
	RadiusF             - Radius of influence (default: 100.0)
	FovF                - Field of view (FOV) angle (default: 90.0)
	ColorF3             - Diffuse color RGB (default: 1.0, 1.0, 1.0)
	ColorMultiplierF    - Diffuse color multiplier for altering intensity, mainly useful for HDR (default: 1.0)
	ShadowMapCountI     - Number of shadow maps used for light source (values: 0, 1, 2, 3, 4; default: 0)]
	ShadowSplitLambdaF  - Constant determining segmentation of view frustum for Parallel Split Shadow Maps (default: 0.5)
	ShadowMapBiasF      - Bias value for shadow mapping to reduce shadow acne (default: 0.005)
	LightingContextStr  - Name of shader context used for computing lighting
	ShadowContextStr    - Name of shader context used for generating shadow map
*/
const (
	_ = iota + 500
	H3DLight_MatResI
	H3DLight_RadiusF
	H3DLight_FovF
	H3DLight_ColorF3
	H3DLight_ColorMultiplierF
	H3DLight_ShadowMapCountI
	H3DLight_ShadowSplitLambdaF
	H3DLight_ShadowMapBiasF
	H3DLight_LightingContextStr
	H3DLight_ShadowContextStr
)

/*	Enum: H3DCamera
		The available Camera node parameters.

	PipeResI         - Pipeline resource used for rendering
	OutTexResI       - 2D Texture resource used as output buffer (can be 0 to use main framebuffer) (default: 0)
	OutBufIndexI     - Index of the output buffer for stereo rendering (values: 0 for left eye, 1 for right eye) (default: 0)
	LeftPlaneF       - Coordinate of left plane relative to near plane center (default: -0.055228457)
	RightPlaneF      - Coordinate of right plane relative to near plane center (default: 0.055228457)
	BottomPlaneF     - Coordinate of bottom plane relative to near plane center (default: -0.041421354f)
	TopPlaneF        - Coordinate of top plane relative to near plane center (default: 0.041421354f)
	NearPlaneF       - Distance of near clipping plane (default: 0.1)
	FarPlaneF        - Distance of far clipping plane (default: 1000)
	ViewportXI       - Position x-coordinate of the lower left corner of the viewport rectangle (default: 0)
	ViewportYI       - Position y-coordinate of the lower left corner of the viewport rectangle (default: 0)
	ViewportWidthI   - Width of the viewport rectangle (default: 320)
	ViewportHeightI  - Height of the viewport rectangle (default: 240)
	OrthoI           - Flag for setting up an orthographic frustum instead of a perspective one (default: 0)
	OccCullingI      - Flag for enabling occlusion culling (default: 0)
*/

const (
	_ = iota + 600
	H3DCamera_PipeResI
	H3DCamera_OutTexResI
	H3DCamera_OutBufIndexI
	H3DCamera_LeftPlaneF
	H3DCamera_RightPlaneF
	H3DCamera_BottomPlaneF
	H3DCamera_TopPlaneF
	H3DCamera_NearPlaneF
	H3DCamera_FarPlaneF
	H3DCamera_ViewportXI
	H3DCamera_ViewportYI
	H3DCamera_ViewportWidthI
	H3DCamera_ViewportHeightI
	H3DCamera_OrthoI
	H3DCamera_OccCullingI
)

/*	Enum: H3DEmitter
		The available Emitter node parameters.

	MatResI        - Material resource used for rendering
	PartEffResI    - ParticleEffect resource which configures particle properties
	MaxCountI      - Maximal number of particles living at the same time
	RespawnCountI  - Number of times a single particle is recreated after dying (-1 for infinite)
	DelayF         - Time in seconds before emitter begins creating particles (default: 0.0)
	EmissionRateF  - Maximal number of particles to be created per second (default: 0.0)
	SpreadAngleF   - Angle of cone for random emission direction (default: 0.0)
	ForceF3        - Force vector XYZ applied to particles (default: 0.0, 0.0, 0.0)
*/
const (
	_ = iota + 700
	H3DEmitter_MatResI
	H3DEmitter_PartEffResI
	H3DEmitter_MaxCountI
	H3DEmitter_RespawnCountI
	H3DEmitter_DelayF
	H3DEmitter_EmissionRateF
	H3DEmitter_SpreadAngleF
	H3DEmitter_ForceF3
)

//used for bools from c interfaces
var Bool = map[int]bool{
	0: false,
	1: true,
}
var Int = map[bool]C.int{
	false: C.int(1),
	true:  C.int(0),
}

func H3dGetVersionString() string {
	verPointer := C.h3dGetVersionString()
	//is this needed?
	defer C.free(unsafe.Pointer(verPointer))

	return C.GoString(verPointer)
}

func H3dCheckExtension(extensionName string) bool {
	cExtName := C.CString(extensionName)
	defer C.free(unsafe.Pointer(cExtName))
	return Bool[int(C.h3dCheckExtension(cExtName))]
}

func H3dGetError() bool {
	return Bool[int(C.h3dGetError())]
}

func H3dInit() bool {
	return Bool[int(C.h3dInit())]
}

func H3dRelease() {
	C.h3dRelease()
}

func H3dRender(cameraNode H3DNode) {
	C.h3dRender(C.H3DNode(cameraNode))
}

func H3dFinalizeFrame() {
	C.h3dFinalizeFrame()
}

func H3dClear() {
	C.h3dClear()
}

func H3dGetMessage(level *int, time *float32) string {
	message := C.h3dGetMessage((*C.int)(unsafe.Pointer(level)),
		(*C.float)(unsafe.Pointer(time)))
	defer C.free(unsafe.Pointer(message))
	return C.GoString(message)
}

func H3dGetOption(param int) float32 {
	return float32(C.h3dGetOption(C.int(param)))
}

func H3dSetOption(param int, value float32) bool {
	return Bool[int(C.h3dSetOption(C.int(param), C.float(value)))]
}

func H3dGetStat(param int, reset bool) float32 {
	return float32(C.h3dGetStat(C.int(param), Int[reset]))
}

func H3dShowOverlays(verts *float32,
	vertCount int,
	colR float32,
	colG float32,
	colB float32,
	colA float32,
	materialRes H3DRes,
	flags int) {
	C.h3dShowOverlays((*C.float)(unsafe.Pointer(verts)),
		C.int(vertCount),
		C.float(colR),
		C.float(colG),
		C.float(colB),
		C.float(colA),
		C.H3DRes(materialRes),
		C.int(flags))
}

func H3dClearOverlays() {
	C.h3dClearOverlays()
}

func H3dGetResType(res H3DRes) int {
	return int(C.h3dGetResType(C.H3DRes(res)))
}

func H3dGetResName(res H3DRes) string {
	return C.GoString(C.h3dGetResName(C.H3DRes(res)))
}

func H3dGetNextResource(resType int, start H3DRes) H3DRes {
	return H3DRes(C.h3dGetNextResource(C.int(resType), C.H3DRes(start)))
}

func H3dFindResource(resType int, name string) H3DRes {
	return H3DRes(C.h3dFindResource(C.int(resType), (*C.char)(unsafe.Pointer(&name))))
}

func H3dAddResource(resType int, name string, flags int) H3DRes {
	return H3DRes(C.h3dAddResource(C.int(resType), (*C.char)(unsafe.Pointer(&name)), C.int(flags)))
}

func H3dCloneResource(sourceRes H3DRes, name string) H3DRes {
	return H3DRes(C.h3dCloneResource(C.H3DRes(sourceRes), (*C.char)(unsafe.Pointer(&name))))
}

func H3dRemoveResource(res H3DRes) int {
	return int(C.h3dRemoveResource(C.H3DRes(res)))
}

func H3dIsResLoaded(res H3DRes) bool {
	return Bool[int(C.h3dIsResLoaded(C.H3DRes(res)))]
}

func H3dLoadResource(res H3DRes, data *string, size int) bool {
	return Bool[int(C.h3dLoadResource(C.H3DRes(res), (*C.char)(unsafe.Pointer(data)), C.int(size)))]

}

func H3dUnloadResource(res H3DRes) {
	C.h3dUnloadResource(C.H3DRes(res))
}

func H3dGetResElemCount(res H3DRes, elem int) int {
	return int(C.h3dGetResElemCount(C.H3DRes(res), C.int(elem)))
}

func H3dFindResElem(res H3DRes, elem int, param int, value string) int {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	return int(C.h3dFindResElem(C.H3DRes(res),
		C.int(elem), C.int(param), cValue))
}

func H3dGetResParamI(res H3DRes, elem int, elemIdx int, param int) int {
	return int(C.h3dGetResParamI(C.H3DRes(res), C.int(elem), C.int(elemIdx), C.int(param)))
}

func H3dSetResParamI(res H3DRes, elem int, elemIdx int, param int, value int) {
	C.h3dSetResParamI(C.H3DRes(res), C.int(elem), C.int(elemIdx), C.int(param), C.int(value))
}

func H3dGetResParamF(res H3DRes, elem int, elemIdx int, param int, compIdx int) float32 {
	return float32(C.h3dGetResParamF(C.H3DRes(res),
		C.int(elem), C.int(elemIdx), C.int(param), C.int(compIdx)))
}

func H3dSetResParamF(res H3DRes, elem int, elemIdx int, param int, compIdx int, value float32) {
	C.h3dSetResParamF(C.H3DRes(res),
		C.int(elem), C.int(elemIdx), C.int(param), C.int(compIdx), C.float(value))
}

func H3dGetResParamStr(res H3DRes, elem int, elemIdx int, param int) string {
	value := C.h3dGetResParamStr(C.H3DRes(res), C.int(elem), C.int(elemIdx), C.int(param))
	defer C.free(unsafe.Pointer(value))
	return C.GoString(value)
}

func H3dSetResParamStr(res H3DRes, elem int, elemIdx int, param int, value string) {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	C.h3dSetResParamStr(C.H3DRes(res), C.int(elem), C.int(elemIdx), C.int(param), cValue)
}

func H3dMapResStream(res H3DRes, elem int, elemIdx int, stream int, read bool, write bool, size int) []byte {
	//TODO: Review this method of getting the stream data from the pointer 

	cStream := C.h3dMapResStream(C.H3DRes(res), C.int(elem), C.int(elemIdx), C.int(stream),
		Int[read], Int[write])

	return C.GoBytes(unsafe.Pointer(cStream), C.int(size))
}

func H3dUnmapResStream(res H3DRes) {
	C.h3dUnmapResStream(C.H3DRes(res))
}

func H3dQueryUnloadedResource(index int) H3DRes {
	return H3DRes(C.h3dQueryUnloadedResource(C.int(index)))
}

func H3dReleaseUnusedResources() {
	C.h3dReleaseUnusedResources()
}

func H3dCreateTexture(name string, width int, height int, fmt int, flags int) H3DRes {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return H3DRes(C.h3dCreateTexture(cName, C.int(width), C.int(height), C.int(fmt), C.int(flags)))
}

func H3dSetShaderPreambles(vertPreamble string, fragPreamble string) {
	cVertPreamble := C.CString(vertPreamble)
	cFragPreamble := C.CString(fragPreamble)
	defer C.free(unsafe.Pointer(cVertPreamble))
	defer C.free(unsafe.Pointer(cFragPreamble))

	C.h3dSetShaderPreambles(cVertPreamble, cFragPreamble)
}

func H3dSetMaterialUniform(materialRes H3DRes, name string, a float32, b float32,
	c float32, d float32) bool {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return Bool[int(C.h3dSetMaterialUniform(C.H3DRes(materialRes), cName,
		C.float(a), C.float(b), C.float(c), C.float(d)))]
}

func H3dResizePipelineBuffers(pipeRes H3DRes, width int, height int) {
	C.h3dResizePipelineBuffers(C.H3DRes(pipeRes), C.int(width), C.int(height))
}

func H3dGetRenderTargetData(pipelineRes H3DRes, targetName string, bufIndex int, width *int,
	height *int, compCount *int, dataBuffer []byte) bool {
	cTargetName := C.CString(targetName)
	defer C.free(unsafe.Pointer(cTargetName))
	//TODO: Not sure if this is correct, or how Go's GC will handle it

	slcHead := (*reflect.SliceHeader)((unsafe.Pointer(&dataBuffer)))

	return Bool[int(C.h3dGetRenderTargetData(C.H3DRes(pipelineRes), cTargetName,
		C.int(bufIndex), (*C.int)(unsafe.Pointer(width)), (*C.int)(unsafe.Pointer(height)),
		(*C.int)(unsafe.Pointer(compCount)), unsafe.Pointer(slcHead.Data), C.int(len(dataBuffer))))]
}

func H3dGetNodeType(node H3DNode) int {
	return int(C.h3dGetNodeType(C.H3DNode(node)))
}

func H3dGetNodeParent(node H3DNode) H3DNode {
	return H3DNode(C.h3dGetNodeParent(C.H3DNode(node)))
}

func H3dSetNodeParent(node H3DNode, parent H3DNode) bool {
	return Bool[int(C.h3dSetNodeParent(C.H3DNode(node), C.H3DNode(parent)))]
}

func H3dGetNodeChild(node H3DNode, index int) bool {
	return Bool[int(C.h3dGetNodeChild(C.H3DNode(node), C.int(index)))]
}

func H3dAddNodes(parent H3DNode, sceneGraphRes H3DRes) H3DNode {
	return H3DNode(C.h3dAddNodes(C.H3DNode(parent), C.H3DRes(sceneGraphRes)))
}

func H3dRemoveNode(node H3DNode) {
	C.h3dRemoveNode(C.H3DNode(node))
}

func H3dCheckNodeTransFlag(node H3DNode, reset bool) bool {
	return Bool[int(C.h3dCheckNodeTransFlag(C.H3DNode(node), Int[reset]))]
}

func H3dGetNodeTransform(node H3DNode, tx *float32, ty *float32, tz *float32,
	rx *float32, ry *float32, rz *float32, sx *float32, sy *float32, sz *float32) {
	C.h3dGetNodeTransform(C.H3DNode(node), (*C.float)(unsafe.Pointer(tx)), (*C.float)(unsafe.Pointer(ty)),
		(*C.float)(unsafe.Pointer(tz)), (*C.float)(unsafe.Pointer(rx)), (*C.float)(unsafe.Pointer(ry)),
		(*C.float)(unsafe.Pointer(rz)), (*C.float)(unsafe.Pointer(sx)), (*C.float)(unsafe.Pointer(sy)),
		(*C.float)(unsafe.Pointer(sz)))
}

func H3dSetNodeTransform(node H3DNode, tx float32, ty float32, tz float32,
	rx float32, ry float32, rz float32, sx float32, sy float32, sz float32) {
	C.h3dSetNodeTransform(C.H3DNode(node), C.float(tx), C.float(ty), C.float(tz),
		C.float(rx), C.float(ry), C.float(rz), C.float(sx), C.float(sy), C.float(sz))
}

func H3dGetNodeTransMats(node H3DNode, relMat []float32, absMat []float32) {
	//TODO: Handle nil pointers, possibly check for nil
	relMatHead := (*reflect.SliceHeader)((unsafe.Pointer(&relMat)))
	absMatHead := (*reflect.SliceHeader)((unsafe.Pointer(&absMat)))

	C.h3dGetNodeTransMats(C.H3DNode(node), unsafe.Pointer(relMatHead.Data), unsafe.Pointer(absMatHead.Data))
}

//func H3dSetNodeTransMat(node H3DNode, mat4x4 []float32) {

//}
