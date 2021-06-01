package main

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
)

const stdErrFile = "./tmp/go-stderr.log"

var stdErrFileHandler *os.File

func RewriteStderrFile() error {
	if runtime.GOOS == "windows" {
		return nil
	}

	file, err := os.OpenFile(stdErrFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	stdErrFileHandler = file

	if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		fmt.Println(err)
		return err
	}
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})
	return nil
}

func testPanic() {
	panic("test panic")
}

func main() {
	RewriteStderrFile()
	testPanic()
}
