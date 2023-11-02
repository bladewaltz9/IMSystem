package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIP   string
	ServerPort uint16
	Name       string
	conn       net.Conn
}

func NewClient(serverIP string, serverPort uint16) *Client {
	// create client object
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
	}

	// connect to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return nil
	}
	client.conn = conn

	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>> connect failed")
		return
	}

	fmt.Println(">>>> connect success")

	select {}
}
