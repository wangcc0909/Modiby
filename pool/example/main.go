package main

import (
	"fmt"
	"github.peaut.limit/pool"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const addr string = "127.0.0.1:8080"

func main() {
	go server()
	time.Sleep(2 * time.Second)
	client()
	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	select {
	case s := <-sg:
		log.Println("退出： ", s.String())
	}
}

func client() {
	factory := func() (interface{}, error) { return net.Dial("tcp", addr) }
	close := func(v interface{}) error { return v.(net.Conn).Close() }
	poolConfig := &pool.Config{
		InitialCap:  2,
		MaxCap:      5,
		MaxIdle:     4,
		Factory:     factory,
		Close:       close,
		IdleTimeout: 15 * time.Second,
	}
	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		fmt.Println("err = ", err)
	}
	v, err := p.Get()
	//do something
	//conn := v.(net.Conn)
	p.Put(v)
	current := p.Len()
	fmt.Println("len = ", current)
}

func server() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error accepting: ", err)
		}
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
	}
}
