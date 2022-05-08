package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func read(conn net.Conn) {
	for {
		data := make([]byte, 1000)
		n, err := conn.Read(data)
		if err != nil {
			return
		}
		fmt.Print(string(data[:n]))
	}
}
func main() {
	if len(os.Args) != 3 {
		fmt.Println("[USAGE]: [Host][Port]")
		return
	}
	host := os.Args[1]
	port := os.Args[2]
	client, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	go read(client)
	for {
		data := make([]byte, 1000)
		n, _ := os.Stdin.Read(data)
		client.Write(data[:n])
	}
}
