package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Server struct {
	UserName string
	Port     string
	Password string
}

func NewServer(username, port, password string) *Server {
	return &Server{
		UserName: username,
		Port:     port,
		Password: password,
	}
}

func (s *Server) StartServer() {

	fmt.Printf("%v\n",s)
	address := fmt.Sprintf(":%s", s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	log.Println("FastStore server is ready for accepting connection on port ", s.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("unable to accept connection:", err)
			continue
		}
		s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("connection closed :", err)
			return
		}

		fmt.Println("msg received is ", msg)

		conn.Write([]byte("Msg received is " + msg))
	}
}
