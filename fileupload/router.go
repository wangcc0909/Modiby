package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Router(app *gin.Engine) {
	app.GET("/checkChunk", func(c *gin.Context) {
		hash := c.Query("hash") //前端通过get传过来的文件hash值
		hashPath := fmt.Sprintf("./uploadFile/%s", hash)
		var chunkList []string
		//检查路径是否存在
		isExistPath, err := PathExists(hashPath)
		if err != nil {
			fmt.Println("获取hash路径错误", err)
		}
		if isExistPath {
			//如果文件夹存在开始读取所有文件并且返回已经存在的文件名
			files, err := ioutil.ReadDir(hashPath)
			// state状态为0说明文件不完整，1为文件已经完全上传完毕并且已经存在文件了
			state := 0
			if err != nil {
				fmt.Println("文件读取错误", err)
			}
			for _, f := range files {
				fileName := f.Name()
				chunkList = append(chunkList, fileName)
				//如果存在一个文件名与hash值一致的文件，说明已经生成了完整的文件，就不需要进行上传了
				fileBaseName := strings.Split(fileName, ".")[0]
				if fileBaseName == hash {
					state = 1
				}
			}
			c.JSON(http.StatusOK, gin.H{
				"state":     state,
				"chunkList": chunkList,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"state":     0,
				"chunkList": chunkList,
			})
		}
	})
	app.POST("/uploadChunk", func(c *gin.Context) {
		fileHash := c.PostForm("hash")                       //文件hash
		file, err := c.FormFile("file")                      //具体文件
		hashPath := fmt.Sprintf("./uploadFile/%s", fileHash) //上传的文件夹路径
		if err != nil {
			fmt.Println("获取上传文件失败", err)
		}
		isExistPath, err := PathExists(hashPath)
		if err != nil {
			fmt.Println("获取hash路径错误", err)
		}
		//如果不存在则进行文件夹创建
		if !isExistPath {
			os.Mkdir(hashPath, os.ModePerm)
		}
		// gin提供的方法，把上传的文件保存到提供的路径
		err = c.SaveUploadedFile(file, fmt.Sprintf("./uploadFile/%s/%s", fileHash, file.Filename))
		if err != nil {
			c.String(400, "0")
			fmt.Println(err)
		} else {
			chunkList := []string{}
			//读取文件夹下的所有文件(得给前端一个反馈，告诉前端我们已经拿到了多少哥区块了)
			files, err := ioutil.ReadDir(hashPath)
			if err != nil {
				fmt.Println("文件读取错误", err)
			}
			for _, f := range files {
				fileName := f.Name()
				//因为本地开发是Mac 所有存在这个文件，进行排除一下
				if f.Name() == ".DS_Store" {
					continue
				}
				chunkList = append(chunkList, fileName)
			}
			c.JSON(http.StatusOK, gin.H{
				"chunkList": chunkList,
			})
		}
	})
	app.GET("mergeChunk", func(c *gin.Context) {
		hash := c.Query("hash")
		fileName := c.Query("fileName")
		hashPath := fmt.Sprintf("./uploadFile/%s", hash)
		isExistPath, err := PathExists(hashPath)
		if err != nil {
			fmt.Println("获取hash路径错误", err)
		}
		if !isExistPath {
			c.JSON(400, gin.H{
				"message": "文件夹不存在",
			})
			return
		}
		isExistFile, err := PathExists(hashPath + "/" + fileName)
		if err != nil {
			fmt.Println("获取hash路径文件错误", err)
		}
		//如果文件已经存在了我们直接返回不进行区块合并了
		if isExistFile {
			c.JSON(http.StatusOK, gin.H{
				"fileUrl": fmt.Sprintf("http://127.0.0.1:9999/uploadFile/%s/%s", hash, fileName),
			})
			return
		}
		files, err := ioutil.ReadDir(hashPath)
		if err != nil {
			fmt.Println("合并问价读取失败", err)
		}
		//创建文件
		complateFile, err := os.Create(hashPath + "/" + fileName)
		//关闭文件
		defer complateFile.Close()
		for _, f := range files {
			if f.Name() == ".DS_Store" {
				continue
			}
			//读取区块文件
			fileBuffer, err := ioutil.ReadFile(hashPath + "/" + f.Name())
			if err != nil {
				fmt.Println("文件打开错误", err)
			}
			// 写入刚刚创建的文件
			complateFile.Write(fileBuffer)
		}
		c.JSON(http.StatusOK, gin.H{
			"fileUrl": fmt.Sprintf("http://127.0.0.1:9999/uploadFile/%s/%s", hash, fileName),
		})
	})
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
