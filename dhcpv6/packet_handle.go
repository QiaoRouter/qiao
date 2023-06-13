package dhcpv6

import (
	"fmt"
	"log"
	dhcpv6 "qiao/dhcpv6/protocol"
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

			// 解析出来 udp 报文中的 dhcp
			dhcp, err := ParseDhcpPacket(udp.Payload)

			// 检查是否为 DHCPv6 Solicit 或 DHCPv6 Request
			if dhcp.Header.MsgType == DHCP_SOLICIT || dhcp.Header.MsgType == DHCP_ADVERTISE {
				fmt.Println("确实是个 dhcp solicit 或者是 dhcp request 报文")
			}

			// 解析 DHCPv6 头部后的 Option，找到其中的 Client Identifier
			// 和 IA_NA 中的 IAID
			// https://www.rfc-editor.org/rfc/rfc8415.html#section-21.2
			// https://www.rfc-editor.org/rfc/rfc8415.html#section-21.4

			// 解析出来两个 Option 包裹
			clientIdPacket, IaNaPacket, _ := ParseDhcpOptions(dhcp.MessageBody)

			// 构造响应对方的 ipv6 报文

			return
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
		// forwardPacket(h, dgrm, ether)
	}
}

func ParseDhcpOptions(buf protocol.Buffer) (*dhcpv6.OptionClientId, *dhcpv6.OptionIaNa, error) {
	parser := protocol.NetParser{
		Buffer:  buf,
		Pointer: 0,
	}

	// 是否两个都已经找到
	flagClientIdentifier, flagIaNa := false, false
	clientIdentifier, iaNa := &dhcpv6.OptionClientId{}, &dhcpv6.OptionIaNa{}

	// 专门儿找这两个
	i := 1
	for !flagClientIdentifier || !flagIaNa {
		fmt.Printf("目前是第 %v 个\n", i)
		i++

		// 先拿出来编号看看
		optionCode, _ := parser.ParseU16()
		if optionCode == DHCP_OPTION_CLIENTID {
			fmt.Println("找着 client id 了")

			// 获取 client id 完整数据
			clientIdentifier.Header.OptionCode = optionCode
			ParseClientId(&parser, clientIdentifier)

			flagClientIdentifier = true
		} else if optionCode == DHCP_OPTION_IA_NA {
			fmt.Println("找着 ia-na 了")

			// 获取 client id 完整数据
			iaNa.Header.OptionCode = optionCode
			ParseIaNa(&parser, iaNa)

			flagIaNa = true
		} else {
			// 跳过非目标 option
			ParseGargabe(&parser)
		}
	}

	return clientIdentifier, iaNa, nil
}

func ParseClientId(parser *protocol.NetParser, clientIdentifier *dhcpv6.OptionClientId) {
	// 解析长度
	clientIdentifier.Header.OptionLen, _ = parser.ParseU16()

	// 解析 uuid
	log.Printf("formor length of uuid: %v\n", len(clientIdentifier.DUID))
	clientIdentifier.DUID, _ = parser.ParseNBytes(int(clientIdentifier.Header.OptionLen))
	log.Printf("current length of uuid after appending: %v\n", len(clientIdentifier.DUID))
}

func ParseIaNa(parser *protocol.NetParser, iaNa *dhcpv6.OptionIaNa) {
	// 解析长度
	iaNa.Header.OptionLen, _ = parser.ParseU16()

	// 解析 IA-ID
	iaNa.IaId, _ = parser.ParseU32()

	// 跳过
	// skip T1, T2, IA_NA Options
	// len = IAID + T1 + T2 + IANA Option
	// skip size = len - IAID (already resolved)
	_ = parser.ParseSkipNBytes(int(iaNa.Header.OptionLen - 4))
}

func ParseGargabe(parser *protocol.NetParser) {
	optionLen, _ := parser.ParseU16()
	_ = parser.ParseSkipNBytes(int(optionLen))
}

// 解析 dhcpv6 报文头部
func ParseDhcpPacket(buf protocol.Buffer) (*DHCPv6Packet, error) {
	var err error
	parser := protocol.NetParser{
		Buffer:  buf,
		Pointer: 0,
	}
	dhcp := &DHCPv6Packet{}

	// dhcp 类型
	dhcp.Header.MsgType, err = parser.ParseU8()
	if err != nil {
		return nil, err
	}

	// txn id 高位
	dhcp.Header.TransactionIdHi, err = parser.ParseU8()
	if err != nil {
		return nil, err
	}

	// txn id 低位
	dhcp.Header.TransactionIdLo, err = parser.ParseU8()
	if err != nil {
		return nil, err
	}

	// 报文内容
	dhcp.MessageBody, err = parser.ParseBuffer()
	if err != nil {
		return nil, err
	}

	return dhcp, err
}
