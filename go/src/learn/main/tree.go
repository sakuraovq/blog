package main

import (
	"fmt"
	queue2 "learn/queue"
	"learn/tree"
)

type myTree struct {
	node *tree.Node
}

func (myNode *myTree) postOrder() {
	if myNode == nil || myNode.node == nil {
		return
	}
	left := &myTree{myNode.node.Left}
	right := &myTree{myNode.node.Right}
	// 后续遍历
	left.postOrder()
	right.postOrder()
	myNode.node.Print()
}

func main() {
	treeNode()

	queueSlice()
}

func queueSlice() {

	quere := queue2.Queue{1}
	quere.Push(2)
	quere.Push(3)
	fmt.Println(quere.Pop())
	fmt.Println(quere.Shift())
	fmt.Println(quere.Empty())
	fmt.Println(quere.Shift())
	fmt.Println(quere.Empty())
	fmt.Println(quere.Shift())
}

func treeNode() {
	var root tree.Node

	// 从最底部开始遍历从下往上遍历
	root = tree.Node{Value: 1}
	root.Left = &tree.Node{}
	root.Right = &tree.Node{2, nil, nil}
	root.Right.Left = new(tree.Node)
	root.Right.Right = tree.CreateNode(3)
	root.Left.Left = tree.CreateNode(2)
	root.Left.Left.SetValue(3)
	root.Left.Left.Left = tree.CreateNode(4)
	root.Left.Left.Right = tree.CreateNode(5)
	//root.Tar()
	fmt.Println("my tree")
	root.TarFunc(func(node *tree.Node) {
		node.Print()
	})
	fmt.Println("END")
	//post := myTree{&root}
	//post.postOrder()
}
