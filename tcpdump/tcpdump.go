// Package tcpdump @Author lanpang
// @Date 2024/9/24 下午12:19:00
// @Desc
package tcpdump

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"log"
	"os"
	"time"
)

var (
	snapshot int32 = 65535 //读取一个数据包的最大值，一般设置成这65535即可
	//promisc        = false            //是否开启混杂模式
	timeout = time.Second * 10 //抓取数据包的超时时间，负数表示立即刷新，一般都设为负数
	writer  *pcapgo.Writer
)

func TcpDump(device, filter, pcapFile string) {
	handle, err := pcap.OpenLive(device, snapshot, false, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// 设置过滤器
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Fatal(err)
	}
	fmt.Println("开始抓包...")

	// 创建一个 .pcap 文件用于保存捕获的数据包
	writer = newWriter(pcapFile)

	// 初始化解码器层和解析器
	var ethLayer layers.Ethernet
	var ipLayer layers.IPv4
	var tcpLayer layers.TCP
	var payload gopacket.Payload

	// 创建解码器解析器
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &ethLayer, &ipLayer, &tcpLayer, &payload)
	var decodedLayers []gopacket.LayerType

	// 实时分析的变量
	var packetCount int
	var totalBytes int

	// 开始捕获包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	startTime := time.Now()

	for packet := range packetSource.Packets() {
		// 写入 .pcap 文件
		if writer != nil {
			if err := writer.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
				log.Println("写入文件失败:", err)
			}
		}

		// 解析数据包
		err := parser.DecodeLayers(packet.Data(), &decodedLayers)
		if err != nil {
			//log.Println("解码失败:", err)
			continue
		}

		// 分析解析出的数据层
		for _, layerType := range decodedLayers {
			switch layerType {
			case layers.LayerTypeIPv4:
				fmt.Printf("IP层 - 源IP: %s, 目标IP: %s\n", ipLayer.SrcIP, ipLayer.DstIP)
			case layers.LayerTypeTCP:
				fmt.Printf("TCP层 - 源端口: %d, 目标端口: %d\n", tcpLayer.SrcPort, tcpLayer.DstPort)
			case layers.LayerTypeEthernet:
				//fmt.Printf("以太网层 - 源MAC: %s, 目标MAC: %s\n", ethLayer.SrcMAC, ethLayer.DstMAC)
			default:
				//fmt.Println("未知层类型:", layerType)
			}
		}

		//fmt.Println(packet)

		//解析包并打印
		//processPacket(packet)
		//
		//实时分析
		packetCount++
		totalBytes += len(packet.Data())

		// 每隔 10 秒打印一次实时统计
		if time.Since(startTime) > 10*time.Second {
			printStats(packetCount, totalBytes, time.Since(startTime))
			startTime = time.Now()
			packetCount = 0
			totalBytes = 0
		}
	}
}

//func packetStatistic() {
//}

// 分析并处理包
func processPacket(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		fmt.Printf("从 %s 到 %s\n", ip.SrcIP, ip.DstIP)
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		fmt.Printf("源端口: %d, 目标端口: %d\n", tcp.SrcPort, tcp.DstPort)
	}

	// 打印数据包的详细信息
	//fmt.Println(packet)
}

// 打印实时统计数据
func printStats(packetCount int, totalBytes int, duration time.Duration) {
	if duration.Seconds() == 0 {
		return
	}
	// 计算流量速率（bps）
	avgRate := float64(totalBytes*8) / duration.Seconds()
	fmt.Printf("\n=== 实时统计 ===\n")
	fmt.Printf("捕获包数: %d\n", packetCount)
	fmt.Printf("总字节数: %d bytes\n", totalBytes)
	fmt.Printf("平均速率: %.2f bps\n", avgRate)
	fmt.Println("================")
}

func getDevices() {
	devices, _ := pcap.FindAllDevs()
	for _, device := range devices {
		fmt.Println("\nname:", device.Name)
		fmt.Println("describe:", device.Description)
		for _, address := range device.Addresses {
			fmt.Println("IP:", address.IP)
			fmt.Println("mask:", address.Netmask)
		}
	}
}

func newWriter(fileName string) *pcapgo.Writer {
	if fileName == "" {
		return nil
	}
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)
	w := pcapgo.NewWriter(f)
	err = w.WriteFileHeader(uint32(snapshot), layers.LinkTypeEthernet)
	if err != nil {
		return nil
	}
	return w
}
