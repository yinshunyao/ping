package pingbase

//	"bytes"
//	"encoding/binary"
//	"fmt"
//	"net"
//"os"
//"strconv"
//	"time"

// Ping包包结构
type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequencyNum uint16
}

//校验Ping包
func CheckSum(data []byte) (rt uint16) {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		//前8bit + 后8bit
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}

	//低8位相加
	if length == 1 {
		sum += uint32(data[index]) // << 8
	}
	rt = uint16(sum) + uint16(sum>>16)
	return ^rt
}

//func PingDomain(domain string, PS, Count int) {
//	var (
//		icmp ICMP
//		//		laddr                     = src}
//		raddr, _                  = net.ResolveIPAddr("ip", domain)
//		max_lan, min_lan, avg_lan float64
//	)

//	conn, err := net.DialIP("ip4:icmp", &src, raddr)
//	if err != nil {
//		fmt.Printf(err.Error())
//		return
//	}
//	defer conn.Close()
//	icmp.Type = 8
//	icmp.Code = 0
//	icmp.Checksum = 0
//	icmp.Identifier = 0
//	icmp.SequencyNum = 0

//	var buffer bytes.Buffer
//	//写入icmp包头
//	binary.Write(&buffer, binary.BigEndian, icmp)
//	//写入icmp包内容
//	binary.Write(&buffer, binary.BigEndian, originaBytes)
//	b := buffer.Bytes()
//	binary.BigEndian.PutUint16(b[2:], CheckSum(b))
//	if debug {
//		fmt.Printf("\n正在Ping %s 具有 %d(%d)字节的数据:\n", raddr.String(), PS, PS+28)
//	}

//	recv := make([]byte, 1024)
//	ret_list := []float64{}
//	dropPack := 0.0
//	max_lan = 3000.0
//	min_lan = 0.0
//	avg_lan = 0.0
//	for i := Count; i > 0; i-- {
//		if _, err := conn.Write(buffer.Bytes()); err != nil {
//			dropPack++
//			time.Sleep(time.Second)
//			continue
//		}
//		t_start := time.Now()
//		conn.SetReadDeadline((time.Now().Add(time.Second * 3)))
//		_, err := conn.Read(recv)
//		if err != nil {
//			dropPack++
//			time.Sleep(time.Second)
//			continue
//		}
//		t_end := time.Now()
//		dur := float64(t_end.Sub(t_start).Nanoseconds()) / 1e6
//		ret_list = append(ret_list, dur)
//		if dur < max_lan {
//			max_lan = dur
//		}
//		if dur < min_lan {
//			min_lan = dur
//		}
//		fmt.Printf("来自 %s 的回复：时间 = %.3fms\n", raddr.String(), dur)
//		time.Sleep(time.Second)
//	}

//	fmt.Printf("丢包率：%.2f%%\n", dropPack/float64(Count))
//	fmt.Printf("rtt min/avg/max = %.3fms/.3fms/%.3fms\n", min_lan, avg_lan, max_lan)
//}
