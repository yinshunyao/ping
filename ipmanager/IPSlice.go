/*
IP数组，无序，基本上不需要按照关键字查询，只是添加和删除
可以用作存储 初始化所有需要Ping的IP对象  以及 Ping结束的所有对象
*/
package ipmanager

import (
	"fmt"
	"sync"
)

//IP列表
type IPSlice struct {
	lock *sync.RWMutex
	ips  []*IPPing
}

//长度获取
func (ipslice *IPSlice) Length() int64 {
	ipslice.lock.RLock()
	defer ipslice.lock.RUnlock()
	if ipslice.ips == nil {
		return 0
	}
	return int64(len(ipslice.ips))
}

//长度获取
func (ipslice *IPSlice) AliveList(print_flag bool) []string {
	ipslice.lock.RLock()
	defer ipslice.lock.RUnlock()
	var result []string
	if ipslice.ips == nil {
		return nil
	}
	index := int64(0)
	for index < int64(len(ipslice.ips)) {
		if ipslice.ips[index].alive && print_flag {
			if print_flag {
				fmt.Printf("%v\n", ipslice.ips[index].ip.String())
			}
			result = append(result, ipslice.ips[index].ip.String())
		}

		index++

	}
	return result
}

// 安全添加
func (ipslice *IPSlice) AddSecurity(ipping *IPPing) {
	ipslice.lock.Lock()
	defer ipslice.lock.Unlock()
	ipslice.ips = append(ipslice.ips, ipping)
}

// 非安全模式添加
func (ipslice *IPSlice) Add(ipping *IPPing) {
	ipslice.ips = append(ipslice.ips, ipping)
}

// 安全的获取某个节点
func (ipslice *IPSlice) PopSecurity(index int64) *IPPing {
	ipslice.lock.Lock()
	defer ipslice.lock.Unlock()
	length := ipslice.Length()

	var ipping *IPPing
	if index >= length || index < 0 {
		return nil
	}

	//删除最后一个元素
	if length == 1 {
		ipping = ipslice.ips[0]
		ipslice.ips = nil
	} else {
		//将待删除的元素替换成最后一个元素
		ipping = ipslice.ips[index]
		ipslice.ips[index] = ipslice.ips[length-1]
		//删除末尾的元素
		ipslice.ips = ipslice.ips[0 : length-1]
	}

	return ipping
}

//获取某个节点，非线程安全模式
func (ipslice *IPSlice) Pop(index int64) *IPPing {
	var ipping *IPPing
	length := ipslice.Length()
	if index >= length || index < 0 {
		return nil
	}

	//删除最后一个元素
	if length == 1 {
		ipping = ipslice.ips[0]
		ipslice.ips = nil
	} else {
		//将待删除的元素替换成最后一个元素
		ipping = ipslice.ips[index]
		ipslice.ips[index] = ipslice.ips[length-1]
		//删除末尾的元素
		ipslice.ips = ipslice.ips[0 : length-1]
	}

	return ipping
}

func GenerateIPSlice() *IPSlice {
	return &IPSlice{lock: new(sync.RWMutex)}
	//	return &IPSlice{lock: new(sync.Mutex), ips: make([1]*IPPing)}
}
