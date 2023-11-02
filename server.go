package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const (
	DURATION = 60 * time.Second
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
			user.C <- msg
		}
		server.mapLock.Unlock()
	}
}

// broadcast user login message
func (server *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + ": " + msg

	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	user := NewUser(conn, server)

	user.Online()

	// check if user connection is alive
	isALive := make(chan bool)
	// create a timer
	timer := time.NewTimer(DURATION)

	// receive message sent by user
	go func() {
		sendBuf := make([]byte, 4096)

		for {
			// read the message sent by user
			n, err := conn.Read(sendBuf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Connect read err: ", err)
				return
			}

			isALive <- true

			// remove the last '\n'
			msg := string(sendBuf[:n-1])

			user.HandleMessage(msg)
		}
	}()

	for {
		select {
		case <-isALive: // user has sent message, reset timer
			timer.Reset(DURATION)
		case <-timer.C: // user is overtime, close connection
			timer.Stop()
			user.SendMessage("长时间未活跃，正在断开连接...\n")
			// close connection
			conn.Close()

			return
		}
	}
}

// get the list of online user
func (server *Server) GetOnlineUserList() string {
	var OnlineUserList string
	server.mapLock.Lock()
	for _, user := range server.OnlineMap {
		OnlineUserList += "[" + user.Addr + "] " + user.Name + "\n"
	}
	server.mapLock.Unlock()
	return OnlineUserList
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
