package main

import (
	"container/ring"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	limitCount              = 10
	limitBucket             = 6
	curCount     int32      = 0
	head         *ring.Ring //环形队列（链表)
	printVersion bool
)

func init() {
	flag.BoolVar(&printVersion, "version", false, "print program build version")
	flag.Parse()
}

func main() {
	if printVersion {
		PrintVersion()
	}

	addr, err := net.ResolveTCPAddr("tcp4", "0.0.0.0:9090")
	checkError(err)
	listen, err := net.ListenTCP("tcp", addr)
	checkError(err)
	defer func() {
		listen.Close()
	}()
	head = ring.New(limitBucket)
	for i := 0; i < limitBucket; i++ {
		head.Value = 0
		head = head.Next()
	}
	go func() {
		timer := time.NewTicker(time.Second * 1)
		for range timer.C {
			subCount := int32(0 - head.Value.(int))
			newCount := atomic.AddInt32(&curCount, subCount)
			arr := [6]int{}
			for i := 0; i < limitBucket; i++ {
				arr[i] = head.Value.(int)
				head = head.Next()
			}
			fmt.Println("move subCount, newCount,arr", subCount, newCount, arr)
			head.Value = 0
			head = head.Next()
		}
	}()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handle(&conn)
	}
}

func handle(conn *net.Conn) {
	defer (*conn).Close()
	n := atomic.AddInt32(&curCount, 1)
	if n > int32(limitCount) {
		atomic.AddInt32(&curCount, -1)
		(*conn).Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\nError, too many request, please try again."))
	} else {
		mu := sync.Mutex{}
		mu.Lock()
		pos := head.Prev()
		val := pos.Value.(int)
		val++
		pos.Value = val
		mu.Unlock()
		time.Sleep(time.Second * 1)
		(*conn).Write([]byte("HTTP/1.1 200 OK\r\n\r\nI can change the world!"))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
