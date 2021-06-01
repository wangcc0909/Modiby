package algo

import (
	"fmt"
	"sync"
)

type TreeNode struct {
	Value string
	Left  *TreeNode
	Right *TreeNode
}

//层次遍历 BFS
type LinkNode struct {
	Next  *LinkNode
	Value *TreeNode
}

type LinkQueue struct {
	Root *LinkNode
	Size int
	lock sync.Mutex
}

func (q *LinkQueue) Add(tree *TreeNode) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.Root == nil {
		q.Root = new(LinkNode)
		q.Root.Value = tree
	} else {
		newNode := new(LinkNode)
		newNode.Value = tree
		nextNode := q.Root
		for nextNode.Next != nil {
			nextNode = nextNode.Next
		}
		nextNode.Next = newNode
	}
	q.Size = q.Size + 1
}

func (q *LinkQueue) Remove() *TreeNode {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.Size == 0 {
		panic("over limit")
	}
	node := q.Root
	v := node.Value
	q.Root = node.Next
	q.Size = q.Size - 1
	return v
}

func LayerOrder(tree *TreeNode) {
	if tree == nil {
		return
	}
	queue := new(LinkQueue)
	queue.Add(tree)
	for queue.Size > 0 {
		node := queue.Remove()
		fmt.Print(node.Value, " ")
		if node.Left != nil {
			queue.Add(node.Left)
		}
		if node.Right != nil {
			queue.Add(node.Right)
		}
	}
}
