package service

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pborman/uuid"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetFileList(c *fiber.Ctx) error {
	var files []string
	err := filepath.Walk("./filearound/files", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, info.Name())
		return nil
	})
	if err != nil {
		return err
	}
	return c.JSON(files)
}

func Download(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return errors.New("file not found")
	}
	fp := fmt.Sprintf("./filearound/files/%s", filename)
	return c.SendFile(fp)
}

func UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	//c.SaveFile(file,fmt.Sprintf("./%s",file.Filename))  //保存文件
	ext := filepath.Ext(file.Filename)
	if ext != ".zip" {
		c.WriteString("文件格式不对")
		return nil
	}
	savePath := fmt.Sprintf("./tmp/%s", file.Filename)
	err = c.SaveFile(file, savePath)
	if err != nil {
		return err
	}
	t := time.Now()
	unzipPath := "./" + t.Format("20060102150405")
	err = unzip(savePath, unzipPath)
	if err != nil {
		return err
	}
	//重命名(相当于排序)
	err = fileSort(unzipPath)
	if err != nil {
		return err
	}
	//压缩文件
	t = time.Now()
	zipFilePath := "./filearound/files/" + t.Format("20060102150405") + ".zip"
	err = Zip(unzipPath, zipFilePath)
	if err != nil {
		return err
	}
	return os.RemoveAll(unzipPath)
}

func fileSort(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			newName := filepath.Dir(path) + "/" + uuid.New() + filepath.Ext(path)
			os.Rename(path, newName)
		}
		return nil
	})
}

//解压 zip
func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}
	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if strings.Contains(path, "MACOSX") {
			continue
		}
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}

func Zip(srcFile string, destZip string) error {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	archive := zip.NewWriter(zipFile)
	defer archive.Close()
	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == "./filearound" || path == "./filearound/tmp" {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})
	return err
}
