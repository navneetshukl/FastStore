package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func ConnectToServer(){
	conn,err:=net.Dial("tcp","localhost:6969")
	if err!=nil{
		panic(err)
	}

	defer conn.Close()

	log.Println("Connected to FastStore server")

	inputReader:=bufio.NewReader(os.Stdin)
	for {
		msg, _ := inputReader.ReadString('\n')

		_, err = conn.Write([]byte(msg))
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

		log.Print("Server replied:", resp)
	}
}

func main(){
	ConnectToServer()
}