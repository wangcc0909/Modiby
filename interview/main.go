package main

import (
	"fmt"
	"math/big"
	"net"
)

type Item struct {
	Value int
	Next  *Item
}

func main() {
	/*item1 := &Item{
		Value: 1,
		Next: &Item{
			Value: 3,
			Next: &Item{
				Value: 5,
			},
		},
	}
	item2 := &Item{
		Value: 2,
		Next: &Item{
			Value: 4,

		},
	}
	result := mergeList(item1,item2)
	PrintList(result)*/
	ip := "192.168.1.53"
	ipInt := net_aton(ip)
	fmt.Println("ip int:", ipInt)
	fmt.Println("ip string:", net_ntoa(ipInt))
}

//如果list1的value比list2的value小就先取list1 如果有一方next为nil则直接将
func mergeList(list1, list2 *Item) *Item {
	if list1 == nil {
		return list2
	}
	if list2 == nil {
		return list1
	}

	var result *Item
	if list1.Value < list2.Value {
		result = list1
		result.Next = mergeList(list1.Next, list2)
	} else {
		result = list2
		result.Next = mergeList(list1, list2.Next)
	}
	return result
}

func PrintList(item *Item) {
	head := item
	fmt.Print(head.Value)
	for head.Next != nil {
		head = head.Next
		fmt.Print(head.Value)
	}
}

func net_aton(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func net_ntoa(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
