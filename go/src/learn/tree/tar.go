package tree

func (node *Node) Tar() {
	if node == nil {
		return
	}
	// 先遍历左边(每次吧节点的自己传递进来了)
	node.Left.Tar()
	// 在显示中间
	node.Print()
	// 最后是右边
	node.Right.Tar()
}

func (node *Node) TarFunc(f func(*Node)) {
	if node == nil {
		return
	}
	node.Left.TarFunc(f)
	f(node)
	node.Right.TarFunc(f)
}
