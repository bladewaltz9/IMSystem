package main

import "net"

type User struct {
	Name string
	Addr string

	Channel chan string
	conn    net.Conn
}

// create a new user
func NewUser(conn net.Conn) *User {
	user := &User{
		Name:    conn.RemoteAddr().String(),
		Addr:    conn.RemoteAddr().String(),
		Channel: make(chan string),
		conn:    conn,
	}

	// start to listen user message
	go user.ListenMessage()

	return user
}

// send msg to user if there is data in user's channel
func (user *User) ListenMessage() {
	for {
		msg := <-user.Channel

		user.conn.Write([]byte(msg + "\n"))
	}
}
