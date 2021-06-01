package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	s := strings.Join([]string{"a"}, " AND ")
	fmt.Println(s)
	http.HandleFunc("/get", resultSuccess)
	http.ListenAndServe(":9090", nil)
}

func resultSuccess(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("hello world"))
}
