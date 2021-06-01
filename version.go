package main

import (
	"fmt"
)

var (
	//Version 项目版本信息
	Version = ""
	//GoVersion GO版本信息
	GoVersion = ""
	//GitCommit git 提交commit id
	GitCommit = ""
	//BuildTime 构建时间
	BuildTime = ""
)

//输出版本信息
func PrintVersion() {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Go Version: %s\n", GoVersion)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Build Time: %s\n", BuildTime)
}
