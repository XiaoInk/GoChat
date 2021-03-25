/*
   Created by XiaoInk at 2021/3/25 12:56
   GitHub: https://github.com/XiaoInk
*/

package main

import (
	"flag"
	"gochat/src"
)

var Ip string
var Port int

func init() {
	flag.StringVar(&Ip, "ip", "127.0.0.1", "监听地址")
	flag.IntVar(&Port, "port", 8888, "监听端口")
}

func main() {
	flag.Parse()

	server := src.NewServer(Ip, Port)
	server.Start()
}
