package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP        string
	Port      uint16
	OnlineMap map[string]*User // online user list
	mapLock   sync.RWMutex
	Message   chan string // channel for broadcasting
}

// create a new server
func NewServer(ip string, port uint16) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// listen to the server's channel
// if there is data, send it to every online user
func (server *Server) ListenMessage() {
	for {
		msg := <-server.Message

		// send msg to every online user
		server.mapLock.Lock()
		for _, user := range server.OnlineMap {
			user.Channel <- msg
		}
		server.mapLock.Unlock()
	}
}

// broadcast user login message
func (server *Server) BroadcastUserLogin(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + " " + msg

	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	user := NewUser(conn)

	// broadcast user login message
	server.BroadcastUserLogin(user, "已上线")

	// time.Sleep(1 * time.Second)
	// add new user to onlineMap
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	// block
	select {}
}

// start server
func (server *Server) Start() {
	// Socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	// close socket listen
	defer listener.Close()

	// start listening to server's Message
	go server.ListenMessage()

	for {
		// accept client connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept err: ", err)
			continue
		}

		fmt.Println("Connect client: ", conn.RemoteAddr().String())

		// do handler
		go server.Handler(conn)
	}
}
