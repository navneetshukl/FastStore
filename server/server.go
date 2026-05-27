package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func StartServer() {
	listener, err := net.Listen("tcp", ":6969")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	log.Println("FastStore server is ready for accepting connection on port 6969")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("unable to accept connection:", err)
			continue
		}
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("connection closed :", err)
			return
		}

		fmt.Println("msg received is ",msg)

		conn.Write([]byte("Msg received is " + msg))
	}
}
