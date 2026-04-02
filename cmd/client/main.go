package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

const host = "127.0.0.1:8080"

func main() {
	var conn net.Conn
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if conn, err = net.Dial("tcp", host); err != nil {
			fmt.Println("failed to connect " + host + " retry...")
			time.Sleep(500 * time.Millisecond)
		} else {
			fmt.Println("connected to " + conn.RemoteAddr().String())
			break
		}
	}
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			wbuffer := scanner.Bytes()
			conn.Write(wbuffer)
			rbuffer := make([]byte, 1024)
			if _, err := conn.Read(rbuffer); err == nil {
				fmt.Println(string(rbuffer))
			} else {
				conn.Close()
				break
			}
		} else {
			continue
		}
	}
}
