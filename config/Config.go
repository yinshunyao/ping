package config

import (
	"net"
)

//运行参数
var (
	One bool = true
	//	Ip      string     = "127.0.0.1"
	Rate    int        = 200
	Retry   int        = 1
	Debug   bool       = false
	Timeout int        = 4
	Src     net.IPAddr = net.IPAddr{IP: net.ParseIP("0.0.0.0")}
)
