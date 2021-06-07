package main

import "fmt"

func test1() (i int) { //如果变量生命在返回值中则是所有运算的值
	i = 0
	defer func() {
		i = i + 1
	}()
	return i
}

func test2() int { //如果变量不定义在返回值中则跟defer无关
	i := 0
	defer func() {
		i = i + 1
	}()
	return i
}

func main() {
	fmt.Println(test1()) // 1
	fmt.Println(test2()) // 0
}
