package ipmanager

/*
IP简单计算
*/
import (
	"Ping/config"
	"Ping/pingbase"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func init() {
	var (
		length int = len(originaBytes)
		index  int = 0
	)
	for index < length {
		originaBytes[index] = byte(index)
		index++
	}
}

type IPPing struct {
	ip    net.IPAddr //IP
	try   int        //重试次数
	alive bool       //存活次数
	scan  bool       //扫描状态
	//	lock  *sync.Mutex
}

//生成IPPing
func GenIPPing(ip_int int64) *IPPing {
	return &IPPing{ip: net.IPAddr{IP: pingbase.IPFromInt(ip_int)}, try: 0, alive: false, scan: false}
}

//判断是否需要重传
func (ipping *IPPing) CheckRetry() bool {
	if ipping.alive || ipping.try >= config.Retry {
		return false
	} else {
		return true
	}
}

/*
Ping
可能用公共的Socket连接进行Ping，也可能在本函数中新建Socket连接同步模式进行Ping
返回结果
表示 是否需要wait回包

*/
func (ipping *IPPing) Ping(conn *net.IPConn) bool {
	var (
		icmp pingbase.ICMP
		//max_lan, min_lan, avg_lan float64
	)

	// 尝试次数加1
	ipping.try += 1

	//本地连接标记
	local_flag := false
	//如果公共连接不存在，新建连接
	if conn == nil {
		local_flag = true
		conn, err = net.DialIP("ip4:icmp", &config.Src, &ipping.ip)
		if err != nil {
			if config.Debug {
				fmt.Printf(err.Error())
			}

			return false
		}
		defer conn.Close()

	}

	icmp.Type = 8
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequencyNum = 0

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	binary.Write(&buffer, binary.BigEndian, originaBytes)
	b := buffer.Bytes()
	binary.BigEndian.PutUint16(b[2:], pingbase.CheckSum(b))

	//	if config.Debug {
	//		fmt.Printf("\n正在Ping %s 具有 %d(%d)字节的数据:\n", ipping.ip.String(), MAX_PG, MAX_PG+28)
	//	}

	recv := make([]byte, 1024)
	//	t_start := time.Now()
	//直接休眠
	if !local_flag { //使用全局的连接
		//发送消息到某个IP
		var (
			n   int
			err error
		)
		try := 0
		for try < 3 {
			try++
			if n, err = conn.WriteToIP(buffer.Bytes(), &ipping.ip); err != nil {
				//				if config.Debug {
				//					fmt.Printf("发送到 %v 异常: %v， %v\n", ipping.ip.String(), n, err.Error())
				//				}
				//发送失败
				time.Sleep(time.Duration(1000) * time.Millisecond)
				continue
				//				return false
			} else {
				return true
			}

		}
		if config.Debug {
			fmt.Printf("重试了%v次发送到 %v 异常: %v, %v\n", try, ipping.ip.String(), n, err.Error())
		}
		// 发送成功
		return false
	} else { //本地新建连接
		//		t_start := time.Now()
		if _, err := conn.Write(buffer.Bytes()); err != nil {
			if config.Debug {
				fmt.Printf("发送到 %s 异常:%v\n", ipping.ip.String(), err.Error())
			}
			return false
		}
		conn.SetReadDeadline((time.Now().Add(time.Duration(config.Timeout) * time.Second)))
		_, err = conn.Read(recv)
		if err != nil {
			//			if config.Debug {
			//				fmt.Printf("来自 %s 的回复超时\n", ipping.ip.String())
			//			}
			return false
		}
		//		t_end := time.Now()
		//		dur := float64(t_end.Sub(t_start).Nanoseconds()) / 1e6
		//		if config.Debug {
		//			fmt.Printf("来自 %s 的回复：时间 = %.3fms\n", ipping.ip.String(), dur)
		//		}

		//Ping成功
		ipping.Success()
		return false
	}
}

//Ping成功处理
func (ipping *IPPing) Success() {
	ipping.alive = true
	ipping.try = config.Retry
}

func (ipping *IPPing) String() string {
	return ipping.ip.String()
}
