package ipmanager

import (
	//	"net"
	"time"
)

const (
	MAX_PG         = 1
	MAX_RATE       = 1e7 //最大速率
	MAX_QUEUE_SIZE = 1e8 //最大队列长度
	TIMEOUT_QUEUE  = float64(1)
)

// 默认的一些参数配置
var (
	err          error
	MS_20        = time.Duration(time.Second.Nanoseconds() / 50)
	originaBytes = make([]byte, MAX_PG) //Ping净荷
)
