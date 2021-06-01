package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var tag string

const HAND_SHAKE_MSG = "I am shake message"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请输入一个客户端标志")
		os.Exit(0)
	}
	tag = os.Args[1]
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 9901}
	dstAddr := &net.UDPAddr{IP: net.ParseIP("10.31.88.137"), Port: 9527}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr) //这里连接server服务
	if err != nil {
		fmt.Println(err)
	}
	if _, err = conn.Write([]byte("hello, I'm new peer: " + tag)); err != nil { //给server发消息
		log.Panic(err)
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data) //获取server的返回信息  data是另一台电脑的ip地址  remoteAddr是server的地址
	if err != nil {
		fmt.Println("error during read ", err)
	}
	conn.Close()
	anotherPeer := parseAddr(string(data[:n]))
	fmt.Printf("local: %s server: %s another: %s \n", srcAddr, remoteAddr, anotherPeer.String())
	//开始打洞
	bidirectionHole(srcAddr, &anotherPeer)
}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func bidirectionHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	//向另一个peer发送一条udp消息（对方peer的nat设备会丢弃该消息，非法来源），用意在自身的nat设备打开一条可进入的通道，这样对方就可以发过来udp消息
	if _, err := conn.Write([]byte(HAND_SHAKE_MSG)); err != nil {
		log.Println("send handshake:", err)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from [" + tag + "]")); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()
	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf("收到数据： %s\n", data[:n])
		}
	}
}
