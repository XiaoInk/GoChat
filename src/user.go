/*
   Created by XiaoInk at 2021/3/25 13:13
   GitHub: https://github.com/XiaoInk
*/

package src

import (
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string // 消息频道
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Name:   addr,
		Addr:   addr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// 监听消息频道
	go func() {
		for {
			user.SendMsg(<-user.C)
		}
	}()
	return user
}

func (u *User) SendMsg(msg string) {
	_, _ = u.conn.Write([]byte(msg + "\n"))
}

func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	u.server.Broadcast(u, "上线啦...")
}

func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	u.server.Broadcast(u, "已下线.")
}

func (u *User) DoMessage(msg string) {
	u.server.Broadcast(u, msg)
}
