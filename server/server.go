package server

import (
	"bufio"
	"encoding/json"
	"fast-store/internals"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	UserName     string `json:"userName"`
	Port         string `json:"port"`
	Password     string `json:"password"`
	CacheService internals.FastStoreCache
}

func NewServer(username, port, password string, cacheService internals.FastStoreCache) *Server {
	return &Server{
		UserName:     username,
		Port:         port,
		Password:     password,
		CacheService: cacheService,
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

		serverResp, err := s.storeData(msg)
		if err != nil {
			conn.Write([]byte("NOT_OK | " + err.Error()))
			continue
		}

		conn.Write([]byte("OK | " + serverResp))
	}
}

func (s *Server) storeData(command string) (string, error) {
	cmds := strings.Split(command, "|")
	if len(cmds) < 2 || len(cmds) > 3 {
		return "", invalid_command
	}

	switch strings.ToLower(cmds[0]) {
	case "set":
		if len(cmds) != 3 {
			return "", invalid_command
		}
		err := s.CacheService.Set(cmds[1], cmds[2])
		if err != nil {
			log.Printf("error in set for %s is %v\n ", command, err)
			return "", server_error
		}
		return "", nil

	case "get":
		if len(cmds) != 2 {
			return "", invalid_command
		}
		val, err := s.CacheService.Get(cmds[1])
		if err != nil {
			log.Printf("error in set for %s is %v\n ", command, err)
			return "", server_error
		}
		return val, nil

	case "del":
		if len(cmds) != 2 {
			return "", invalid_command
		}
		err := s.CacheService.Delete(cmds[1])
		if err != nil {
			log.Printf("error in set for %s is %v\n ", command, err)
			return "", server_error
		}
		return "", nil

	default:
		return "", invalid_command
	}
}
