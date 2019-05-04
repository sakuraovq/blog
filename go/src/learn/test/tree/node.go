package tree

import "fmt"

type Node struct {
	Value       int
	Left, Right *Node
}

func CreateNode(value int) *Node {
	return &Node{Value: value}
}

func (node *Node) SetValue(value int) {
	node.Value = value
}

func (node *Node) Print() {
	fmt.Println(node)
}

func (node *Node) String() string {
	fmt.Println(node.Value)
	return  ""
}
