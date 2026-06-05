package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	invalidCommand = "invalid command"
)

type Client struct {
	UserName string `json:"userName"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

func NewClient(username, port, password string) *Client {
	return &Client{
		UserName: username,
		Port:     port,
		Password: password,
	}
}

func (c *Client) ConnectToServer() {
	address := fmt.Sprintf("localhost:%s", c.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	jsonReq, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("error in marshalling %v\n", err)
		return
	}

	fmt.Printf("Original Client request is %v\n", c)

	fmt.Printf("Marshalled Client request is %s\n", string(jsonReq))

	_, err = conn.Write(jsonReq)
	if err != nil {
		log.Fatalf("error in writing to server %v\n", err)
		return

	}

	// serverResp := []byte{}
	serverResp, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatalf("error in reading from server %v\n", err)
		return
	}
	if string(serverResp) == "invalid username" {
		log.Println("UserName is incorrect")
		return

	}

	if string(serverResp) == "invalid password" {
		log.Println("Password is incorrect")
		return

	}
	log.Println("FastStore Server : ", string(serverResp))

	inputReader := bufio.NewReader(os.Stdin)
	for {
		msg, _ := inputReader.ReadString('\n')

		validCommand := processRequest(msg)
		if validCommand == invalidCommand {
			log.Print("Server replied : ", "invalid command")
			continue
		}

		_, err = conn.Write([]byte(validCommand))
		if err != nil {
			log.Println("Write error:", err)
			return
		}

		serverReader := bufio.NewReader(conn)
		resp, err := serverReader.ReadString('\n')
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		log.Print("Server replied : ", resp)
	}
}

// processRequest read the user msg and convert to valid db command
func processRequest(msg string) string {
	splittedMsg := strings.Split(msg, " ")
	var req []string
	cmd := ""

	for _, v := range splittedMsg {
		if len(v) != 0 {
			v = strings.ToLower(v)
			req = append(req, v)
		}

	}
	if len(req) > 3 || len(req) < 1 {
		cmd = invalidCommand
	}

	switch req[0] {
	case "set":
		if len(req) != 3 {
			cmd = invalidCommand
		} else {
			cmd = fmt.Sprintf("%s|%s|%s", req[0], req[1], req[2])
		}
	case "get":
		if len(req) != 2 {
			cmd = invalidCommand
		} else {
			cmd = fmt.Sprintf("%s|%s", req[0], req[1])
		}
	case "del":
		if len(req) != 2 {
			cmd = invalidCommand
		} else {
			cmd = fmt.Sprintf("%s|%s", req[0], req[1])
		}
	default:
		cmd = invalidCommand

	}

	return cmd
}

func main() {

	port := flag.String("port", "6369", "FastStore server port")
	userName := flag.String("username", "", "username for db server")
	password := flag.String("password", "", "password for db server")
	flag.Parse()

	fmt.Println("Port is ", *port)

	client := NewClient(*userName, *port, *password)

	client.ConnectToServer()
}
