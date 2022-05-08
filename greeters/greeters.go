package chat

import (
	"log"
	"net"
	"os"
	"strings"
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

	conn.Write([]byte(str))
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
