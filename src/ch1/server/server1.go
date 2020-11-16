// web server1
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler1)
	log.Fatal(http.ListenAndServe("localhost:20000", nil))
}

// 处理程序回显请求URL r的路径部分
func handler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}
