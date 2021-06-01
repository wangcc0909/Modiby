package main

import (
	"fmt"
	"sync"
)

type order interface {
	Add(node *TreeNode)
	Remove() *TreeNode
}

type TreeNode struct {
	Data  string
	Left  *TreeNode
	Right *TreeNode
}

//先序遍历
func PreOrder(tree *TreeNode) {
	if tree == nil {
		return
	}
	fmt.Print(tree.Data, " ")
	PreOrder(tree.Left)
	PreOrder(tree.Right)
}

//中序遍历
func MidOrder(tree *TreeNode) {
	if tree == nil {
		return
	}
	MidOrder(tree.Left)
	fmt.Print(tree.Data, " ")
	MidOrder(tree.Right)
}

//后序遍历
func PostOrder(tree *TreeNode) {
	if tree == nil {
		return
	}
	PostOrder(tree.Left)
	PostOrder(tree.Right)
	fmt.Print(tree.Data, " ")
}

//层次遍历
type LinkNode struct {
	Next  *LinkNode
	Value *TreeNode
}

type LinkQueue struct {
	root *LinkNode
	size int
	lock sync.Mutex
}

//入队
func (q *LinkQueue) Add(tree *TreeNode) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.root == nil {
		q.root = new(LinkNode)
		q.root.Value = tree
	} else {
		newNode := new(LinkNode)
		newNode.Value = tree

		nextNode := q.root
		for nextNode.Next != nil {
			nextNode = nextNode.Next
		}
		nextNode.Next = newNode
	}
	q.size += 1
}

//出队
func (q *LinkQueue) Remove() *TreeNode {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.size == 0 {
		panic("over limit")
	}
	topNode := q.root
	v := topNode.Value
	q.root = topNode.Next
	q.size = q.size - 1
	return v
}

func LayerOrder(tree *TreeNode) {
	if tree == nil {
		return
	}
	linkQueue := new(LinkQueue)
	linkQueue.Add(tree)
	for linkQueue.size > 0 {
		element := linkQueue.Remove()
		fmt.Print(element.Data, " ")
		if element.Left != nil {
			linkQueue.Add(element.Left)
		}
		if element.Right != nil {
			linkQueue.Add(element.Right)
		}
	}
}
