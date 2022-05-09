Net-Cat

This program starts a group chat using TCP protocol. Package includes both server and client

INSTALLATION

To install a program, simply clone it to any directory. You need to have Go installed to run it 

USAGE


SERVER

To start a server, go into the package directory and launch a server using 

go run server.go [port]

or compile it 

go build -o TCPChat server.go
chmod +x TCPChat
./TCPChat [port]

if no port specified, the server will listen on 8989

CLIENT

To start a client, go into the client directory and launch a client using either 

go run client.go [host] [port] 

or compile it and run

go build -o client client.go
chmopd +x client
./client [host] [port]

Host and port spectification is mandatory


