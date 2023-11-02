package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIP string, serverPort int) *Client {
	// create client object
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		flag:       -1,
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

func (client *Client) Menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag < 0 || flag > 3 {
		fmt.Println(">>>>请输入范围内的数字(0-3)<<<<")
		return false
	} else {
		client.flag = flag
		return true
	}
}

func (client *Client) run() {
	for client.flag != 0 {
		for client.Menu() != true {
		}

		switch client.flag {
		case 1:
			fmt.Println("公聊模式...")
			break
		case 2:
			fmt.Println("私聊模式...")
			break
		case 3:
			fmt.Println("更新用户名...")
			break
		}
	}
}

var serverIP string
var serverPort int

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "set server ip(default: 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "set server port(default: 8888)")
}

func main() {
	flag.Parse()

	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println(">>>> connect failed")
		return
	}

	fmt.Println(">>>> connect success")

	client.run()
}
