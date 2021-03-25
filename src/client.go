/*
   Created by XiaoInk at 2021/3/25 15:18
   GitHub: https://github.com/XiaoInk
*/

package src

import (
	"fmt"
	"log"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Flag       int
	conn       net.Conn
}

func NewClient(ServerIp string, ServerPort int) *Client {
	return &Client{
		ServerIp:   ServerIp,
		ServerPort: ServerPort,
		Flag:       -1,
	}
}

func (c *Client) Run() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerIp, c.ServerPort))
	if err != nil {
		log.Println("net.Dial err: ", err)
	}
	defer conn.Close()

	log.Println("服务器连接成功")

	// 打印菜单
	for {
		if c.Flag != 0 {
			if c.Menu() {
				break
			}
		}
	}
}

func (c *Client) Menu() bool {
	fmt.Println("1. 群聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 修改用户名")
	fmt.Println("0. 退出")

	_, _ = fmt.Scanln(&c.Flag)

	switch c.Flag {
	case 1: // 群聊模式
		fmt.Println("群聊模式")
	case 2: // 私聊模式
		fmt.Println("私聊模式")
	case 3: // 修改用户名
		fmt.Println("修改用户名")
	case 0: // 退出
		return true
	}
	return false
}
