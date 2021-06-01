package main

import (
	"fmt"
	"time"
)

func algo(arrs [][]string) {
	arr1 := arrs[0]
	result := make([]string, 0)

	for _, data := range arr1 {
		result = append(result, data)
	}

	for i := 1; i < len(arrs); i++ {
		temp := haha(result, arrs[i])
		result = temp
	}
	fmt.Println(result)
}

func haha(s []string, arr []string) []string {
	temp := make([]string, 0)
	for _, r := range s {
		for _, data := range arr {
			fmt.Println(r + data)
			temp = append(temp, r+data)
		}
	}
	return temp
}

var (
	a = 0
	b = 1
)

func fibonacci() {
	for b < 10 {
		println(b)
		a, b = b, a+b
	}
}

func reverse(s []rune) []rune {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func main() {
	/*algo([][]string{
		{"1","2","3"},
		{"4","5"},
		{"6"},
		{"7","8"},
	})*/
	//fibonacci()
	mm := NewMyMay()

	for i := 0; i < 10; i++ {
		go mm.Store(fmt.Sprintf("i-%d", i), i)
	}

	time.Sleep(20)
	for i := 0; i < 10; i++ {
		go func(i int) {
			r := mm.Get(fmt.Sprintf("i-%d", i))
			fmt.Println(r)
		}(i)
	}
	fmt.Println(mm)
}

type MyMap struct {
	Data map[string]interface{}
	ch   chan func()
}

func NewMyMay() *MyMap {
	m := &MyMap{
		Data: make(map[string]interface{}),
		ch:   make(chan func()),
	}
	go func() {
		for {
			(<-m.ch)()
		}
	}()
	return m
}

func (m *MyMap) Store(k string, v interface{}) {
	m.ch <- func() {
		m.Data[k] = v
	}
}

func (m *MyMap) Get(k string) (data interface{}) {
	m.ch <- func() {
		if res, ok := m.Data[k]; ok {
			data = res
		}
	}
	return
}

func (m *MyMap) Delete(k string) {
	m.ch <- func() {
		delete(m.Data, k)
	}
}
