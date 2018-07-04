package pingpool

import (
	"fmt"
	//	"fmt"
	"Ping/config"
	"net"
	"runtime"
	"sync"
)

type PingPool struct {
	lock  *sync.Mutex
	conn  []*PingConn
	count int //轮询计数
	alive int
}

func (pool *PingPool) AllConn() []*net.IPConn {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	if pool.alive > 0 {
		var (
			conn  []*net.IPConn
			index int = 0
		)
		for index < pool.alive {
			conn = append(conn, pool.conn[index].conn)
			index++
		}
		return conn
	} else {
		return nil
	}

}

// 轮询获取空闲连接
func (pool *PingPool) GetConn() *net.IPConn {

	pool.lock.Lock()
	defer pool.lock.Unlock()
	if pool.alive <= 0 {
		return nil
	}
	var index = pool.count % pool.alive
	pool.count = index + 1
	pool.conn[index].count += 1
	return pool.conn[index].conn
}

//打印连接池发送数据情况
func (pool *PingPool) Print() {

	pool.lock.Lock()
	defer pool.lock.Unlock()
	var index = 0
	for index < pool.alive {
		fmt.Printf("连接%v:发包总计-%v\n", index, pool.conn[index].count)
		index++
	}
}

//产生连接池
func GeneratePool(src *net.IPAddr, conn_max_size int) *PingPool {
	var (
		index = 0
		err   error
		conn  *net.IPConn
		pool  *PingPool = &PingPool{lock: new(sync.Mutex), count: 0, alive: 0}
	)

	//windows不支持
	if runtime.GOOS != "linux" || !config.One {
		return pool
	}

	for index < conn_max_size {
		// windows不支持
		conn, err = net.ListenIP("ip4:icmp", src)
		if err != nil {
			//			fmt.Printf("监听ICMP异常：%v", err.Error())
			index += 1
			continue
		} else {

			conn.SetReadBuffer(1024 * 1024 * 12)
			conn.SetWriteBuffer(1024 * 1024 * 12)
			pool.conn = append(pool.conn, &PingConn{conn: conn, count: int64(0)})
			pool.alive += 1
		}
		index += 1
	}
	return pool

}
