// Package check @Author lanpang
// @Date 2024/9/6 下午2:28:00
// @Desc
package check

var hosts = make(map[string]*Host)

func HostCheck(baseUrl string) {
	// CPU 使用率
	getHostData(baseUrl, "cpu_usage_active")
	// 内存 使用率
	getHostData(baseUrl, "mem_used_percent")
	// 内存 大小
	getHostData(baseUrl, "mem_total")
	// 入网流量
	getHostData(baseUrl, "rate", "net_bytes_recv", "eth0")
	// 出网流量
	getHostData(baseUrl, "rate", "net_bytes_sent", "eth0")
	// 系统盘
	getHostData(baseUrl, "disk_used_percent", "/")
	// 数据盘
	getHostData(baseUrl, "disk_used_percent", "/data")
	//fmt.Println(hosts)
	// 输出表格
	tableRender(hosts)
}
