package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"github.com/saracen/go7z"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

const (
	TypeFileZip   = ".zip"
	TypeFileTarGz = ".tar.gz"
	TypeFileTgz   = ".tgz"
	TypeFile7z    = ".7z"
)

type TemplateContens struct {
	Name     string `json:"Name"`
	FileType string `json:"file_type"`
}

var (
	uploadFileKey = "uploadfile"
)
var dataMaps []map[string]interface{}

type CiHandler struct {
}

// NewCiHandler : New Ci Handler
func NewCiHandler() *CiHandler {
	h := CiHandler{}
	return &h
}

func main() {
	fmt.Println(os.Getwd())
	dir := "./testdata"
	fileName := "/test.7z"
	var err error
	if strings.HasSuffix(fileName, ".zip") {
		err = unzip(dir+fileName, dir)
	} else if strings.HasSuffix(fileName, "tar.gz") || strings.HasSuffix(fileName, "tgz") {
		err = UnTarGz(dir+fileName, dir)
	} else if strings.HasSuffix(fileName, ".7z") || strings.HasSuffix(fileName, ".rar") {
		err = un7z(dir+fileName, dir)
	}
	fmt.Println(err)
	//var s []string
	//s,_ = getFileList(dir,s)
	//fmt.Println(s)
	//CleanSlice()
	//http.HandleFunc("/upload", UploadFile)
	//http.ListenAndServe(":8088", nil)

	/*ci := NewCiHandler()
	value := reflect.ValueOf(ci)
	fmt.Println("value = ",value)
	typ := value.Type()
	//这里就是struct名字
	fmt.Println("type = ",typ.String()[len("*main."):])*/
	/*err := os.RemoveAll("./unarr/tmp1598951825")
	if err != nil {
		fmt.Println(err.Error())
	}*/
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile(uploadFileKey)
	if err != nil {
		//ignore the error handler
		fmt.Println(err.Error())
	}
	defer file.Close()
	dst, err := os.Create(header.Filename)
	if err != nil {
		w.Write([]byte("request error"))
		return
	}
	_, err = io.Copy(dst, file)
	if err != nil {
		w.Write([]byte("copy error"))
		return
	}
	fileName := header.Filename
	dir := "./tmp" + fmt.Sprint(time.Now().Unix())

	err = os.Mkdir(dir, os.ModePerm)
	if strings.HasSuffix(fileName, TypeFileZip) {
		err = unzip(fileName, dir)
	} else if strings.HasSuffix(fileName, TypeFileTarGz) || strings.HasSuffix(fileName, TypeFileTgz) {
		err = UnTarGz(fileName, dir)
	} else if strings.HasSuffix(fileName, TypeFile7z) {
		err = un7z(fileName, dir)
	}
	if err != nil {
		w.Write([]byte("archive error"))
		return
	}
	dataMaps = dataMaps[0:0]
	getFileList(dir, "")
	fmt.Println(dataMaps)
	w.Write([]byte("success,good job"))
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

//解压 tar.gz 和 tgz
func UnTarGz(srcGz string, dstSrc string) error {
	dstDir := path.Clean(dstSrc) + string(os.PathSeparator)
	fr, err := os.Open(srcGz)
	if err != nil {
		return err
	}
	defer fr.Close()

	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	if gr == nil {
		return fmt.Errorf("call gzip.NewReader return nil")
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		fi := hdr.FileInfo()
		dstFullPath := dstDir + hdr.Name
		if hdr.Typeflag == tar.TypeDir {
			_ = os.MkdirAll(dstFullPath, fi.Mode().Perm())
			_ = os.Chmod(dstFullPath, fi.Mode().Perm())
		} else {
			tempDir := dstFullPath[:strings.LastIndex(dstFullPath, "\\")+1]
			_ = os.MkdirAll(tempDir, os.ModePerm)
			if err := unTarFile(dstFullPath, tr); err != nil {
				fmt.Println(err.Error())
				return err
			}
			_ = os.Chmod(dstFullPath, fi.Mode().Perm())
		}
	}
	return nil
}

func unTarFile(dstFile string, tr *tar.Reader) error {
	fw, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer fw.Close()
	_, err = io.Copy(fw, tr)
	if err != nil {
		return err
	}
	return nil
}

func un7z(src string, dst string) error {
	/*a, err := unarr.NewArchive(src)
	if err != nil {
		return err
	}
	defer a.Close()
	_, err = a.Extract(dst)
	return err*/
	sz, err := go7z.OpenReader(src)
	if err != nil {
		panic(err)
	}
	defer sz.Close()

	for {
		hdr, err := sz.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			panic(err)
		}

		// If empty stream (no contents) and isn't specifically an empty file...
		// then it's a directory.
		if hdr.IsEmptyStream && !hdr.IsEmptyFile {
			os.MkdirAll(dst+"/"+hdr.Name, os.ModePerm)
			continue
		}

		// Create file
		f, err := os.Create(dst + "/" + hdr.Name)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := io.Copy(f, sz); err != nil {
			panic(err)
		}
	}
	return nil
}

func CleanSlice() {
	var cSlice = []int{1, 2, 3}
	fmt.Println("清空前元素>>: \n", cSlice)
	cSlice = cSlice[0:0]
	fmt.Println("清空后元素>>: \n", cSlice)
	cSlice = append(cSlice, 5)
	fmt.Println("清空后添加元素>>: \n", cSlice)
}

func getFileList(filePath string, tmpeType string) {
	fs, _ := ioutil.ReadDir(filePath)
	for _, file := range fs {
		if file.IsDir() {
			fileName := file.Name()
			if strings.EqualFold(fileName, "dir1") || strings.EqualFold(fileName, "Outgress GW") || strings.EqualFold(fileName, "Common") {
				getFileList(filePath+"/"+file.Name()+"/", fileName)
			} else {
				getFileList(filePath+"/"+file.Name()+"/", tmpeType)
			}
		} else {
			fileName := file.Name()
			tc := TemplateContens{
				Name:     fileName,
				FileType: tmpeType,
			}
			dataMap := Struct2Map(tc)
			dataMaps = append(dataMaps, dataMap)
		}
	}
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

func getAllFile(dir string, files []TemplateContens) ([]TemplateContens, error) {

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, TemplateContens{
				Name:     info.Name(),
				FileType: info.Name(),
			})
		}
		return nil
	})
	return files, nil
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
