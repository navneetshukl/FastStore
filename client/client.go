package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	invalidCommand = "invalid command"
)

func ConnectToServer() {
	conn, err := net.Dial("tcp", "localhost:6969")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	log.Println("Connected to FastStore server")

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

//processRequest read the user msg and convert to valid db command
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
	ConnectToServer()
}
