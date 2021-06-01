package main

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {
	head := &TreeNode{Data: "A"}
	head.Left = &TreeNode{Data: "B"}
	head.Right = &TreeNode{Data: "C"}
	head.Left.Left = &TreeNode{Data: "D"}
	head.Left.Right = &TreeNode{Data: "E"}
	head.Right.Left = &TreeNode{Data: "F"}
	head.Right.Right = &TreeNode{Data: "G"}
	fmt.Println("先序排序： ")
	PreOrder(head)
	fmt.Println("\n中序排序： ")
	MidOrder(head)
	fmt.Println("\n后序排序： ")
	PostOrder(head)
	fmt.Println("\n层次遍历： ")
	LayerOrder(head)
	t.Log("success")
}
