package chat

import (
	"log"
	"net"
	"os"
)

func Greet(conn net.Conn, filename string) {
	file, _ := os.Open(filename)
	data := make([]byte, 10000)
	n, _ := file.Read(data)
	conn.Write([]byte(data[:n]))
	file.Close()
}

func WriteLog(s string) {
	file, _ := os.OpenFile("log.txt", os.O_WRONLY|os.O_APPEND, 0666)
	_, err := file.WriteString(s + "\n")
	if err != nil {
		log.Fatal()
	}
	file.Close()
}
