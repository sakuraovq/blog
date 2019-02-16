package main

import "fmt"

type Node struct {
	value       int
	left, right *Node
}

func createNode(value int) *Node {
	return &Node{value: value}
}

func (node *Node) SetValue(value int) {
	node.value = value
}

func (node *Node) print() {
	fmt.Println(node)
}

func (node *Node) tar() {
	if node == nil {
		return
	}
	// 先遍历左边(每次吧节点的自己传递进来了)
	node.left.tar()
	// 在显示中间
	node.print()
	// 最后是右边
	node.right.tar()
}

func main() {
	var root Node

	root = Node{value: 1}
	root.left = &Node{}
	root.right = &Node{2, nil, nil}
	root.right.left = new(Node)
	root.right.right = createNode(3)
	root.left.right = createNode(2)

	root.left.right.SetValue(3)


	//root.left.right.print()

	root.tar()
}

func batch() {
	nodes := []Node{
		{value: 3},
		{},
		{5, nil, nil},
	}
	fmt.Println(nodes)
}
