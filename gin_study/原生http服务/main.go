package main

import (
	"fmt"
	"io"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------", r.Method, r.URL.String())

	if r.Method != "GET" {
		byteData, _ := io.ReadAll(r.Body)
		fmt.Println("-----------", string(byteData)) // 转换为字符串以便查看
	}

	fmt.Println(r.Header)

	w.Write([]byte("hello world"))
}

func main() {
	http.HandleFunc("/index", Index) // 添加前导斜杠
	fmt.Println("服务器启动在: http://127.0.0.1:8080/index")
	http.ListenAndServe("127.0.0.1:8080", nil)
}
