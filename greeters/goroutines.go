package chat

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

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
			Usernum.Lock()
			Users--
			Usernum.Unlock()
			mtx.Unlock()
			fmt.Println("con closed 1")
			return
		}

		mtx.Unlock()

	}
}

func Writer(conn net.Conn, Mess chan string, user User) {
	fmt.Println("Writer")
	Usernum.Lock()
	UserJoin := user.Username[:len(user.Username)-1] + " has joined our chat..."
	WriteLog(UserJoin)
	for i := 0; i < Users; i++ {
		Mess <- UserJoin
	}
	Usernum.Unlock()
	for {
		data := make([]byte, 200)
		message, err := conn.Read(data)
		if message == 1 {
			continue
		}
		if err != nil {
			log.Print(err)

			conn.Close()
			Usernum.Lock()
			UserLeft := user.Username[:len(user.Username)-1] + " has left our chat..."
			WriteLog(UserLeft)
			for i := 0; i < Users; i++ {
				Mess <- UserLeft
			}
			Usernum.Unlock()

			fmt.Println("con closed")
			return
		}
		Usernum.Lock()
		UserMessage := "[" + time.Now().Format(time.RFC822) + "][" + user.Username[:len(user.Username)-1] + "]:" + string(data[:message-1])
		WriteLog(UserMessage)
		for i := 0; i < Users; i++ {
			Mess <- UserMessage
			// fmt.Println(i)
		}
		Usernum.Unlock()
	}
}
