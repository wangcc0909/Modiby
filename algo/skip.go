package algo

import (
	"math/rand"
	"sync"
)

type SkipListNode struct { //跳跃表节点定义
	key   int
	value interface{}     //key value为跳跃表中每个节点的键值对
	next  []*SkipListNode //指向多层连标节点的数组
}

type SkipList struct {
	head, tail    *SkipListNode //跳跃表的起始节点指针地址
	length, level int           //跳跃遍长度和层数
	mut           sync.RWMutex  //用于集合并发访问
	rand          *rand.Rand    //内部保留的随机数
}

const P uint32 = 4

func (list *SkipList) random() int {
	//当新增节点时随机生成层数，定义一个平衡P叉树
	//有多种实现算法，redis及leveldb中一般采用P=4的平衡四叉树
	level := 1
	for level < list.length && ((list.rand.Uint32() % P) == 0) {
		level++
	}
	if level < list.level {
		return level
	} else {
		return list.level
	}
}

func (list SkipList) AddNode(key int, value interface{}) {
	list.mut.Lock()
	defer list.mut.Unlock()
	//生成随机数层
	level := list.random()
	//定位新元素的插入点
	update := make([]*SkipListNode, level)
	node := list.head
	//从高到低逐层进行定位（这里默认0是最底层）
	for index := level - 1; index >= 0; index-- {
		for {
			nextNode := node.next[index]
			//n最开始指向最高层的head的next节点，从head开始逐个比较
			if nextNode == list.tail || nextNode.key > key {
				update[index] = node
				break
			} else if nextNode.key == key {
				//此处简化算法 如果key值相同则覆盖 即保留唯一的key值节点
				nextNode.value = value
				return
			} else {
				//如未达到队尾且新增元素key值大于当前节点key值，则继续遍历列表
				node = nextNode
			}
		}
	}
	//生成并初始化新节点
	newNode := &SkipListNode{key, value, make([]*SkipListNode, level)}
	for index, node := range update {
		node.next[index], newNode.next[index] = newNode, node.next[index]
	}
	list.length++
}
