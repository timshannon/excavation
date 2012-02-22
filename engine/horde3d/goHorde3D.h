// *************************************************************************************************
//
// Horde3D
//   Next-Generation Graphics Engine
// --------------------------------------
// Copyright (C) 2006-2009 Nicolas Schulz
//
// This software is distributed under the terms of the Eclipse Public License v1.0.
// A copy of the license may be obtained at: http://www.eclipse.org/legal/epl-v10.html
//
// *************************************************************************************************
//
//Modified by Tim Shannon for Go (http://www.golang.org) compatibility.
/*	Title: Horde3D API */

#pragma once

#ifndef DLL
#	if defined( WIN32 ) || defined( _WINDOWS )
#		define DLL extern "C" __declspec( dllimport )
#	else
#  if defined( __GNUC__ ) && __GNUC__ >= 4
#   define DLL extern "C" __attribute__ ((visibility("default")))
#  else
#		define DLL extern "C"
#  endif
#	endif
#endif


/*	Topic: Conventions
		Some conventions for the API.
	
	Horde3D uses a right-handed coordinate system, where y is the up-axis and the positive z-axis is
	pointing out of the screen. The rotations are specified in degrees and rotation direction is
	counter-clockwise when looking down the axis from the the positive end towards the origin. View
	vectors like the camera or light direction are always pointing along the negative z-axis when no
	transformation is applied. Matrices in Horde are stored in a column-major memory layout. When
	Euler angles are used, the rotation order is YXZ [*vec], so the z-rotation is applied first.
*/


// *************************************************************************************************
// Core Library
// *************************************************************************************************

/* Group: Typedefs and constants */

/*	Constants: Typedefs
	H3DRes   - handle to resource (type: int32)
	H3DNode  - handle to scene node (type: int32)
*/
typedef int H3DRes;
typedef int H3DNode;


/*	Constants: Predefined constants
	H3DRootNode  - Scene root node handle
*/
const H3DNode H3DRootNode = 1;



/* Group: Basic functions */
/* Function: h3dGetVersionString
		Returns the engine version string.
	
	Details:
		This function returns a pointer to a string containing the current version of Horde3D.
	
	Parameters:
		none
		
	Returns:
		pointer to the version string
*/
const char *h3dGetVersionString();

/* Function: h3dCheckExtension
		Checks if an extension is part of the engine library.
	
	Details:
		This function checks if a specified extension is contained in the DLL/shared object of the engine.
	
	Parameters:
		extensionName  - name of the extension
		
	Returns:
		true if extension is implemented, otherwise false
*/
int h3dCheckExtension( const char *extensionName );

/* Function: h3dGetError
		Checks if an error occured.
	
	Details:
		This function checks if an error occured in a previous API function call. If an error
		flag is set, the function resets the flag and returns true. The function will solely
		report errors that originate from a violated precondition, like an invalid parameter
		that is passed to an API function. Errors that happen during the execution of a function,
		for example failure of initializing the engine due to a missing hardware feature, can
		be catched by checking the return value of the corresponding API function.
		More information about the error can be retrieved by checking the message queue,
		provided that the message level is set accordingly.
	
	Parameters:
		none
		
	Returns:
		true in there was an error, otherwise false
*/
int h3dGetError();

/* Function: h3dInit
		Initializes the engine.
	
	Details:
		This function initializes the graphics engine and makes it ready for use. It has to be the
		first call to the engine except for getVersionString. In order to successfully initialize
		the engine the calling application must provide a valid OpenGL context. The function can be
		called several times on different rendering contexts in order to initialize them.
	
	Parameters:
		none
		
	Returns:
		true in case of success, otherwise false
*/
int h3dInit();

/* Function: h3dRelease
		Releases the engine.
	
	Details:
		This function releases the engine and frees all objects and associated memory. It should be
		called when the application is destroyed.
	
	Parameters:
		none
		
	Returns:
		nothing
*/
void h3dRelease();


/* Function: h3dRender
		Main rendering function.
	
	Details:
		This is the main function of the engine. It executes all the rendering, animation and other
		tasks. The function can be called several times per frame, for example in order to write to different
		output buffers.
	
	Parameters:
		cameraNode  - camera node used for rendering scene
		
	Returns:
		nothing
*/
void h3dRender( H3DNode cameraNode );

/* Function: h3dFinalizeFrame
		Marker for end of frame.
	
	Details:
		This function tells the engine that the current frame is finished and that all
		subsequent rendering operations will be for the next frame.
	
	Parameters:
		none
		
	Returns:
		nothing
*/
void h3dFinalizeFrame();

/* Function: h3dClear
		Removes all resources and scene nodes.
	
	Details:
		This function removes all nodes from the scene graph except the root node and releases all resources.
	
		*Warning: All resource and node IDs are invalid after calling this function*
	
	Parameters:
		none
		
	Returns:
		nothing
*/
void h3dClear();


// --- General functions ---
/* Function: h3dGetMessage
		Gets the next message from the message queue.
	
	Details:
		This function returns the next message string from the message queue and writes additional information
		to the specified variables. If no message is left over in the queue an empty string is returned.
	
	Parameters:
		level  - pointer to variable for storing message level indicating importance (can be NULL)
		time   - pointer to variable for storing time when message was added (can be NULL)
		
	Returns:
		message string or empty string if no message is in queue
*/
const char *h3dGetMessage( int *level, float *time );

/* Function: h3dGetOption
		Gets an option parameter of the engine.
	
	Details:
		This function gets a specified option parameter and returns its value.
	
	Parameters:
		param  - option parameter
		
	Returns:
		current value of the specified option parameter
*/
float h3dGetOption(int H3DOptionsList );

/* Function: h3dSetOption
		Sets an option parameter for the engine.
	
	Details:
		This function sets a specified option parameter to a specified value.
	
	Parameters:
		param  - option parameter
		value  - value of the option parameter
		
	Returns:
		true if the option could be set to the specified value, otherwise false
*/
int h3dSetOption(int H3DOptionsList, float value );

/* Function: h3dGetStat
		Gets a statistic value of the engine.
	
	Details:
		This function returns the value of the specified statistic. The reset flag makes it possible
		to reset the statistic value after reading.
	
	Parameters:
		param  - statistic parameter
		reset  - flag specifying whether statistic value should be reset
		
	Returns:
		current value of the specified statistic parameter
*/
float h3dGetStat(int H3DStatsList , int reset );

/* Function: h3dShowOverlays
		Displays overlays on the screen.
	
	Details:
		This function displays one or more overlays with a specified material and color.
		An overlay is a screen-space quad that can be used to render 2D GUI elements. The overlay coordinate
		system has its origin (0, 0) at the top-left corner of the screen and its maximum (aspect, 1)
		at the bottom-right corner. As the x coordinate of the maximum corresponds to the aspect ratio
		of the viewport, the size of overlays can always be the same, even when different screen formats
		(standard 4:3, widescreen 16:9, etc.) are used. Texture coordinates are using a system where the
		coordinates (0, 0) correspond to the lower left corner of the image.
		Overlays are drawn in the order in which they are pushed using this function. Overlays with
		the same state will be batched together, so it can make sense to group overlays that have the
		same material, color and flags in order to achieve best performance.
		Note that the overlays have to be removed manually using the function h3dClearOverlays.
	
	Parameters:
		verts                   - vertex data (x, y, u, v), interpreted as quads
		vertCount               - number of vertices (must be multiple of 4)
		colR, colG, colB, colA  - color (and transparency) of overlays
		materialRes             - material resource used for rendering
		flags                   - overlay flags (reserved for future use)
		
	Returns:
		nothing
*/
void h3dShowOverlays( const float *verts, int vertCount, float colR, float colG, float colB,
                          float colA, H3DRes materialRes, int flags );

/* Function: h3dClearOverlays
		Removes all overlays.
	
	Details:
		This function removes all overlays that were added with h3dShowOverlays.
	
	Parameters:
		none
		
	Returns:
		nothing
*/
void h3dClearOverlays();


/* Group: General resource management functions */
/* Function: h3dGetResType
		Returns the type of a resource.
	
	Details:
		This function returns the type of a specified resource. If the resource handle is invalid,
		the function returns the resource type 'Unknown'.
	
	Parameters:
		res  - handle to the resource
		
	Returns:
		type of the resource
*/
int h3dGetResType( H3DRes res );

/* Function: h3dGetResName
		Returns the name of a resource.
	
	Details:
		This function returns a pointer to the name of a specified resource. If the resource handle
		is invalid, the function returns an empty string.

		*Important Note: The pointer is const and allows only read access to the data. Do never try to modify the
		data of the pointer since that can corrupt the engine's internal states!*
	
	Parameters:
		res  - handle to the resource
		
	Returns:
		name of the resource or empty string in case of failure
*/
const char *h3dGetResName( H3DRes res );

/* Function: h3dGetNextResource
		Returns the next resource of the specified type.
	
	Details:
		This function searches the next resource of the specified type and returns its handle.
		The search begins after the specified start handle. If a further resource of the queried type
		does not exist, a zero handle is returned. The function can be used to iterate over all
		resources of a given type by using as start the return value of the previous iteration step.
		The first iteration step should start at 0 and iteration can be ended when the function returns 0.
	
	Parameters:
		type   - type of resource to be searched (H3DResTypes::Undefined for all types)
		start  - resource handle after which the search begins (can be 0 for beginning of resource list)
		
	Returns:
		handle to the found resource or 0 if it does not exist
*/
H3DRes h3dGetNextResource( int type, H3DRes start );

/* Function: h3dFindResource
		Finds a resource and returns its handle.
	
	Details:
		This function searches the resource of the specified type and name and returns its handle. If
		the resource is not available in the resource manager a zero handle is returned.
	
	Parameters:
		type  - type of the resource
		name  - name of the resource
		
	Returns:
		handle to the resource or 0 if not found
*/
H3DRes h3dFindResource( int type, const char *name );

/* Function: h3dAddResource
		Adds a resource.
	
	Details:
		This function tries to add a resource of a specified type and name to the resource manager. If
		a resource of the same type and name is already existing, the handle to the existing resource is returned
		instead of creating a new one and the user reference count of the resource is increased.
	
	Parameters:
		type   - type of the resource
		name   - name of the resource
		flags  - flags used for creating the resource
		
	Returns:
		handle to the resource to be added or 0 in case of failure
*/
H3DRes h3dAddResource( int type, const char *name, int flags );

/* Function: h3dCloneResource
		Duplicates a resource.
	
	Details:
		This function duplicates a specified resource. In the cloning process a new resource with the
		specified name is added to the resource manager and filled with the data of the specified source
		resource. If the specified name for the new resource is already in use, the function fails and
		returns 0. If the name string is empty, a unique name for the resource is generated automatically.
	
	Parameters:
		sourceRes  - handle to resource to be cloned
		name       - name of new resource (can be empty for auto-naming)
		
	Returns:
		handle to the cloned resource or 0 in case of failure
*/
H3DRes h3dCloneResource( H3DRes sourceRes, const char *name );

/* Function: h3dRemoveResource
		Removes a resource.
	
	Details:
		This function decreases the user reference count of a specified resource. When the user reference
		count is zero and there are no internal references, the resource can be released and removed using
		the API fuction releaseUnusedResources.
	
	Parameters:
		res	- handle to the resource to be removed
		
	Returns:
		the number of references that the application is still holding after removal or -1 in case of an error
*/
int h3dRemoveResource( H3DRes res );

/* Function: h3dIsResLoaded
		Checks if a resource is loaded.
	
	Details:
		This function checks if the specified resource has been successfully loaded.
	
	Parameters:
		res  - handle to the resource to be checked
		
	Returns:
		true if resource is loaded, otherwise or in case of failure false
*/
int h3dIsResLoaded( H3DRes res );

/* Function: h3dLoadResource
		Loads a resource.
	
	Details:
		This function loads data for a resource that was previously added to the resource manager.
		If data is a NULL-pointer the resource manager is told that the resource doesn't have any data
		(e.g. the corresponding file was not found). In this case, the resource remains in the unloaded state
		but is no more returned when querying unloaded resources. When the specified resource is already loaded,
		the function returns false.
	
	Parameters:
		res   - handle to the resource for which data will be loaded
		data  - pointer to the data to be loaded
		size  - size of the data block
		
	Returns:
		true in case of success, otherwise false
*/
int h3dLoadResource( H3DRes res, const char *data, int size );

/* Function: h3dUnloadResource
		Unloads a resource.
	
	Details:
		This function unloads a previously loaded resource and restores the default values
		it had before loading. The state is set back to unloaded which makes it possible to load
		the resource again.
	
	Parameters:
		res  - handle to resource to be unloaded
		
	Returns:
		nothing
*/
void h3dUnloadResource( H3DRes res );


/* Function: h3dGetResElemCount
		Gets the number of elements in a resource.
	
	Details:
		This function returns how many elements of the specified element type a specified
		resource has.
	
	Parameters:
		res      - handle to the resource to be accessed
		elem     - element type
		
	Returns:
		number of elements
*/
int h3dGetResElemCount( H3DRes res, int elem );

/* Function: h3dFindResElem
		Finds a resource element with the specified property value.
	
	Details:
		This function searches in a specified resource for the first element of the specified
		type that has the property with the specified name set to the specified search value.
		If such element is found, its index is returned, otherwise the function returns -1. All
		string comparisons done for the search are case-sensitive.
	
	Parameters:
		res      - handle to the resource to be accessed
		elem     - element type
		param    - parameter name
		value    - parameter value to be searched for
		
	Returns:
		index of element or -1 if element not found
*/
int h3dFindResElem( H3DRes res, int elem, int param, const char *value );

/* Function: h3dGetResParamI
		Gets an integer property of a resource element.
	
	Details:
		This function returns a specified property of the specified resource element.
		The property must be of the type int.
	
	Parameters:
		res      - handle to the resource to be accessed
		elem     - element type
		elemIdx  - index of element
		param    - parameter to be accessed
		
	Returns:
		value of the parameter
*/
int h3dGetResParamI( H3DRes res, int elem, int elemIdx, int param );

/* Function: h3dSetResParamI
		Sets an integer property of a resource element.
	
	Details:
		This function sets a specified property of the specified resource element to the
		specified value. The property must be of the type int.
	
	Parameters:
		res      - handle to the resource to be modified
		elem     - element type
		elemIdx  - index of element
		param    - parameter to be modified
		value    - new value for the specified parameter
		
	Returns:
		 nothing
*/
void h3dSetResParamI( H3DRes res, int elem, int elemIdx, int param, int value );

/* Function: h3dGetResParamF
		Gets a float property of a resource element.
	
	Details:
		This function returns the specified component of the specified resource element property.
		The property must be of the type float. Float properties can have several indexed components,
		for example the x, y and z components of a vector.
	
	Parameters:
		res      - handle to the resource to be accessed
		elem     - element type
		elemIdx  - index of element
		param    - parameter to be accessed
		compIdx  - component of the parameter to be accessed
		
	Returns:
		value of the parameter
*/
float h3dGetResParamF( H3DRes res, int elem, int elemIdx, int param, int compIdx );

/* Function: h3dSetResParamF
		Sets a float property of a resource element.
	
	Details:
		This function sets the specified component of the specified resource element property to the
		specified value. The property must be of the type float.
	
	Parameters:
		res      - handle to the resource to be modified
		elem     - element type
		elemIdx  - index of element
		param    - parameter to be modified
		compIdx  - component of the parameter to be modified
		value    - new value for the specified parameter
		
	Returns:
		nothing
*/
void h3dSetResParamF( H3DRes res, int elem, int elemIdx, int param, int compIdx, float value );

/* Function: h3dGetResParamStr
		Gets a string property of a resource element.
	
	Details:
		This function returns a specified property of the specified resource element.
		The property must be of the type string (const char *).

		*Important Note: The pointer is const and allows only read access to the data. Do never try to modify the
		data of the pointer since that can corrupt the engine's internal states!*
	
	Parameters:
		res      - handle to the resource to be accessed
		elem     - element type
		elemIdx  - index of element
		param    - parameter to be accessed
		
	Returns:
		value of the property or empty string if no such property exists
*/
const char *h3dGetResParamStr( H3DRes res, int elem, int elemIdx, int param );

/* Function: h3dSetResParamStr
		Sets a string property of a resource element.
	
	Details:
		This function sets a specified property of the specified resource element to the
		specified value. The property must be of the type string (const char *).
	
	Parameters:
		res      - handle to the resource to be modified
		elem     - element type
		elemIdx  - index of element
		param    - parameter to be modified
		value    - new value for the specified parameter
		
	Returns:
		nothing
*/
void h3dSetResParamStr( H3DRes res, int elem, int elemIdx, int param, const char *value );

/* Function: h3dMapResStream
		Maps the stream of a resource element.
	
	Details:
		This function maps the specified stream of a specified resource element and returns a
		pointer to the stream data. The required access to the data can be specified with the read
		write parameters. If read is false, the pointer will usually not contain meaningful data.
		Not all resource streams can be mapped with both read and write access. If it is not
		possible to map the stream, the function will return a NULL pointer. A mapped stream should
		be unmapped again as soon as possible but always before subsequent API calls are made. It
		is only possible to map one stream per resource at a time.
	
	Parameters:
		res      - handle to the resource to be mapped
		elem     - element type
		elemIdx  - index of element
		stream   - stream to be mapped
		read     - flag indicating whether read access is required
		write    - flag indicating whether write access is required
		
	Returns:
		pointer to stream data or NULL if stream cannot be mapped
*/
void *h3dMapResStream( H3DRes res, int elem, int elemIdx, int stream, int read, int write );

/* Function: h3dUnmapResStream
		Unmaps a previously mapped resource stream.
	
	Details:
		This function unmaps a resource stream that has been mapped before.
	
	Parameters:
		res  - handle to the resource to be unmapped
		
	Returns:
		nothing
*/
void h3dUnmapResStream( H3DRes res );

/* Function: h3dQueryUnloadedResource
		Returns handle to an unloaded resource.
	
	Details:
		This function looks for a resource that is not yet loaded and returns its handle.
		If there are no unloaded resources or the zero based index specified is greater than the number
		of the currently unloaded resources, 0 is returned.
	
	Parameters:
		index  - index of unloaded resource within the internal list of unloaded resources (starting with 0) 
		
	Returns:
		handle to an unloaded resource or 0
*/
H3DRes h3dQueryUnloadedResource( int index );

/* Function: h3dReleaseUnusedResources
		Frees resources that are no longer used.
	
	Details:
		This function releases resources that are no longer used. Unused resources were either told
		to be released by the user calling removeResource or are no more referenced by any other
		engine objects.
	
	Parameters:
		none
		
	Returns:
		nothing
*/
void h3dReleaseUnusedResources();


/* Group: Specific resource management functions */
/* Function: h3dCreateTexture
		Creates a Texture resource.
	
	Details:
		This function tries to create and add a Texture resource with the specified name to the resource
		manager. If a Texture resource with the same name is already existing, the function fails. The
		texture is initialized with the specified dimensions and the resource is declared as loaded. This
		function is especially useful to create dynamic textures (e.g. for displaying videos) or output buffers
		for render-to-texture.
	
	Parameters:
		name        - name of the resource
		width       - width of the texture image
		height      - height of the texture image
		fmt         - texture format (see stream formats)
		flags       - flags used for creating the resource
		
	Returns:
		handle to the created resource or 0 in case of failure
*/
H3DRes h3dCreateTexture( const char *name, int width, int height, int fmt, int flags );

/* Function: h3dSetShaderPreambles
		Sets preambles of all Shader resources.
	
	Details:
		This function defines a header that is inserted at the beginning of all shaders. The preamble
		is used when a shader is compiled, so changing it will not affect any shaders that are already
		compiled. The preamble is useful for setting platform-specific defines that can be employed for
		creating several shader code paths, e.g. for supporting different hardware capabilities.
	
	Parameters:
		vertPreamble  - preamble text of vertex shaders (default: empty string)
		fragPreamble  - preamble text of fragment shaders (default: empty string)
		
	Returns:
		nothing
*/
void h3dSetShaderPreambles( const char *vertPreamble, const char *fragPreamble );

/* Function: h3dSetMaterialUniform
		Sets a shader uniform of a Material resource.
	
	Details:
		This function sets the specified shader uniform of the specified material to the specified values.
	
	Parameters:
		materialRes  - handle to the Material resource to be accessed
		name         - name of the uniform as defined in Material resource
		a, b, c, d   - values of the four components
		
	Returns:
		true if uniform was found, otherwise false
*/
int h3dSetMaterialUniform( H3DRes materialRes, const char *name, float a, float b, float c, float d );

/* Function: h3dResizePipelineBuffers
		Changes the size of the render targets of a pipeline.
	
	Details:
		This function sets the base width and height which affects render targets with relative (in percent) size 
		specification. Changing the base size is usually desired after engine initialization and when the window
		is being resized. Note that in case several cameras use the same pipeline resource instance, the change
		will affect all cameras.
	
	Parameters:
	    pipeRes  - the pipeline resource instance to be changed
		width    - base width in pixels used for render targets with relative size
		height   - base height in pixels used for render targets with relative size
		
	Returns:
		nothing
*/
void h3dResizePipelineBuffers( H3DRes pipeRes, int width, int height );

/* Function: h3dGetRenderTargetData
		Reads back the pixel data of a render target buffer.
	
	Details:
		This function reads back the pixels of a specified render target. If no pipeline resource is specified,
		the color data of the backbuffer is read back as RGBA float. Otherwise, the specified buffer of the
		specified render target is read back and stored in dataBuffer as RGBA float values.
		To compute the required size of the output buffer to which the data is written, the function can be called
		with a NULL pointer for dataBuffer and pointers to variables where the buffer width, height and the number
		of components will be stored.
		As this function has a considerable performance overhead, it is only intended for debugging purposes and screenshots.
		For more information about the render buffers, refer to the Pipeline documentation.
		
	Parameters:
		pipelineRes  - handle to pipeline resource (0 for backbuffer)
		targetName   - name of render target to be accessed (ignored for backbuffer)
		bufIndex     - index of buffer to be accessed (32 for depth buffer)
		width        - pointer to variable where the width of the buffer will be stored (can be NULL)
		height       - pointer to variable where the height of the buffer will be stored (can be NULL)
		compCount    - pointer to variable where the number of components will be stored (can be NULL)
		dataBuffer   - pointer to array where the pixel data will be stored (can be NULL)
		bufferSize   - size of dataBuffer array in bytes
		
	Returns:
		true if specified render target could be found, otherwise false
*/
int h3dGetRenderTargetData( H3DRes pipelineRes, const char *targetName, int bufIndex,
                                 int *width, int *height, int *compCount, void *dataBuffer, int bufferSize );


/* Group: General scene graph functions */
/* Function: h3dGetNodeType
		Returns the type of a scene node.
	
	Details:
		This function returns the type of a specified scene node. If the node handle is invalid,
		the function returns the node type 'Unknown'.
	
	Parameters:
		node  - handle to the scene node
		
	Returns:
		type of the scene node
*/
int h3dGetNodeType( H3DNode node );
	
/* Function: h3dGetNodeParent
		Returns the parent of a scene node.
	
	Details:
		This function returns the handle to the parent node of a specified scene node. If the specified
		node handle is invalid or the root node, 0 is returned.
	
	Parameters:
		node  - handle to the scene node
		
	Returns:
		handle to parent node or 0 in case of failure
*/
H3DNode h3dGetNodeParent( H3DNode node );


/* Function: h3dSetNodeParent
		Relocates a node in the scene graph.
	
	Details:
		This function relocates a scene node. It detaches the node from its current parent and attaches
		it to the specified new parent node. If the attachment to the new parent is not possible, the
		function returns false. Relocation is not possible for the H3DRootNode.
	
	Parameters:
		node    - handle to the scene node to be relocated
		parent  - handle to the new parent node
		
	Returns:
		true if node could be relocated, otherwise false
*/
int h3dSetNodeParent( H3DNode node, H3DNode parent );

/* Function: h3dGetNodeChild
		Returns the handle to a child node.
	
	Details:
		This function looks for the n-th (index) child node of a specified node and returns its handle. If the child
		doesn't exist, the function returns 0.
	
	Parameters:
		node   - handle to the parent node
		index  - index of the child node
		
	Returns:
		handle to the child node or 0 if child doesn't exist
*/
H3DNode h3dGetNodeChild( H3DNode node, int index );



/* Function: h3dAddNodes
		Adds nodes from a SceneGraph resource to the scene.
	
	Details:
		This function creates several new nodes as described in a SceneGraph resource and
		attaches them to a specified parent node. If an invalid scenegraph resource is specified
		or the scenegraph resource is unloaded, the function returns 0.
	
	Parameters:
		parent         - handle to parent node to which the root of the new nodes will be attached
		sceneGraphRes  - handle to loaded SceneGraph resource
		
	Returns:
		handle to the root of the created nodes or 0 in case of failure
*/
H3DNode h3dAddNodes( H3DNode parent, H3DRes sceneGraphRes );

/* Function: h3dRemoveNode
		Removes a node from the scene.
	
	Details:
		This function removes the specified node and all of it's children from the scene.
	
	Parameters:
		node  - handle to the node to be removed
		
	Returns:
		nothing
*/
void h3dRemoveNode( H3DNode node );

/* Function: h3dCheckNodeTransFlag
		Checks if a scene node has been transformed by the engine.
	
	Details:
		This function checks if a scene node has been transformed by the engine since the last
		time the transformation flag was reset. Therefore, it stores a flag that is set to true when a
		setTransformation function is called explicitely by the application or when the node transformation
		has been updated by the animation system. The function also makes it possible to reset the
		transformation flag.
	
	Parameters:
		node   - handle to the node to be accessed
		reset  - flag indicating whether transformation flag shall be reset
		
	Returns:
		true if node has been transformed, otherwise false
*/
int h3dCheckNodeTransFlag( H3DNode node, int reset );

/* Function: h3dGetNodeTransform
		Gets the relative transformation of a node.
	
	Details:
		This function gets the translation, rotation and scale of a specified scene node object. The
		coordinates are in local space and contain the transformation of the node relative to its parent.
	
	Parameters:
		node        - handle to the node which will be accessed
		tx, ty, tz  - pointers to variables where translation of the node will be stored (can be NULL)
		rx, ry, rz  - pointers to variables where rotation of the node in Euler angles
		              will be stored (can be NULL)
		sx, sy, sz  - pointers to variables where scale of the node will be stored (can be NULL)
		
	Returns:
		nothing
*/
void h3dGetNodeTransform( H3DNode node, float *tx, float *ty, float *tz,
                              float *rx, float *ry, float *rz, float *sx, float *sy, float *sz );

/* Function: h3dSetNodeTransform
		Sets the relative transformation of a node.
	
	Details:
		This function sets the relative translation, rotation and scale of a specified scene node object.
		The coordinates are in local space and contain the transformation of the node relative to its parent.
	
	Parameters:
		node        - handle to the node which will be modified
		tx, ty, tz  - translation of the node
		rx, ry, rz  - rotation of the node in Euler angles
		sx, sy, sz  - scale of the node
		
	Returns:
		nothing
*/
void h3dSetNodeTransform( H3DNode node, float tx, float ty, float tz,
                              float rx, float ry, float rz, float sx, float sy, float sz );

/* Function: h3dGetNodeTransMats
		Returns the transformation matrices of a node.
	
	Details:
		This function stores a pointer to the relative and absolute transformation matrices
		of the specified node in the specified pointer varaibles.
	
	Parameters:
		node    - handle to the scene node to be accessed
		relMat  - pointer to a variable where the address of the relative transformation matrix will be stored
		          (can be NULL if matrix is not required)
		absMat  - pointer to a variable where the address of the absolute transformation matrix will be stored
		          (can be NULL if matrix is not required)
		
		
	Returns:
		nothing
*/
void h3dGetNodeTransMats( H3DNode node, const float **relMat, const float **absMat );

/* Function: h3dSetNodeTransMat
		Sets the relative transformation matrix of a node.
	
	Details:
		This function sets the relative transformation matrix of the specified scene node. It is basically the
		same as setNodeTransform but takes directly a matrix instead of individual transformation parameters.
	
	Parameters:
		node    - handle to the node which will be modified
		mat4x4  - pointer to a 4x4 matrix in column major order
		
	Returns:
		nothing
*/
void h3dSetNodeTransMat( H3DNode node, const float *mat4x4 );

/* Function: h3dGetNodeParamI
		Gets a property of a scene node.
	
	Details:
		This function returns a specified property of the specified node.
		The property must be of the type int or H3DRes.
	
	Parameters:
		node   - handle to the node to be accessed
		param  - parameter to be accessed
		
	Returns:
		value of the parameter
*/
int h3dGetNodeParamI( H3DNode node, int param );

/* Function: h3dSetNodeParamI
		Sets a property of a scene node.
	
	Details:
		This function sets a specified property of the specified node to a specified value.
		The property must be of the type int or H3DRes.
	
	Parameters:
		node   - handle to the node to be modified
		param  - parameter to be modified
		value  - new value for the specified parameter
		
	Returns:
		nothing
*/
void h3dSetNodeParamI( H3DNode node, int param, int value );

/* Function: h3dGetNodeParamF
		Gets a property of a scene node.
	
	Details:
		This function returns a specified property of the specified node. Float properties can have
		several indexed components, for example the x, y and z components of a vector.
		The property must be of the type float.
	
	Parameters:
		node     - handle to the node to be accessed
		param    - parameter to be accessed
		compIdx  - component of the parameter to be accessed
		
	Returns:
		value of the parameter
*/
float h3dGetNodeParamF( H3DNode node, int param, int compIdx );

/* Function: h3dSetNodeParamF
		Sets a property of a scene node.
	
	Details:
		This function sets a specified property of the specified node to a specified value.
		The property must be of the type float.
	
	Parameters:
		node     - handle to the node to be modified
		param    - parameter to be modified
		compIdx  - component of the parameter to be modified
		value    - new value for the specified parameter
		
	Returns:
		nothing
*/
void h3dSetNodeParamF( H3DNode node, int param, int compIdx, float value );

/* Function: h3dGetNodeParamStr
		Gets a property of a scene node.
	
	Details:
		This function returns a specified property of the specified node.
		The property must be of the type string (const char *).

		*Important Note: The pointer is const and allows only read access to the data. Do never try to modify the
		data of the pointer since that can corrupt the engine's internal states!*
	
	Parameters:
		node   - handle to the node to be accessed
		param  - parameter to be accessed
		
	Returns:
		value of the property or empty string if no such property exists
*/
const char *h3dGetNodeParamStr( H3DNode node, int param );

/* Function: h3dSetNodeParamStr
		Sets a property of a scene node.
	
	Details:
		This function sets a specified property of the specified node to a specified value.
		The property must be of the type string (const char *).
	
	Parameters:
		node   - handle to the node to be modified
		param  - parameter to be modified
		value  - new value for the specified parameter
		
	Returns:
		nothing
*/
void h3dSetNodeParamStr( H3DNode node, int param, const char *value );

/* Function: h3dGetNodeFlags
		Gets the scene node flags.

	Details:
		This function returns a bit mask containing the set scene node flags.
	
	Parameters:
		node  - handle to the node to be accessed
		
	Returns:
		flag bitmask
*/
int h3dGetNodeFlags( H3DNode node );

/* Function: h3dSetNodeFlags
		Sets the scene node flags.
	
	Details:
		This function sets the flags of the specified scene node.
	
	Parameters:
		node       - handle to the node to be modified
		flags      - new flag bitmask
		recursive  - specifies whether flags should be applied recursively to all child nodes
		
	Returns:
		nothing
*/
void h3dSetNodeFlags( H3DNode node, int flags, int recursive );

/* Function: h3dGetNodeAABB
		Gets the bounding box of a scene node.
	
	Details:
		This function stores the world coordinates of the axis aligned bounding box of a specified node in
		the specified variables. The bounding box is represented using the minimum and maximum coordinates
		on all three axes.
	
	Parameters:
		node              - handle to the node which will be accessed
		minX, minY, minZ  - pointers to variables where minimum coordinates will be stored
		maxX, maxY, maxZ  - pointers to variables where maximum coordinates will be stored
		
	Returns:
		nothing
*/
void h3dGetNodeAABB( H3DNode node, float *minX, float *minY, float *minZ,
                         float *maxX, float *maxY, float *maxZ );

/* Function: h3dFindNodes
		Finds scene nodes with the specified properties.
	
	Details:
		This function loops recursively over all children of startNode and adds them to an internal list
		of results if they match the specified name and type. The result list is cleared each time this
		function is called. The function returns the number of nodes which were found and added to the list.
	
	Parameters:
		startNode  - handle to the node where the search begins
		name       - name of nodes to be searched (empty string for all nodes)
		type       - type of nodes to be searched (H3DNodeTypes::Undefined for all types)
		
	Returns:
		number of search results
*/
int h3dFindNodes( H3DNode startNode, const char *name, int type );

/* Function: h3dGetNodeFindResult
		Gets a result from the findNodes query.
	
	Details:
		This function returns the n-th (index) result of a previous findNodes query. The result is the handle
		to a scene node with the poperties specified at the findNodes query. If the index doesn't exist in the
		result list the function returns 0.
	
	Parameters:
		index  - index of search result
		
	Returns:
		handle to scene node from findNodes query or 0 if result doesn't exist
*/
H3DNode h3dGetNodeFindResult( int index );

/* Function: h3dCastRay
		Performs a recursive ray collision query.
	
	Details:
		This function checks recursively if the specified ray intersects the specified node or one of its children.
		The function finds intersections relative to the ray origin and returns the number of intersecting scene
		nodes. The ray is a line segment and is specified by a starting point (the origin) and a finite direction
		vector which also defines its length. Currently this function is limited to returning intersections with Meshes.
		For Meshes, the base LOD (LOD0) is always used for performing the ray-triangle intersection tests.
	
	Parameters:
		node        - node at which intersection check is beginning
		ox, oy, oz  - ray origin
		dx, dy, dz  - ray direction vector also specifying ray length
		numNearest  - maximum number of intersection points to be stored (0 for all)
		
	Returns:
		number of intersections
	*/
int h3dCastRay( H3DNode node, float ox, float oy, float oz, float dx, float dy, float dz, int numNearest );

/*	Function: h3dGetCastRayResult
		Returns a result of a previous castRay query.

	Details:
		This functions is used to access the results of a previous castRay query. The index is used to access
		a specific result. The intersection data is copied to the specified variables.

	Parameters:
		index         - index of result to be accessed (range: 0 to number of results returned by castRay)
		node          - handle of intersected node
		distance      - distance from ray origin to intersection point
		intersection  - coordinates of intersection point (float[3] array)

	Returns:
		true if index was valid and data could be copied, otherwise false
*/
int h3dGetCastRayResult( int index, H3DNode *node, float *distance, float *intersection );

/*	Function: h3dCheckNodeVisibility
		Checks if a node is visible.

	Details:
		This function checks if a specified node is visible from the perspective of a specified
		camera. The function always checks if the node is in the camera's frustum. If checkOcclusion
		is true, the function will take into account the occlusion culling information from the previous
		frame (if occlusion culling is disabled the flag is ignored). The flag calcLod determines whether the
		detail level for the node should be returned in case it is visible. The function returns -1 if the node
		is not visible, otherwise 0 (base LOD level) or the computed LOD level.

	Parameters:
		node            - node to be checked for visibility
		cameraNode      - camera node from which the visibility test is done
		checkOcclusion  - specifies if occlusion info from previous frame should be taken into account
		calcLod         - specifies if LOD level should be computed

	Returns:
		computed LOD level or -1 if node is not visible
*/
int h3dCheckNodeVisibility( H3DNode node, H3DNode cameraNode, int checkOcclusion, int calcLod );


/* Group: Group-specific scene graph functions */
/* Function: h3dAddGroupNode
		Adds a Group node to the scene.
	
	Details:
		This function creates a new Group node and attaches it to the specified parent node.
	
	Parameters:
		parent  - handle to parent node to which the new node will be attached
		name    - name of the node
		
	Returns:
		handle to the created node or 0 in case of failure
*/
H3DNode h3dAddGroupNode( H3DNode parent, const char *name );


/* Group: Model-specific scene graph functions */
/* Function: h3dAddModelNode
		Adds a Model node to the scene.
	
	Details:
		This function creates a new Model node and attaches it to the specified parent node.
	
	Parameters:
		parent       - handle to parent node to which the new node will be attached
		name         - name of the node
		geometryRes  - Geometry resource used by Model node
		
	Returns:
		handle to the created node or 0 in case of failure
*/
H3DNode h3dAddModelNode( H3DNode parent, const char *name, H3DRes geometryRes );

/* Function: h3dSetupModelAnimStage
		Configures an animation stage of a Model node.
	
	Details:
		This function is used to setup the specified animation stage/slot of the specified Model node.
		
		The function can be used for animation blending and mixing. There is a fixed number of stages
		(by default 16) on which different animations can be played. The start node determines the first node
		(Joint or Mesh) to which the animation is recursively applied. If the start node is an empty string, the
		animation affects all animatable nodes (Joints and Meshes) of the model. If a NULL-handle is used for
		animationRes, the stage is cleared and the currently set animation is removed.

		The layer determines the priority of the animation and how the weights are distributed. See
		h3dSetModelAnimParams for more information.
		
		A simple way to do animation mixing is using additive animations. If a stage is configured to be
		additive, the engine calculates the difference between the current frame and the first frame in the
		animation and adds this delta, scaled by the weight factor, to the current transformation of the joints
		or meshes. Additive animations completely ignore the layer settings. They are only applied by the engine
		if a non-additive animation is assigned to the model as well.
	
	Parameters:
		modelNode     - handle to the Model node to be modified
		stage         - index of the animation stage to be configured
		animationRes  - handle to Animation resource (can be 0)
		layer         - layer id
		startNode     - name of first node to which animation shall be applied (or empty string)
		additive      - flag indicating whether stage is additive
		
	Returns:
		nothing
*/
void h3dSetupModelAnimStage( H3DNode modelNode, int stage, H3DRes animationRes, int layer,
                                 const char *startNode, int additive );

/* Function: h3dSetModelAnimParams
		Sets the animation stage parameters of a Model node.
	
	Details:
		This function sets the current animation time and weight for a specified stage of the
		specified model. The time corresponds to the frames of the animation and the animation is
		looped if the time is higher than the maximum number of frames in the Animation resource.
		The weight is used for animation blending and determines how much influence the stage has compared
		to the other active stages. All weights of stages that are on the same layer are normalized.
		Stages with a higher layer id are animated before stages with lower layers. The blend weight
		is distributed across the layers. If the weight sum of a layer is smaller than 1.0, the remaining
		weight is propagated to the next lower layer. So if a layer uses a weight of 100%, the lower layers
		will not get any weight and consequently not contribute to the overall animation.
	
	Parameters:
		modelNode  - handle to the Model node to be modified
		stage      - index of the animation stage to be modified
		time       - new animation time/frame
		weight     - new blend weight
		
	Returns:
		nothing
*/
void h3dSetModelAnimParams( H3DNode modelNode, int stage, float time, float weight );

/* Function: h3dSetModelMorpher
		Sets the weight of a morph target.
	
	Details:
		This function sets the weight of a specified morph target. If the target parameter
		is an empty string the weight of all morph targets in the specified Model node is modified.
		If the specified morph target is not found the function returns false.
	
	Parameters:
		modelNode  - handle to the Model node to be modified
		target     - name of morph target
		weight     - new weight for morph target
		
	Returns:
		true if morph target was found, otherwise false
*/
int h3dSetModelMorpher( H3DNode modelNode, const char *target, float weight );


/* Group: Mesh-specific scene graph functions */
/* Function: h3dAddMeshNode
		Adds a Mesh node to the scene.
	
	Details:
		This function creates a new Mesh node and attaches it to the specified parent node.
	
	Parameters:
		parent       - handle to parent node to which the new node will be attached
		name         - name of the node
		materialRes  - material resource used by Mesh node
		batchStart   - first triangle index of mesh in Geometry resource of parent Model node
		batchCount   - number of triangle indices used for drawing mesh
		vertRStart   - first vertex in Geometry resource of parent Model node
		vertREnd     - last vertex in Geometry resource of parent Model node
		
	Returns:
		handle to the created node or 0 in case of failure
*/
H3DNode h3dAddMeshNode( H3DNode parent, const char *name, H3DRes materialRes, 
                            int batchStart, int batchCount, int vertRStart, int vertREnd );


/* Group: Joint-specific scene graph functions */
/* Function: h3dAddJointNode
		Adds a Joint node to the scene.
	
	Details:
		This function creates a new Joint node and attaches it to the specified parent node.
	
	Parameters:
		parent      - handle to parent node to which the new node will be attached
		name        - name of the node
		jointIndex  - index of joint in Geometry resource of parent Model node
		
	Returns:
		handle to the created node or 0 in case of failure
*/
H3DNode h3dAddJointNode( H3DNode parent, const char *name, int jointIndex );


/* Group: Light-specific scene graph functions */
/* Function: h3dAddLightNode
		Adds a Light node to the scene.
	
	Details:
		This function creates a new Light node and attaches it to the specified parent node.
		The direction vector of the untransformed light node is pointing along the the negative
		z-axis. The specified material resource can define uniforms and projective textures.
		Furthermore it can contain a shader for doing lighting calculations if deferred shading
		is used. If no material is required the parameter can be zero. The context names
		define which shader contexts are used when rendering shadow maps or doing light calculations for
		forward rendering configurations.
	
	Parameters:
		parent           - handle to parent node to which the new node will be attached
		name             - name of the node
		materialRes      - material resource for light configuration or 0 if not used
		lightingContext  - name of the shader context used for doing light calculations
		shadowContext    - name of the shader context used for doing shadow map rendering
		
	Returns:
		handle to the created node or 0 in case of failure
*/
H3DNode h3dAddLightNode( H3DNode parent, const char *name, H3DRes materialRes,
                             const char *lightingContext, const char *shadowContext );


/* Group: Camera-specific scene graph functions */
/* Function: h3dAddCameraNode
		Adds a Camera node to the scene.
	
	Details:
		This function creates a new Camera node and attaches it to the specified parent node.
	
	Parameters:
		parent       - handle to parent node to which the new node will be attached
		name         - name of the node
		pipelineRes  - pipeline resource used for rendering
		
	Returns:
		handle to the created node or 0 in case of failure
*/
H3DNode h3dAddCameraNode( H3DNode parent, const char *name, H3DRes pipelineRes );

/* Function: h3dSetupCameraView
		Sets the planes of a camera viewing frustum.
	
	Details:
		This function calculates the view frustum planes of the specified camera node using the specified view
		parameters.
	
	Parameters:
		cameraNode  - handle to the Camera node which will be modified
		fov         - field of view (FOV) angle
		aspect      - aspect ratio
		nearDist    - distance of near clipping plane
		farDist     - distance of far clipping plane 
		
	Returns:
		nothing
*/
void h3dSetupCameraView( H3DNode cameraNode, float fov, float aspect, float nearDist, float farDist );

/* Function: h3dGetCameraProjMat
		Gets the camera projection matrix.
	
	Details:
		This function gets the camera projection matrix used for bringing the geometry to
		screen space and copies it to the specified array.
	
	Parameters:
		cameraNode  - handle to Camera node
		projMat     - pointer to float array with 16 elements
		
	Returns:
		nothing
*/
void h3dGetCameraProjMat( H3DNode cameraNode, float *projMat );


/* Group: Emitter-specific scene graph functions */
/* Function: h3dAddEmitterNode
		Adds a Emitter node to the scene.
	
	Details:
		This function creates a new Emitter node and attaches it to the specified parent node.
	
	Parameters:
		parent             - handle to parent node to which the new node will be attached
		name               - name of the node
		materialRes        - handle to Material resource used for rendering
		particleEffectRes  - handle to ParticleEffect resource used for configuring particle properties
		maxParticleCount   - maximal number of particles living at the same time
		respawnCount       - number of times a single particle is recreated after dying (-1 for infinite)
		
		
	Returns:
		handle to the created node or 0 in case of failure
*/
H3DNode h3dAddEmitterNode( H3DNode parent, const char *name, H3DRes materialRes,
                               H3DRes particleEffectRes, int maxParticleCount, int respawnCount );

/* Function: h3dAdvanceEmitterTime
		Advances the time value of an Emitter node.
	
	Details:
		This function advances the simulation time of a particle system and continues the particle simulation
		with timeDelta being the time elapsed since the last call of this function. The specified
		node must be an Emitter node.
	
	Parameters:
		emitterNode  - handle to the Emitter node which will be modified
		timeDelta    - time delta in seconds
		
	Returns:
		nothing
*/
void h3dAdvanceEmitterTime( H3DNode emitterNode, float timeDelta );

/* Function: h3dHasEmitterFinished
		Checks if an Emitter node is still alive.
	
	Details:
		This function checks if a particle system is still active and has living particles or
		will spawn new particles. The specified node must be an Emitter node. The function can be
		used to check when a not infinitely running emitter for an effect like an explosion can be
		removed from the scene.
	
	Parameters:
		emitterNode  - handle to the Emitter node which is checked
		
	Returns:
		true if Emitter will no more emit any particles, otherwise or in case of failure false
*/
int h3dHasEmitterFinished( H3DNode emitterNode );
