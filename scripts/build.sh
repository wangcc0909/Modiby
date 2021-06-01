#!/bin/bash
cd ..

VERSION=1.0.1
path="github.peaut.limit"

GoVersion=$(go version)
GitCommit=$(git rev-parse --short HEAD)
echo $GoVersion
echo $GitCommit

go build -ldflags "-X main.Version=$VERSION -X main.GoVersion=$GoVersion -X main.BuildTime = `date +"%Y-%m-%d %H:%m:%s"` -X main.GitCommit=$GitCommit" -o limit