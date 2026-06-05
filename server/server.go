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

	conn, err := listener.Accept()
	if err != nil {
		log.Println("unable to accept connection:", err)
		return
	}
	defer conn.Close()

	data := Server{}

	err = json.NewDecoder(conn).Decode(&data)
	if err != nil {
		log.Println("error in decoding:", err)
		return
	}

	fmt.Printf("Request from client is %+v\n", data)
	if data.UserName != s.UserName {
		fmt.Println(data.UserName, " ", s.UserName)
		conn.Write([]byte("Invalid UserName \n"))
		return
	}

	if data.Password != s.Password {
		fmt.Println("invalid password")
		conn.Write([]byte("Invalid Password \n"))
		return
	}

	conn.Write([]byte("FastStore server is ready to accept request \n"))

	log.Println("FastStore server is ready for accepting connection on port ", s.Port)
	conn.Close()

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
