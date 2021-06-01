package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8008")
	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return
	}
	defer conn.Close()
	go func() { //读取服务器返回的数据
		msg := make([]byte, 4096)
		for {
			n, err := conn.Read(msg)
			if n == 0 || err == io.EOF {
				fmt.Println("退出")
				os.Exit(1)
			}
			if err != nil {
				fmt.Println("read server message error: ", err)
				continue
			}
			fmt.Println(string(msg[:n]))
		}
	}()
	//获取用户写入数据
	buf := make([]byte, 4096)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println("user input error")
			continue
		}
		conn.Write(buf[:n])
	}
}
