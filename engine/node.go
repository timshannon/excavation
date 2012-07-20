package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"errors"
)

type Node struct {
	horde3d.H3DNode
}

func AddNodes(parent Node, sceneResource Resource) *Node, error {
	node := new(Node)
	node.H3DNode = horde3d.AddNodes(parent.H3DNode, sceneResource.H3DRes)

	return node
}

//This function returns the type of a specified scene node.  If the node handle is invalid, 
//the function returns the node type Unknown.
func (n *Node) Type() int { return horde3d.GetNodeType(n.H3DNode) }

func (n *Node) Parent() *Node {
	parent := new(Node)
	parent.H3DNode = horde3d.GetNodeParent(n.H3DNode)
	return parent
}

func (n *Node) SetParent(parent Node) bool { return horde3d.SetNodeParent(n.H3DNode, parent.H3DNode) }
