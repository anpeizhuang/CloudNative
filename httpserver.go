package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
接收客户端 request，并将 request 中带的 header 写入 response header
读取当前系统的环境变量中的 VERSION 配置，并写入 response header
Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
当访问 localhost/healthz 时，应返回 200
*/

func main() {

	//http.HandleFunc("/", rootHandler)
	//err := http.ListenAndServe(":80", nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/setheader", respHandler)
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
	//RespCode := http.StatusText(err.Code)
	//fmt.Println("Response Code : " + string(RespCode))
}

func respHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("set resp header")
	//接收客户端 request，并将 request 中带的 header 写入 response header
	for k, v := range r.Header {
		w.Header().Add(k, strings.Join(v, " "))
	}

	//读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	var VERSION string
	VERSION = os.Getenv("VERSION")
	w.Header().Add("VERSION", VERSION)

	io.WriteString(w, "set resp header ok\n")

	//Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = strings.Split(r.RemoteAddr, ":")[0]
	}
	//io.WriteString(w, IPAddress+"\n")
	fmt.Println("IPAddress : " + IPAddress + ";\n")

	//RespCode := http.StatusCode
	//RespCode := http.StatusText()
	//RespCode := w.Header().Get("Status")
	//fmt.Println("Response Code : " + string(RespCode))

}

func healthz(w http.ResponseWriter, r *http.Request) {
	//当访问 localhost/healthz 时，应返回 200
	io.WriteString(w, "200\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}

	w.Header().Add("foo", "bar1")

	io.WriteString(w, "===================Details of the http request header:============\n")

	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}
