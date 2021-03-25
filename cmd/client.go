/*
   Created by XiaoInk at 2021/3/25 15:20
   GitHub: https://github.com/XiaoInk
*/

package main

import (
	"flag"
	"gochat/src"
)

var ServerIp string
var ServerPort int

func init() {
	flag.StringVar(&ServerIp, "ip", "127.0.0.1", "服务器地址")
	flag.IntVar(&ServerPort, "port", 8888, "服务器端口")
}

func main() {
	flag.Parse()

	client := src.NewClient(ServerIp, ServerPort)
	client.Run()
}
