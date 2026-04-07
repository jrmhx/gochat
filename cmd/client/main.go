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
	scannerBuffer := make([]byte, 0, 1024)
	scanner.Buffer(scannerBuffer, 10240)
	buffer := make([]byte, 1024)
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
			if len(scanner.Bytes()) == 0 {
				continue
			}
			conn.Write(scanner.Bytes())
			if n, err := conn.Read(buffer); err == nil {
				fmt.Println(string(buffer[:n]))
			} else {
				conn.Close()
				break
			}
		} else {
			continue
		}
	}
}
