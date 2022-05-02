package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type User struct {
	username string
}

func Messanger(conn net.Conn, Mess chan string, user User) {
	var mtx sync.Mutex
	fmt.Println("Messanger")
	file, _ := os.OpenFile("log.txt", os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	for {
		mtx.Lock()
		message := <-Mess
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Print(err)
			usernum.Lock()
			users--
			usernum.Unlock()
			mtx.Unlock()
			fmt.Println("con closed 1")
			return
		}

		_, err = file.WriteString(message)
		if err != nil {
			log.Fatal()
		}
		mtx.Unlock()
	}
}

func Writer(conn net.Conn, Mess chan string, user User) {
	fmt.Println("Writer")
	for {
		data := make([]byte, 200)
		message, err := conn.Read(data)
		if message == 1 {
			continue
		}
		if err != nil {
			log.Print(err)

			conn.Close()
			usernum.Lock()
			for i := 0; i < users; i++ {
				Mess <- user.username[:len(user.username)-1] + " has left our chat..." + "\n"
			}
			usernum.Unlock()

			fmt.Println("con closed")
			return
		}
		usernum.Lock()
		for i := 0; i < users; i++ {
			Mess <- "[" + time.Now().Format(time.RFC822) + "][" + user.username[:len(user.username)-1] + "]:" + string(data[:message])
		}
		usernum.Unlock()
	}
}

var (
	users   int
	usernum sync.Mutex
)

func Greet(conn net.Conn, filename string) {
	file, _ := os.Open(filename)
	data := make([]byte, 10000)
	n, _ := file.Read(data)
	for _, line := range strings.Split(string(data[:n]), "\n") {
		conn.Write([]byte(line + "\n"))
	}
	file.Close()
}

func main() {
	l, err := net.Listen("tcp", ":8989")
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

		Greet(conn, "greeting.txt")
		Greet(conn, "log.txt")
		conn.Write([]byte("Your username: "))
		data := make([]byte, 100)
		n, err := conn.Read(data)
		if err != nil {
			log.Print(err)
		}

		user := User{username: string(data[:n])}
		go Writer(conn, Mess, user)
		go Messanger(conn, Mess, user)
		Mess <- user.username[:len(user.username)-1] + " has joined our chat..." + "\n"
		time.Sleep(time.Second)
		usernum.Lock()
		users++
		usernum.Unlock()
		fmt.Println(runtime.NumGoroutine())
	}
}
