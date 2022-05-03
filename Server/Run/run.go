package Run

import (
	"fmt"
	"log"
	"net"
	"os"
	chat "server/greeters"
	"time"
)

func Run() error {
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(0)
	}
	port := "8989"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	Mess := make(chan string)
	os.Create("log.txt")
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		conn.Write([]byte("Your username: "))
		data := make([]byte, 100)
		chat.Greet(conn, "greeting.txt")
		chat.Greet(conn, "log.txt")
		n, err := conn.Read(data)
		if err != nil {
			log.Print(err)
		}

		user := chat.User{Username: string(data[:n])}

		go chat.Writer(conn, Mess, user)
		go chat.Messanger(conn, Mess, user)
		chat.Usernum.Lock()
		chat.Users++
		chat.Usernum.Unlock()
		time.Sleep(time.Second * 1)
	}
}
