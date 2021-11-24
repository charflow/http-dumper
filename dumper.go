package httpdumper

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
)

type DumpedRequest struct {
	HTTPData []byte `json:"-" `
	Client   string `json:"client"`
}

func Dump(req *http.Request) (*DumpedRequest, error) {
	var buf = &bytes.Buffer{}
	if e := req.Write(buf); e != nil {
		return nil, e
	}
	bs := buf.Bytes()
	return &DumpedRequest{bs, req.RemoteAddr}, nil
}

func DoRequest(target string, data []byte) error {
	conn, e := net.Dial("tcp", target)
	if e != nil {
		return e
	}
	defer conn.Close()

	_, e = conn.Write(data)
	func() {
		_, e := io.ReadAll(conn)
		if e != nil {
			log.Println("response err:", e)
		}
		// log.Println("response:", string(res))
	}()
	return e
}
