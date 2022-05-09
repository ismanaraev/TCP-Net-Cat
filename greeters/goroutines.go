package chat

import (
	"log"
	"net"
	"sync"
	"time"
)

//takes data from Writer by channel and sends it to clients
func Messanger(conn net.Conn, Mess chan string, user User) {
	var mtx sync.Mutex
	for {
		mtx.Lock()
		message := <-Mess

		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Print(err)
			Usernum.Lock()
			Users--
			Usernum.Unlock()
			mtx.Unlock()
			return
		}

		mtx.Unlock()

	}
}

//This function Takes input from client through conn and sends it to Messanger through channel
func Writer(conn net.Conn, Mess chan string, user User) {
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

			return
		}
		Usernum.Lock()
		UserMessage := "[" + time.Now().Format("2006-01-02 15:04:05") + "][" + user.Username[:len(user.Username)-1] + "]:" + string(data[:message-1])
		WriteLog(UserMessage)
		for i := 0; i < Users; i++ {
			Mess <- UserMessage
		}
		Usernum.Unlock()
	}
}
