package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

//用户上线
//用户广播消息
//查询在线用户
//用户下线
//超时退出

type Client struct {
	C    chan string
	Name string
	Addr string
}

var userMap map[string]Client

var message = make(chan string)

func makeMsg(clt Client, msg string) (buf string) {
	buf = "[" + clt.Addr + "]" + clt.Name + ": " + msg
	return
}

func handleConnect(conn net.Conn) {
	defer conn.Close()
	cltAddr := conn.RemoteAddr().String()
	clt := Client{
		C:    make(chan string),
		Name: cltAddr,
		Addr: cltAddr,
	}
	userMap[cltAddr] = clt
	go func() {
		for msg := range clt.C {
			conn.Write([]byte(msg))
		}
	}()
	message <- makeMsg(clt, "login")

	buf := make([]byte, 4096)
	isQuit := make(chan struct{})
	isActive := make(chan struct{})
	go func() {
		for {
			n, err := conn.Read(buf)
			if err == io.EOF || n == 0 {
				isQuit <- struct{}{}
				return
			}
			if err != nil {
				fmt.Println("conn.Read error:", err)
				continue
			}
			fmt.Println("收到", cltAddr, "的消息：", string(buf[:n]))
			msg := string(buf[:n-1])
			if string(msg) == "who" {
				conn.Write([]byte("online user list: \n"))
				for addr, user := range userMap {
					conn.Write([]byte(addr + ":" + user.Name + "\n"))
				}
			} else if len(msg) > 8 && string(msg[:6]) == "rename" {
				newName := strings.Split(msg, "|")[1]
				clt.Name = newName
				userMap[cltAddr] = clt
				conn.Write([]byte("rename successful"))
			} else {
				message <- makeMsg(clt, string(msg))
			}
			isActive <- struct{}{}
		}
	}()

	for {
		select {
		case <-isQuit:
			close(clt.C)
			delete(userMap, cltAddr)
			message <- makeMsg(clt, "logout")
			return
		case <-isActive:

		case <-time.After(100 * time.Second):
			delete(userMap, cltAddr)
			message <- makeMsg(clt, "logout")
			return
		}
	}
}

//管理全局map和message
func Message() {
	for {
		msg := <-message
		for _, client := range userMap {
			client.C <- msg
		}
	}
}

func init() {
	userMap = make(map[string]Client)
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8008")
	if err != nil {
		fmt.Println("net.Listen error: ", err)
		return
	}
	defer listener.Close()
	go Message()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener error:", err)
			continue
		}
		go handleConnect(conn)
	}
}
