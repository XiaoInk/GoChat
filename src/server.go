/*
   Created by XiaoInk at 2021/3/25 12:52
   GitHub: https://github.com/XiaoInk
*/

package src

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	return &Server{
		ip,
		port,
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Println("net.Listen err: ", err)
	}
	defer listener.Close()

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

func (s *Server) Handler(conn net.Conn) {
	log.Println("客户端连接成功...")
}
