package main

import (
	"fmt"
	"log"
	"net"
	"os"
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
	for {
		mtx.Lock()
		message := <-Mess

		// fmt.Printf("%v", message)
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

		mtx.Unlock()

	}
}

func WriteLog(s string) {
	file, _ := os.OpenFile("log.txt", os.O_WRONLY|os.O_APPEND, 0666)
	_, err := file.WriteString(s + "\n")
	if err != nil {
		log.Fatal()
	}
	file.Close()
}

func Writer(conn net.Conn, Mess chan string, user User) {
	fmt.Println("Writer")
	usernum.Lock()
	UserJoin := user.username[:len(user.username)-1] + " has joined our chat..."
	WriteLog(UserJoin)
	for i := 0; i < users; i++ {
		Mess <- UserJoin
	}
	usernum.Unlock()
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
			UserLeft := user.username[:len(user.username)-1] + " has left our chat..."
			WriteLog(UserLeft)
			for i := 0; i < users; i++ {
				Mess <- UserLeft
			}
			usernum.Unlock()

			fmt.Println("con closed")
			return
		}
		usernum.Lock()
		UserMessage := "[" + time.Now().Format(time.RFC822) + "][" + user.username[:len(user.username)-1] + "]:" + string(data[:message-1])
		WriteLog(UserMessage)
		for i := 0; i < users; i++ {
			Mess <- UserMessage
			// fmt.Println(i)
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
	str := ""
	split := strings.Split(string(data[:n]), "\n")
	for index, line := range split {
		if index != len(split)-1 {
			str += line + "\n"
		} else if index > 1 {
			str += line
		} else {
			str += "\n"
		}
	}

	conn.Write([]byte(str[:len(str)-1]))
	file.Close()
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
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
		Greet(conn, "greeting.txt")
		Greet(conn, "log.txt")
		n, err := conn.Read(data)
		if err != nil {
			log.Print(err)
		}

		user := User{username: string(data[:n])}

		go Writer(conn, Mess, user)
		go Messanger(conn, Mess, user)
		usernum.Lock()
		users++
		usernum.Unlock()
		time.Sleep(time.Second * 1)
	}
}
