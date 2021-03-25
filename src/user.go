/*
   Created by XiaoInk at 2021/3/25 13:13
   GitHub: https://github.com/XiaoInk
*/

package src

import (
	"net"
	"strings"
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
	if msg == "?who" { // 查询在线用户列表
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			u.SendMsg("[" + user.Addr + "]" + user.Name + ": 在线")
		}
		u.server.mapLock.Unlock()
	} else if strings.HasPrefix(msg, "?rename:") { // 修改用户名
		newName := strings.Split(msg, ":")[1]
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("该用户名已被占用")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMsg("您的用户名已修改成功")
		}
	} else if strings.HasPrefix(msg, "?to:") { // 私聊
		toUserName, toUserMsg := strings.Split(msg, ":")[1], strings.Split(msg, ":")[2]
		toUser, ok := u.server.OnlineMap[toUserName]
		if ok {
			if len(strings.TrimSpace(toUserMsg)) > 0 {
				toUser.SendMsg("<From " + toUser.Name + "> " + toUserMsg)
				return
			}
			u.SendMsg("消息体不能为空，请输入: \"?to:Name:Content\"")
		} else {
			u.SendMsg("未知的用户名或用户不在线")
		}
	} else {
		u.server.Broadcast(u, msg)
	}
}
