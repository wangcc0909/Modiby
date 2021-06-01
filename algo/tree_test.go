package algo

import "testing"

func TestTree(t *testing.T) {
	head := &TreeNode{Value: "A"}
	head.Left = &TreeNode{Value: "B"}
	head.Right = &TreeNode{Value: "C"}
	head.Left.Left = &TreeNode{Value: "D"}
	head.Left.Right = &TreeNode{Value: "E"}
	head.Right.Left = &TreeNode{Value: "F"}
	head.Right.Right = &TreeNode{Value: "G"}
	LayerOrder(head)
	t.Log("end... ")
}
