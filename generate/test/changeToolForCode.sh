#!/usr/bin/env gosl

import "fmt"
import "strings"
import "bufio"
import "io/ioutil"
import "os"

if len(os.Args) != 4 {
    fmt.Println("请检查输入的读取文件写入文件路径，以及所需要的")
    fmt.Println("..changeTool <readPath> <writePath> <name>")
}

var strs []string

//命令传入读取文件
file_bytes, err := ioutil.ReadFile(os.Args[1])
if err != nil {
    fmt.Println(err.Error() + \n")
    return
}

//按行读取
lines := strings.Split(string(file_bytes),"\n")

//根据 = 进行切割， 判断 = 右边有没有//字符串， 之后进行切割组合
for _, line := range lines {
    items := strings.Split(string(line), "=")
    for _, item := range items {
        if strings.Contains(item," //") {
            items := strings.Split(string(item), " //")
            items[0] = strings.Replace(items[0]," ","", -1)
            items[1] = strings.Replace(items[1]," ","", -1)

            str := "    " + items[0] + ": " + "\"" + items[1] + "\""
            strs = append(strs,str)
        }
    }
}

// 打开命令行传入的文件， 进行追加内容
f,err := os.OpenFile(os.Args[2],os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
if err != nil {
    fmt.Println(err.Error() + "\n")
    return
}

defer f.Close()

w := bufio.NewWrite(f)

fmt.Fprintln(w,"\n"+os.Args[3]+":")
for _, str := range strs {
    fmt.Fprintln(w, str)
}

w.Flush()

fmt.Println("success")