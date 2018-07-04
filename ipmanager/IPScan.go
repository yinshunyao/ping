/*IP扫描队列，将需要扫描的IP放入队列；可以从队列中获取需要Ping的对象*/
package ipmanager

import (
	"Ping/config"
	queue "Ping/queue"
	"fmt"
	"time"
)

type Scan struct {
	q *queue.Queue
}

func (scan *Scan) Put(ipping *IPPing) {
	var err_put error
	for {
		err_put = scan.q.Put(ipping, TIMEOUT_QUEUE)
		//入队成功，退出
		if err_put == nil {
			if config.Debug {
				//				fmt.Printf("%v入队成功\n", ipping.String())
			}
			break
		}

		if config.Debug {
			fmt.Printf("%v入队失败: %v", ipping.String(), err_put.Error())
		}
		//入队失败，等待20ms再次入队
		time.Sleep(MS_20)
	}
}

func (scan *Scan) Get() *IPPing {
	val, err_get := scan.q.Get(TIMEOUT_QUEUE)
	//获取成功，退出
	if err_get != nil {
		//		if config.Debug {
		//			fmt.Printf("获取IP异常：%v\n", err_get.Error())
		//		}
		return nil
	}
	return val.(*IPPing)
}

func (scan *Scan) Size() int {
	return scan.q.Size()
}

func GenerateIPScan() *Scan {
	return &Scan{q: queue.New(MAX_QUEUE_SIZE)}
}
