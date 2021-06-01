package balance

import (
	"fmt"
	"testing"
)

func TestRandomBalance(t *testing.T) {
	rlb := LoadBalanceFactory(Random)
	rlb.Add("a")
	rlb.Add("b")
	rlb.Add("c")
	rlb.Add("d")
	for i := 0; i < 20; i++ {
		rs, _ := rlb.Get()
		fmt.Println(rs)
	}
}

func TestRoundRobinBalance(t *testing.T) {
	rlb := LoadBalanceFactory(RoundRobin)
	rlb.Add("a")
	rlb.Add("b")
	rlb.Add("c")
	rlb.Add("d")
	for i := 0; i < 20; i++ {
		rs, _ := rlb.Get()
		fmt.Println(rs)
	}
}

func TestWeightBalance(t *testing.T) {
	rlb := LoadBalanceFactory(Weight)
	rlb.Add("a", "1")
	rlb.Add("b", "2")
	rlb.Add("c", "5")
	var count = make(map[string]int)
	for i := 0; i < 200000; i++ {
		rs, _ := rlb.Get()
		count[rs]++
	}
	fmt.Println(count)
}
