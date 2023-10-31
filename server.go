package main

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port uint16
}

// create a new server
func NewServer(ip string, port uint16) *Server {
	server := &Server{
		IP:   ip,
		Port: port,
	}

	return server
}

func (server *Server) Handler(conn net.Conn) {
	// do Handler
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

	for {
		// accept client connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept err: ", err)
			continue
		}

		fmt.Println("Client address: ", conn.RemoteAddr())

		// do handler
		go server.Handler(conn)
	}
}
