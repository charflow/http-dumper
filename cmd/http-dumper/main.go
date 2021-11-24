package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	httpdumper "github.com/charflow/http-dumper"
)

func index(w http.ResponseWriter, req *http.Request) {
	dumped, _ := httpdumper.Dump(req)

	// 日志
	fmt.Println(string(dumped.HTTPData))
	fmt.Println("--------------------------------------------------")

	// 响应
	fmt.Fprintf(w, string(dumped.HTTPData))
	fmt.Fprintf(w, "\n--------------------------------------------------\n")
	fmt.Fprintf(w, MarshalJSON(dumped))

	// 自我请求
	if req.Header.Get("X-Pass") == "" {
		newData := bytes.Replace(dumped.HTTPData, []byte("\r\n\r\n"), []byte("\r\nX-Pass: true\r\n\r\n"), 1)
		// fmt.Println(string(newData))
		go httpdumper.DoRequest(listen, newData)
	}
}

var (
	listen = GetEnvDefault("LISTEN", ":10086")
)

func main() {
	http.HandleFunc("/", index)
	log.Printf("serve on %v", listen)
	http.ListenAndServe(listen, nil)
}

func GetEnvDefault(key, defaultV string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return defaultV
}

func MarshalJSON(in interface{}) string {
	bs, e := json.Marshal(in)
	if e != nil {
		panic(e)
	}
	return string(bs)
}
