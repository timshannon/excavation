package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"errors"
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
func (n *Node) Transform() (tx, ty, tz, rx, ry, rz, sx, sy, sz float32) {
	horde3d.GetNodeTransform(n.H3DNode, &tx, &ty, &tz, &rx, &ry, &rz, &sx, &sy, &sz)
	return
}

//This function sets the relative translation, rotation and scale of a 
//specified scene node object.  The coordinates are in local space and 
//contain the transformation of the node relative to its parent.
func (n *Node) SetTransform(tx, ty, tz, rx, ry, rz, sx, sy, sz float32) {
	horde3d.SetNodeTransform(n.H3DNode, tx, ty, tz, rx, ry, rz, sx, sy, sz)
}

//TODO: use Matrix type?  goMatrix?
//Gets the relative transformation matrix of the node
func (n *Node) RelativeTransMat() []float32 {
	relative := make([]float32, 16)
	horde3d.GetNodeTransMats(n.H3DNode, relative, nil)
	return relative
}

//Gets the absolute transformation matrix of the node
func (n *Node) AbsoluteTransMat() []float32 {
	absolute := make([]float32, 16)
	horde3d.GetNodeTransMats(n.H3DNode, nil, absolute)
	return absolute
}

//Sets the relative transformation matrix of the node
func (n *Node) SetRelativeTransMat(matrix []float32) {
	horde3d.SetNodeTransMat(n.H3DNode, matrix)
}

//Returns the bounds of a box that encompasses the entire scene node
func (n *Node) BoundingBox() (minX, minY, minZ, maxX, maxY, maxZ float32) {
	horde3d.GetNodeAABB(n.H3DNode, &minX, &minY, &minZ, &maxX, &maxY, &maxZ)
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

type Group struct{ *Node }
type Model struct{ *Node }
type Mesh struct{ *Node }
type Joint struct{ *Node }
type Light struct{ *Node }
type Camera struct{ *Node }
type Emitter struct{ *Node }
