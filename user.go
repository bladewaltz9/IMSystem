package main

import "net"

type User struct {
	Name string
	Addr string

	Channel chan string
	conn    net.Conn
	server  *Server
}

// create a new user
func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Name:    conn.RemoteAddr().String(),
		Addr:    conn.RemoteAddr().String(),
		Channel: make(chan string),
		conn:    conn,
		server:  server,
	}

	// start to listen user message
	go user.ListenMessage()

	return user
}

// user online function
func (user *User) Online() {
	// broadcast user login message
	user.server.Broadcast(user, "已上线")

	// time.Sleep(1 * time.Second)
	// add new user to onlineMap
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()
}

// user offline function
func (user *User) Offline() {
	// delete user in onlineMap
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	user.server.Broadcast(user, "已下线")
}

// Handle user's Message
func (user *User) HandleMessage(msg string) {
	// broadcast message sent by user
	user.server.Broadcast(user, msg)
}

// send msg to user if there is data in user's channel
func (user *User) ListenMessage() {
	for {
		msg := <-user.Channel

		user.conn.Write([]byte(msg + "\n"))
	}
}
