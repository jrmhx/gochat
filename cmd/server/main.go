package main

import (
	"fmt"
	"net"
)

const host = "127.0.0.1:8080"

func echoHandler(ln *net.Listener) {
	// handle conn logic
	var conn net.Conn
	var err error
	conn, err = (*ln).Accept()
	if err != nil {
		return
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		if n, err := conn.Read(buffer); err == nil {
			fmt.Println(conn.RemoteAddr().String() + " " + string(buffer[:n]))
			conn.Write(buffer[:n])
		} else {
			return
		}
	}
}

func main() {
	ln, e := net.Listen("tcp", host)
	if e != nil {
		panic(fmt.Sprintf("cannot listen %v", host))
	}
	echoHandler(&ln)
}
