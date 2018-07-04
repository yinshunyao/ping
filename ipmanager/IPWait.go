/*等待结果的IP字典，按照IP的字符串形式保存*/
package ipmanager

import (
	"Ping/config"
	"sync"
)

type Wait struct {
	lock *sync.Mutex
	dict map[string]*IPPing
}

//添加入队
func (wait *Wait) Add(ipping *IPPing) {
	//	wait.lock.Lock()
	//	defer wait.lock.Unlock()
	//不存在的时候
	if _, ok := wait.dict[ipping.String()]; !ok {
		wait.dict[ipping.String()] = ipping
	}

}

func (wait *Wait) AddSecurity(ipping *IPPing) {
	wait.lock.Lock()
	defer wait.lock.Unlock()
	//不存在的时候
	if _, ok := wait.dict[ipping.String()]; !ok {
		wait.dict[ipping.String()] = ipping
	}

}

//添加入队
func (wait *Wait) Rmv(addr string) *IPPing {
	//	wait.lock.Lock()
	//	defer wait.lock.Unlock()
	var (
		ipping *IPPing
		ok     bool
	)

	//如果处于等待区，删除
	if ipping, ok = wait.dict[addr]; ok {
		delete(wait.dict, addr)
		return ipping
	} else {
		return nil
	}

}

func (wait *Wait) RmvSecurity(addr string) *IPPing {
	wait.lock.Lock()
	defer wait.lock.Unlock()
	var (
		ipping *IPPing
		ok     bool
	)

	//如果处于等待区，删除
	if ipping, ok = wait.dict[addr]; ok {
		delete(wait.dict, addr)
		return ipping
	} else {
		return nil
	}

}

/*失败
整数返回  0表示不需要做任何处理  1表示需要重新Ping  -1表示结束，不需要再Ping了
*/
func (wait *Wait) Fail(addr string) (*IPPing, int) {
	wait.lock.Lock()
	defer wait.lock.Unlock()
	var (
		ipping *IPPing
		ok     bool
	)
	if ipping, ok = wait.dict[addr]; ok {
		//当失败，且不需要再次重试时，删除返回
		if ipping.try >= config.Retry {
			delete(wait.dict, addr)
			return ipping, -1
		} else { // 失败，还需要重试，不用删除，返回ipping
			return ipping, 1
		}

	}
	return nil, 0
}

// 成功
func (wait *Wait) Success(addr string) *IPPing {
	wait.lock.Lock()
	defer wait.lock.Unlock()
	var (
		ipping *IPPing
		ok     bool
	)

	//如果处于等待区，删除
	if ipping, ok = wait.dict[addr]; ok {
		ipping.Success()
		delete(wait.dict, addr)
	} else {
		return nil
	}
	return ipping
}

func GenerateIPWait() *Wait {
	return &Wait{lock: new(sync.Mutex), dict: make(map[string]*IPPing)}
}
