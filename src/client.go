/*
   Created by XiaoInk at 2021/3/25 15:18
   GitHub: https://github.com/XiaoInk
*/

package src

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	name       string
	flag       int
	chatMsg    string
	conn       net.Conn
}

func NewClient(ServerIp string, ServerPort int) *Client {
	return &Client{
		ServerIp:   ServerIp,
		ServerPort: ServerPort,
		flag:       -1,
	}
}

func (c *Client) Run() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerIp, c.ServerPort))
	if err != nil {
		log.Println("net.Dial err: ", err)
	}
	c.conn = conn
	defer conn.Close()

	log.Println("服务器连接成功")

	// 监听服务器发送的消息
	go func() {
		_, _ = io.Copy(os.Stdout, c.conn)
	}()

	// 打印菜单
	for {
		if c.flag != 0 {
			if c.Menu() {
				break
			}
		}
	}
}

func (c *Client) Menu() bool {
	fmt.Println(">>> 1. 公聊模式")
	fmt.Println(">>> 2. 私聊模式")
	fmt.Println(">>> 3. 修改用户名")
	fmt.Println(">>> 0. 退出")

	_, _ = fmt.Scanln(&c.flag)

	switch c.flag {
	case 1: // 群聊模式
		c.PublicChat()
	case 2: // 私聊模式
		c.PrivateChat()
	case 3: // 修改用户名
		c.Rename()
	case 0: // 退出
		return true
	}
	return false
}

func (c *Client) PublicChat() {
	fmt.Println(">>> 您已进入公聊模式，输入 ?exit 退出公聊模式")
	for {
		_, _ = fmt.Scanln(&c.chatMsg)
		if c.chatMsg == "?exit" {
			break
		}

		c.SendMsg(c.chatMsg) // 向服务器发送消息
	}
}

func (c *Client) PrivateChat() {
	var toUser string

	c.SendMsg("?who") // 获取在线用户列表

	fmt.Println(">>> 您已进入私聊模式，请输入聊天对象[用户名]，输入 ?exit 退出私聊模式")
	_, _ = fmt.Scanln(&toUser)
	if toUser == "?exit" {
		return
	}

	fmt.Println(">>> 请输入聊天内容")
	for {
		_, _ = fmt.Scanln(&c.chatMsg)
		if c.chatMsg == "?exit" {
			break
		}

		c.SendMsg("?to:" + toUser + ":" + c.chatMsg) // 向给定用户发送消息
	}
}

func (c *Client) Rename() {
	fmt.Println(">>> 请输入新用户名，输入 ?exit 返回上一级")
	_, _ = fmt.Scanln(&c.name)
	c.SendMsg("?rename:" + c.name)
}

func (c *Client) SendMsg(msg string) {
	_, _ = c.conn.Write([]byte(msg + "\n"))
}
