// Package cmd @Author lanpang
// @Date 2024/8/16 下午1:42:00
// @Desc
package cmd

import (
	"github.com/spf13/cobra"
	"vhagar/tcpdump"
)

var (
	device   string
	port     string
	host     string
	pcapFile string
)

var tcpdumpCmd = &cobra.Command{
	Use:   "tcpdump",
	Short: "抓包工具",
	Long:  `抓包工具`,
	Run: func(cmd *cobra.Command, args []string) {
		filter := getFilter(port, host)
		tcpdump.TcpDump(device, filter, pcapFile)
	},
}

func init() {
	rootCmd.AddCommand(tcpdumpCmd)
	tcpdumpCmd.Flags().StringVarP(&device, "device", "i", "en0", "网卡名")
	tcpdumpCmd.Flags().StringVarP(&port, "port", "p", "", "端口")
	tcpdumpCmd.Flags().StringVarP(&host, "host", "o", "", "主机")
	tcpdumpCmd.Flags().StringVarP(&pcapFile, "file", "f", "", "保存到pcap文件")
}

func getFilter(p, h string) string {
	filter := ""
	if p != "" {
		filter = "tcp and port " + p
	}
	if h != "" {
		if filter != "" {
			filter = filter + " and host " + h
		} else {
			filter = "host " + h
		}
	}
	return filter
}
