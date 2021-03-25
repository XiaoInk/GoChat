/*
   Created by XiaoInk at 2021/3/25 12:52
   GitHub: https://github.com/XiaoInk
*/

package src

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port      int
	Message   chan string      // 广播消息
	OnlineMap map[string]*User // 在线用户
	mapLock   sync.RWMutex
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		Message:   make(chan string),
		OnlineMap: make(map[string]*User),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Println("net.Listen err: ", err)
	}
	defer listener.Close()

	// 监听广播消息
	go func() {
		for {
			msg := <-s.Message
			s.mapLock.Lock()
			for _, user := range s.OnlineMap {
				user.C <- msg
			}
			s.mapLock.Unlock()
		}
	}()

	// 监听客户端连接请求
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept err: ", err)
			continue
		}

		// 处理业务逻辑
		go s.Handler(conn)
	}
}

// 业务处理
func (s *Server) Handler(conn net.Conn) {
	log.Println("客户端连接成功...")

	user := NewUser(conn, s)
	user.Online()

	isActive := make(chan bool)

	// 监听用户输入
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				log.Println("conn.Read err: ", err)
				return
			}

			// 处理用户输入
			user.DoMessage(string(buf[:n-1])) // buf[:n-1] 去除尾部换行符

			isActive <- true
		}
	}()

	// 超时强制关闭
	for {
		select {
		case <-isActive:
		case <-time.After(300 * time.Second):
			user.SendMsg("连接超时退出.")
			close(user.C)
			_ = conn.Close()
			return
		}
	}
}

// 消息广播
func (s *Server) Broadcast(user *User, msg string) {
	s.Message <- "[" + user.Addr + "]" + user.Name + ": " + msg
}
