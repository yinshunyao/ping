//main
package main

import (
	"Ping/config"
	"Ping/pingcontrol"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var err error
var control *pingcontrol.Control = pingcontrol.GenerateControl()

func print_help() {
	fmt.Println("[输入-h查看帮助]")
	fmt.Println("参数：%s", os.Args)
	fmt.Println("\t-d,--debug\tbool\t调试模式，默认false")
	fmt.Println("\t-t,--timeout\tint\t超时定时器，单位秒")
	fmt.Println("\t-r,--rate\tint\t发包速率，定义每秒的最大发包个数")
	fmt.Println("\t-n,--number\tint\t每个IP最多重试次数")
	fmt.Println("\t-s,--src\tstring\t源IP")
	fmt.Println("\t-o,--one\tbool\tlinux下是否启用单socket")
}

func main() {
	length := len(os.Args)
	if length <= 0 {
		print_help()
	}

	var index = 0
	var command string
	for index < length {
		command = os.Args[index]
		//超时定时器参数读取
		if command == "-t" || command == "--timeout" {
			config.Timeout, err = strconv.Atoi(os.Args[index+1])
			if err != nil {
				fmt.Println("定时器参数必须是整数，输入：%v", os.Args[index+1])
				print_help()
				return
			}
		} else if command == "-n" || command == "--number" {
			config.Retry, err = strconv.Atoi(os.Args[index+1])
			if err != nil {
				fmt.Println("重试次数参数必须是整数，输入：%v", os.Args[index+1])
				print_help()
				return
			}

		} else if command == "-d" || command == "--debug" {
			config.Debug, err = strconv.ParseBool(os.Args[index+1])
			if err != nil {
				fmt.Println("debug参数必须是bool，输入：%v", os.Args[index+1])
				print_help()
				return
			}

		} else if command == "-o" || command == "--one" {
			config.One, err = strconv.ParseBool(os.Args[index+1])
			if err != nil {
				fmt.Println("one参数必须是bool，输入：%v", os.Args[index+1])
				print_help()
				return
			}

		} else if command == "-r" || command == "--rate" {
			config.Rate, err = strconv.Atoi(os.Args[index+1])
			if err != nil {
				fmt.Println("发包速率参数必须是整数，输入：%v", os.Args[index+1])
				print_help()
				return
			}

		} else if command == "-s" || command == "--src" {
			config.Src = net.IPAddr{IP: net.ParseIP(os.Args[index+1])}
			//			if err != nil {
			//				fmt.Println("源IP输入错误：%v", os.Args[index+1])
			//				print_help()
			//				return
			//			}

		} else {
			//尝试按照IP解析，如果报错暂时不理会
			//			IPParse(os.Args[index])
			//			ipmanager.InitQueue(os.Args[index])
			control.Parse(os.Args[index])
			index += 1
			continue
		}

		index += 2
	}

	//获取入参端口
	//	v := flag.String("v", "1.0.0", "版本2.0.0")
	//	ips := flag.String("ip", "127.0.0.1", "IP列表，格式是单个IP 8.8.8.8或者IP范围8.8.8.0-8.8.8.255")
	//	Ping(ip, 32, count_retry)
	t_start := time.Now()
	//	fmt.Println("开始扫描")
	//先监听
	//	if !control.Run() {
	//		return
	//	}
	control.Run()
	t_end := time.Now()
	fmt.Printf("ping %v ip finish, cost: %v, result is:\n", control.FinishCount(), t_end.Sub(t_start))
	control.Print()
}
