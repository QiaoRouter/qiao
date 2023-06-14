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

			// 先构造一个 dhcp 的回复报文
			// 1. 生成 dhcp 实体
			// 2. 调用成员函数 -> toIpv6UdpPacket
			// 3. 内部构造 udp header -> 将 dhcp 序列化 -> 将 udp 序列化 -> 返回
			// 4. 内部构造 ipv6 header

			replyDhcp := DHCPv6Packet{}

			// 构造头部：msg-type + txn-id
			if dhcp.Header.MsgType == DHCP_SOLICIT {
				replyDhcp.Header.MsgType = DHCP_ADVERTISE
			} else if dhcp.Header.MsgType == DHCP_REQUEST {
				replyDhcp.Header.MsgType = DHCP_REPLY
			}
			replyDhcp.Header.TransactionIdHi = dhcp.Header.TransactionIdHi
			replyDhcp.Header.TransactionIdLo = dhcp.Header.TransactionIdLo
			MakeDhcpv6(&replyDhcp, clientIdPacket, IaNaPacket, hal.MacAddr(h.IfName))

			// 打包发送
			// src: eui64(localMac);
			// dst: ip6->ip6_src
			// @asu fixme: 这里要用 eui64 把本地 mac 解析成 ipv6 地址，作为 ipv6 header 的 src
			ipv6Dgrm := replyDhcp.ToIpv6UdpPacket(protocol.Ipv6Addr{}, dgrm.Header.Src)

			go hal.SendIpv6(h.IfIndex, ipv6Dgrm, ether.Header.SrcHost)
			return
		}
		if dgrm.Header.NextHeader == protocol.IPProtocolICMPV6 {
			return
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

func forwardPacket(h *hal.IfHandle, dgrm *protocol.Ipv6Datagram, ether *protocol.EthernetFrame) {
	ttl := dgrm.Header.HopLimit
	if ttl <= 1 {
		// 发送 ICMP Time Exceeded 消息
		// 将接受到的 IPv6 packet 附在 ICMPv6 头部之后。
		// 如果长度大于 1232 字节，则取前 1232 字节：
		// 1232 = IPv6 Minimum MTU(1280) - IPv6 Header(40) - ICMPv6 Header(8)
		// 意味着发送的 ICMP Time Exceeded packet 大小不大于 IPv6 Minimum MTU
		// 不会因为 MTU 问题被丢弃。
		// 详见 RFC 4443 Section 3.3 Time Exceeded Message
		// 计算 Checksum 后由自己的 IPv6 地址发送给源 IPv6 地址。
		var payloadLen int
		if ether.Payload.Length()+protocol.ICMPv6HeaderLen+protocol.Ipv6HeaderLen > protocol.Ipv6MinimumMTU {
			payloadLen = protocol.Ipv6MinimumMTU - protocol.ICMPv6HeaderLen - protocol.Ipv6HeaderLen
		} else {
			payloadLen = int(ether.Payload.Length())
		}
		buf := ether.Payload.Prefix(payloadLen)
		icmpv6 := protocol.MakeICMPv6Packet(protocol.ICMPv6TypeDestinationUnreachable,
			protocol.ICMPv6CodeDestinationNetworkUnreachable, buf)
		ipv6DgrmReply := icmpv6.ToIpv6Datagram(h.IPv6, dgrm.Header.Src, 255)
		hal.SendIpv6(h.IfIndex, ipv6DgrmReply, ether.Header.SrcHost)
	} else {
		e := PrefixQuery(dgrm.Header.Dst)
		if e != nil {
			nextIpAddr := e.Nexthop
			if nextIpAddr.AllZero() {
				nextIpAddr = dgrm.Header.Dst
			}

			dstMac, err := hal.GetNeighborMacAddr(e.IfIndex, nextIpAddr)
			if err != nil {
				fmt.Printf("h.GetNeighborMacAddr fail, err: %+v\n, nextIP: %+v\n",
					err, nextIpAddr.String())
				return
			}
			dgrm.Header.HopLimit--
			hal.SendIpv6(e.IfIndex, dgrm, dstMac)
		} else {
			// 没有找到路由
			// todo 回复icmpv6 route not found消息
			fmt.Printf("没有找到路由\n")
		}
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
	_, _ = parser.ParseNBytes(int(iaNa.Header.OptionLen - 4))
}

func ParseGargabe(parser *protocol.NetParser) {
	optionLen, _ := parser.ParseU16()
	_, _ = parser.ParseNBytes(int(optionLen))
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
	dhcp.Header.TransactionIdLo, err = parser.ParseU16()
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
