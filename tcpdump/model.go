// Package tcpdump @Author lanpang
// @Date 2024/9/24 下午1:56:00
// @Desc
package tcpdump

import "net"

type Interface struct {
	Name        string
	Description string
	Flags       uint32
	Addresses   []InterfaceAddress
}

type InterfaceAddress struct {
	IP        net.IP
	Netmask   net.IPMask // Netmask may be nil if we were unable to retrieve it.
	Broadaddr net.IP     // Broadcast address for this IP may be nil
	P2P       net.IP     // P2P destination address for this IP may be nil
}
