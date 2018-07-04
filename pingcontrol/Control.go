package pingcontrol

/*
流程控制函数实现
*/
import (
	"Ping/config"
	"Ping/ipmanager"
	"Ping/pingbase"
	"Ping/pingpool"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

//队列定义
var (
	conn_max_size int = 16
)

type Control struct {
	ip_count int64
	pool     *pingpool.PingPool
	todo     *ipmanager.IPSlice
	scan     *ipmanager.Scan
	wait     *ipmanager.Wait
	finish   *ipmanager.IPSlice
}

//解析IP参数，初始化todo列表，不需要加锁
func (control *Control) Parse(ips string) {
	ip_paramse := strings.Split(ips, "-")
	start := net.ParseIP(ip_paramse[0])
	if start == nil {
		if config.Debug {
			fmt.Println("无法识别参数：%s", ips)
		}
		return
	}

	ip_min := pingbase.IPInt(start)
	//单个IP
	if len(ip_paramse) == 1 {
		control.todo.Add(ipmanager.GenIPPing(ip_min))
		control.ip_count += 1
		return
	}

	end := net.ParseIP(ip_paramse[1])
	//结束IP非法
	if end == nil {
		if config.Debug {
			fmt.Println("无法识别参数：%s", ips)
		}
		return
	}

	//起始和结束IP一样
	if end.Equal(start) {
		return
	}
	//生成IP列表，按需生成，Ping的时候随机Ping
	ip_max := pingbase.IPInt(end)
	for ip_min <= ip_max {
		control.todo.Add(ipmanager.GenIPPing(ip_min))
		ip_min += 1
		control.ip_count += 1
	}
}

func (control *Control) result(addr string, ping_ok bool) {
	//Ping 成功，结束
	if ping_ok {
		ipping := control.wait.Success(addr)
		if ipping != nil {
			control.finish.AddSecurity(ipping)
		}

	} else { // Ping失败，重试或者结束
		ipping, flag := control.wait.Fail(addr)
		if flag > 0 {
			control.scan.Put(ipping)
		} else if flag < 0 {
			// 需要安全添加结果

			if ipping != nil {
				control.finish.AddSecurity(ipping)
			}
		}
		//理论上addr已经找不到对应的Ping实例
	}

}

//接收icmp包进行处理
func (control *Control) recv(conn *net.IPConn) {
	var (
		//		n        int
		//		oobn     int
		//		flags    int
		addr     *net.IPAddr
		err_recv error
	)

	recv := make([]byte, 1024)
	oob := make([]byte, 1024*1000)
	//未结束
	for control.finish.Length() < control.ip_count {
		_, _, _, addr, err_recv = conn.ReadMsgIP(recv, oob)
		if err_recv != nil {
			if config.Debug {
				fmt.Printf("ICMP监听异常:%v\n", err_recv.Error())
			}
			continue
		}
		//Ping成功
		go control.result(addr.IP.String(), true)
	}

	defer conn.Close()

}

// 针对某个IP进行Ping，并进行队列处理
func (control *Control) ping(ipping *ipmanager.IPPing, conn *net.IPConn) {
	var (
		wait bool
	)

	if config.Debug {
		fmt.Printf("%v扫描开始， %v-%v\n", ipping.String(), control.finish.Length(), control.ip_count)
	}

	//移动到等待队列中，重传时自动去重
	control.wait.AddSecurity(ipping)
	//Ping
	wait = ipping.Ping(conn)
	//仅发送，需要启动定时器等待
	if wait {
		time.Sleep(time.Duration(config.Timeout) * time.Second)
	}

	//需要重传
	if ipping.CheckRetry() {
		control.scan.Put(ipping)
		//		if config.Debug {
		//			fmt.Printf("%v需要重传， scan size-%v\tfinish size-%v\tall-%v\n", ipping.String(), control.scan.Size(), control.finish.Length(), control.ip_count)
		//		}
	} else { //不需要重传
		//		if config.Debug {
		//			fmt.Printf("%v扫描结束， %v-%v\n", ipping.String(), control.finish.Length(), control.ip_count)
		//		}
		//接收Ping包线程可能移动ipping到finish数组中
		ipping = control.wait.RmvSecurity(ipping.String())
		if ipping != nil {
			control.finish.AddSecurity(ipping)
		} else {
			if !wait {
				fmt.Printf("获取完成的Ping任务异常")
			}
		}

	}
}

func (control *Control) toScan() {
	var (
		size   int64 = control.todo.Length()
		index  int64
		ipping *ipmanager.IPPing
	)

	for size > 0 {
		rand.Seed(time.Now().UnixNano())
		index = rand.Int63n(size)
		ipping = control.todo.Pop(index)
		if ipping != nil {
			control.scan.Put(ipping)
		} else {
			break
		}

		size = control.todo.Length()
	}

}

func (control *Control) send() {
	var (
		dur    int64
		n      int64
		t_dur  time.Duration
		ipping *ipmanager.IPPing
	)

	n = time.Second.Nanoseconds()
	if config.Rate < ipmanager.MAX_RATE {
		dur = n / int64(config.Rate)
	} else {
		fmt.Printf("发包速率不应该超过%v", ipmanager.MAX_RATE)
		return
	}
	t_dur = time.Duration(dur)

	//会获取锁
	for control.finish.Length() < control.ip_count {
		//循环检查是否有新的IP需要发送
		ipping = control.scan.Get()
		if ipping == nil {
			//pass
			if config.Debug {
				fmt.Printf("获取Ping任务失败:todo-%v, done-%v\n", control.ip_count, control.finish.Length())
			}
			continue
		}
		//异步Ping
		go control.ping(ipping, control.pool.GetConn())
		time.Sleep(t_dur)
	}

	if config.Debug {
		fmt.Printf("发包结束：todo-%v, done-%v\n", control.ip_count, control.finish.Length())
		control.pool.Print()
	}
}

//主运行函数
func (control *Control) Run() bool {
	//先初始化连接池
	if config.Debug {
		fmt.Printf("需要扫描%v个IP\n", control.ip_count)
		fmt.Printf("队列长度：%.0f\n", ipmanager.MAX_QUEUE_SIZE)
	}
	// todo 根据速率来设定连接池大小
	conn_max_size = int(config.Rate) / 50
	control.pool = pingpool.GeneratePool(&config.Src, conn_max_size)
	//监听连接池的所有端口
	var conns []*net.IPConn = control.pool.AllConn()
	if conns == nil {
		if config.Debug {
			fmt.Println("连接池初始化失败，每个Ping包将新建连接")
		}
	} else {
		if config.Debug {
			fmt.Printf("%v个连接池初始化成功\n", len(conns))
		}

		//监听所有线程池，接收所有Ping回包
		for index := 0; index < len(conns); index++ {
			go control.recv(conns[index])
		}
	}
	//异步增加队列，从todo队列到scan队列
	control.toScan()
	//发送Ping包
	control.send()
	return true
}

//统计总共扫描了多少个IP
func (control *Control) FinishCount() int64 {
	return control.finish.Length()
}

//打印存活结果
func (control *Control) Print() {
	control.finish.AliveList(true)
}

func GenerateControl() *Control {
	return &Control{
		ip_count: int64(0),
		todo:     ipmanager.GenerateIPSlice(),
		scan:     ipmanager.GenerateIPScan(),
		wait:     ipmanager.GenerateIPWait(),
		finish:   ipmanager.GenerateIPSlice(),
	}
}
