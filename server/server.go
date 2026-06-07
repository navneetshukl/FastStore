package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Server struct {
	UserName string `json:"userName"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

func NewServer(username, port, password string) *Server {
	return &Server{
		UserName: username,
		Port:     port,
		Password: password,
	}
}

func (s *Server) StartServer() {
	fmt.Printf("%v\n", s)
	address := fmt.Sprintf(":%s", s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	log.Println("FastStore server is listening on port", s.Port)

	for {
		fmt.Println("inside infinite loop")
		conn, err := listener.Accept()
		if err != nil {
			log.Println("unable to accept connection:", err)
			continue
		}

		if !s.authenticateClient(conn) {
			conn.Close()
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) authenticateClient(conn net.Conn) bool {
	data := Server{}
	err := json.NewDecoder(conn).Decode(&data)
	if err != nil {
		log.Println("error in decoding:", err)
		conn.Write([]byte("Invalid request format\n"))
		return false
	}

	fmt.Printf("Request from client is %+v\n", data)

	if data.UserName != s.UserName {
		fmt.Println(data.UserName, " ", s.UserName)
		conn.Write([]byte("Invalid UserName \n"))
		return false
	}

	if data.Password != s.Password {
		fmt.Println("invalid password")
		conn.Write([]byte("Invalid Password \n"))
		return false
	}

	conn.Write([]byte("FastStore server is ready to accept request \n"))
	return true
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
