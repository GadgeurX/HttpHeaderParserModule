package main

import (
	"Airttp/modules"
	"net/rpc"
	"net"
	"log"
	"fmt"
	"bytes"
	"strings"
)

type Http int

func main() {
	http := new(Http)

	server := rpc.NewServer()
	server.RegisterName("Http", http)

	l, e := net.Listen("tcp", ":5000")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	fmt.Println("server start")
	server.Accept(l)
}

func (t *Http) Module(params modules.ModuleParams, result *modules.ModuleParams) error {
	fmt.Print("New Request : ")
	result.Copy(params)
	result.Req.Header = bytes.Split(result.Req.Raw, []byte("\r\n\r\n"))[0]
	headerLines := bytes.Split(result.Req.Header, []byte("\r\n"))
	requestLine := headerLines[0]
	result.Req.Method = string(bytes.Split(requestLine, []byte(" "))[0])
	result.Req.Uri = string(bytes.Split(requestLine, []byte(" "))[1])
	headerLines = append(headerLines[:0], headerLines[0+1:]...)
	result.Req.Headers = make(map[string]string)
	for _, line := range headerLines {
		key := strings.TrimSpace(string(bytes.Split(line, []byte(":"))[0]))
		value := strings.TrimSpace(string(bytes.Split(line, []byte(":"))[1]))
		result.Req.Headers[key] = value
	}
	fmt.Println("OK")
	return nil
}