package pingpool

import (
	"net"
)

type PingConn struct {
	conn  *net.IPConn //socket连接
	count int64       //使用计数
	//	enable bool
}
