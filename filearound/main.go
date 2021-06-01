package main

import (
	"github.com/sirupsen/logrus"
	"github.peaut.limit/filearound/web"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go func() {
		web.Run()
	}()
	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	select {
	case s := <-sg:
		logrus.Infof("got signal: %s", s.String())
	}
}
