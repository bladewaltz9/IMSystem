package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string

	C      chan string
	conn   net.Conn
	server *Server
}

// create a new user
func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Name:   conn.RemoteAddr().String(),
		Addr:   conn.RemoteAddr().String(),
		C:      make(chan string),
		conn:   conn,
		server: server,
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

// send message to user
func (user *User) SendMessage(msg string) {
	user.conn.Write([]byte(msg))
}

// Handle user's Message
func (user *User) HandleMessage(msg string) {
	if msg == "who" {
		// send online user list to the user
		OnlineUserList := user.server.GetOnlineUserList()
		user.SendMessage(OnlineUserList)
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		user.Rename(msg[7:])
	} else if len(msg) > 3 && msg[:3] == "to|" {
		msgList := strings.Split(msg, "|")
		if len(msgList) != 3 || len(msgList[1]) == 0 {
			user.SendMessage("格式错误！Usage: to|<username>|<message>\n")
			return
		} else if len(msgList[2]) == 0 {
			user.SendMessage("发送内容不能为空！\n")
			return
		}

		name, content := msgList[1], msgList[2]
		// get user by name and check if it's exists.
		remoteUser, ok := user.server.OnlineMap[name]
		if !ok {
			user.SendMessage("用户名不存在！\n")
			return
		}

		// send content to remote user
		remoteUser.SendMessage("[from " + user.Name + "]: " + content + "\n")
	} else {
		// broadcast message sent by user
		user.server.Broadcast(user, msg)
	}

}

// rename user
func (user *User) Rename(name string) bool {
	// check if the username already exists
	if _, value := user.server.OnlineMap[name]; value {
		user.SendMessage("用户名已存在！\n")
		return false
	} else {
		// delete old map
		user.server.mapLock.Lock()
		delete(user.server.OnlineMap, user.Name)
		user.server.mapLock.Unlock()

		// add new map
		user.server.OnlineMap[name] = user
		user.Name = name
		user.SendMessage("用户名已改为 " + user.Name + "\n")
		return true
	}
}

// send msg to user if there is data in user's channel
func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
