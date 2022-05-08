package Run

import (
	"fmt"
	"log"
	"net"
	"os"
	chat "server/greeters"
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
		if chat.Users == 10 {
			conn.Write([]byte("Server full"))
			conn.Close()
			continue
		}
		var user chat.User
		go func() {
			var ct int
			for {
				conn.Write([]byte("Your username: "))
				data := make([]byte, 100)
				n, err := conn.Read(data)
				if err != nil {
					log.Print(err)
				}
				fmt.Printf("n is %d", n)
				if n != 1 {
					user = chat.User{Username: string(data[:n])}
					break
				}
				if ct == 2 {
					conn.Close()
					return
				}
				ct++
			}
			chat.Greet(conn, "greeting.txt")
			chat.Greet(conn, "log.txt")
			go chat.Writer(conn, Mess, user)
			go chat.Messanger(conn, Mess, user)
			chat.Usernum.Lock()
			chat.Users++
			chat.Usernum.Unlock()
			//time.Sleep(time.Second * 1)
		}()
	}
}
