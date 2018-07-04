# ping
go语言版本，支持IP段Ping，支持自定义速率，当前速率控制不太精准


支持多IP段和单个IP输入
./Ping 8.7.0.0-8.8.255.255   8.3.3.0  -d true  -r 5000 -t 4 


	-d,--debug	bool	调试模式，默认false
	-t,--timeout	int	超时定时器，单位秒
	-r,--rate	int	发包速率，定义每秒的最大发包个数
	-n,--number	int	每个IP最多重试次数
	-s,--src	string	源IP
	-o,--one	bool	linux下是否启用单socket

