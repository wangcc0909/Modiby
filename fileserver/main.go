package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/download/goland-2020.2.2.dmg", downloadFile)
	http.ListenAndServe(":9090", nil)
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Path
	log.Println(s)
	dir, _ := os.Getwd()
	file := filepath.Join(dir, s)
	//fileName := path.Base(file)
	//fileName = url.QueryEscape(fileName) //防止中文乱码
	f, err := os.Open(file)
	if err != nil {
		log.Println("open file error")
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		log.Println("get file stat error")
		http.NotFound(w, r)
		return
	}
	w.Header().Add("Accept-Ranges", "bytes")
	w.Header().Add("Content-Disposition", "attachment; filename="+info.Name())
	var start, end int64
	if req := r.Header.Get("Range"); req != "" {
		if strings.Contains(req, "bytes=") && strings.Contains(req, "-") {
			fmt.Sscanf(req, "bytes=%d-%d", &start, &end)
			if end == 0 {
				end = info.Size() - 1
			}
			if start > end || start < 0 || end < 0 || end >= info.Size() {
				w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
				log.Println("download start:", start, " end:", end, " size:", info.Size())
				return
			}
			w.Header().Add("Content-Length", strconv.FormatInt(end-start+1, 10))
			w.Header().Add("Content_Range", fmt.Sprintf("bytes %v-%v/%v", start, end, info.Size()))
			w.WriteHeader(http.StatusPartialContent)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.Header().Add("Content-Length", strconv.FormatInt(info.Size(), 10))
		start = 0
		end = info.Size() - 1
	}
	w.Header().Add("Content-Type", "application/octet-stream")
	_, err = f.Seek(start, 0)
	if err != nil {
		log.Println("download:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	n := 512
	buf := make([]byte, n)
	for {
		if end-start+1 < int64(n) {
			n = int(end - start + 1)
		}
		_, err := f.Read(buf[:n])
		if err != nil {
			log.Println("1:", err.Error())
			if err != io.EOF {
				log.Println("error: ", err.Error())
			}
			return
		}
		err = nil
		_, err = w.Write(buf[:n])
		if err != nil {
			return
		}
		start += int64(n)
		if start >= end+1 {
			return
		}
	}
}
