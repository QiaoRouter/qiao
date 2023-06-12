package dhcpv6

import (
	"fmt"
	"qiao/hal"
	"qiao/protocol"
)

// gate handle
func (e *Engine) receivePacketAndHandleIt(h *hal.IfHandle) {
	for true {
		ipv6Dgrm, ether, err := h.NextPacket()
		if err != nil {
			continue
		}
		go e.HandleIpv6(h, ipv6Dgrm, ether)
	}
}

// handle ipv6
func (e *Engine) HandleIpv6(h *hal.IfHandle, dgrm *protocol.Ipv6Datagram, ether *protocol.EthernetFrame) {
	dstIsMe := false
	for _, ifh := range hal.IfHandles {
		if dgrm.Header.Dst.Equals(ifh.IPv6) {
			dstIsMe = true
		}
	}
	multicastAddr, _ := protocol.ParseIpv6("ff02::9")
	if dgrm.Header.Dst.Equals(multicastAddr) {
		dstIsMe = true
	}
	if dstIsMe {
		// 处理 udp 报文
		if dgrm.Header.NextHeader == protocol.IPProtocolUdp {
			// 解析出来 udp -> 处理
			udp, err := protocol.ParseUdp(dgrm.Payload)
			if err != nil {
				fmt.Printf("protocol.ParseDhcpv6 fail, err: %+v\n", err)
			}
			// 检查端口号 ---> 必须是 server
			if udp.DstPort == DHCP_SERVER {
				fmt.Println("确实是DHCP Solicit 或者 DHCP request 报文")
			}
		}
		if dgrm.Header.NextHeader == protocol.IPProtocolICMPV6 {

		}
	} else {
		// forwarding
		// 目标地址不是我，考虑转发给下一跳
		// 检查是否是组播地址（ff00::/8），不需要转发组播分组
		if dgrm.Header.Dst.Octet[0] == 0xff {
			return
		}
		forwardPacket(h, dgrm, ether)
	}
}
