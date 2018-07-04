package pingbase

import (
	"net"
)

//整数返回IP
func IPFromInt(ip int64) net.IP {
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

//IP 转换成整数，仅支持IPv4
func IPInt(ip net.IP) int64 {
	ip_byte := ip.To4()
	return int64(ip_byte[0])<<24 + int64(ip_byte[1])<<16 + int64(ip_byte[2])<<8 + int64(ip_byte[3])
}

//IP添加整数转换成IP
func IPAdd(ip net.IP, offset int64) net.IP {
	return IPFromInt(IPInt(ip) + offset)
}
