package main

import (
	"fmt"
	"net"

	"github.com/jrmhx/gochat/cmd/server/connmgr"
)

const (
	host             = "127.0.0.1:8080"
	MaxMsgNum        = 1024
	MaxMsgBufferSize = 1024
	MaxActiveConn    = 1
)

func echoHandler(conn net.Conn, cm *connmgr.ConnMgr) {
	defer conn.Close()
	buffer := make([]byte, MaxMsgBufferSize)
	for {
		if n, err := conn.Read(buffer); err == nil {
			fmt.Println(conn.RemoteAddr().String() + " " + string(buffer[:n]))
			conn.Write(buffer[:n])
		} else {
			fmt.Println(conn.RemoteAddr().String() + " is dead")
			cm.Delete(conn.RemoteAddr())
			return
		}
	}
}

func listenConn(ln net.Listener, cm *connmgr.ConnMgr) {
	for {
		if conn, err := ln.Accept(); err == nil {
			cm.Add(conn.RemoteAddr(), conn)
			go echoHandler(conn, cm)
		} else {
			fmt.Println(err)
		}
	}
}

func main() {
	ln, e := net.Listen("tcp", host)
	if e != nil {
		panic(fmt.Sprintf("cannot listen %v", host))
	}
	cm := connmgr.NewConnMgr(MaxActiveConn)
	go listenConn(ln, cm)

	for {

	}
}
