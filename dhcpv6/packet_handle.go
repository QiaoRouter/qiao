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

			// 构造响应对方的 ipv6 报文
			// 构造头部：msg-type + txn-id
			if dhcp.Header.MsgType == DHCP_SOLICIT {
				replyDhcp.Header.MsgType = DHCP_ADVERTISE
			} else if dhcp.Header.MsgType == DHCP_REQUEST {
				replyDhcp.Header.MsgType = DHCP_REPLY
			}

			replyDhcp.Header.TransactionIdHi = dhcp.Header.TransactionIdHi
			replyDhcp.Header.TransactionIdLo = dhcp.Header.TransactionIdLo

			// 构造 options
			replyDhcpLen, curLen := 0, 0

			// 1. Server Identifier：根据本路由器在本接口上的 MAC 地址生成。
			//    - https://www.rfc-editor.org/rfc/rfc8415.html#section-21.3
			//    - Option Code: 2
			//    - Option Length: 14
			//    - DUID Type: 1 (Link-layer address plus time) -> 16 bits
			//    - Hardware Type: 1 (Ethernet) -> 16 bits
			//    - DUID Time: 0 -> 32 bits
			//    - Link layer address: MAC Address -> 46 bits = 6B

			// 1. Server Identifier 头部
			replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, DHCP_OPTION_SERVERID)
			replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, DHCP_OPTIONS_DUID_LEN)

			// 1. Server Identifier 身体
			replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, 1) // DUID Type: 1 (Link-layer address plus time)
			replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, 1) // Hardware Type: 1 (Ethernet)
			replyDhcp.MessageBody.Octet = protocol.ConcatU32(replyDhcp.MessageBody.Octet, 0) // DUID Time: 0 -> 32 bits
			localMac := hal.MacAddr(h.IfName)
			replyDhcp.MessageBody.Octet = protocol.ConcatMac(replyDhcp.MessageBody.Octet, &localMac)

			// 1. 计算长度
			curLen = DHCP_OPTIONS_HDR_LEN + DHCP_OPTIONS_DUID_LEN
			replyDhcpLen += curLen

			// 2. Client Identifier
			//    - https://www.rfc-editor.org/rfc/rfc8415.html#section-21.2
			//    - Option Code: 1
			//    - Option Length: 和 Solicit/Request 中的 Client Identifier
			//    一致
			//    - DUID: 和 Solicit/Request 中的 Client Identifier 一致

			// 2. 头部
			replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, DHCP_OPTION_CLIENTID)
			replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, clientIdPacket.Header.OptionLen)

			// 2. 身体 -> 只需要一个 uuid
			replyDhcp.MessageBody.Octet = protocol.ConcatNBytes(replyDhcp.MessageBody.Octet, int(clientIdPacket.Header.OptionLen), clientIdPacket.DUID)

			// 3. 计算长度
			curLen = DHCP_OPTIONS_HDR_LEN + DHCP_OPTIONS_DUID_LEN
			replyDhcpLen += curLen

			// 3. Identity Association for Non-temporary
			// Address：记录服务器将会分配给客户端的 IPv6 地址。
			//    - https://www.rfc-editor.org/rfc/rfc8415.html#section-21.4
			//    - Option Code: 3
			//    - Option Length: 40

			//    - IAID: 和 Solicit/Request 中的 Identity Association for
			//    Non-temporary Address 一致
			//    - T1: 0
			//    - T2: 0
			//    - IA_NA options:
			//      - https://www.rfc-editor.org/rfc/rfc8415.html#section-21.6
			//      - Option code: 5 (IA address)
			//      - Length: 24
			//      - IPv6 Address: fd00::1:2
			//      - Preferred lifetime: 54000s
			//      - Valid lifetime: 86400s

			// 1. 头部

			// 2. 身体

			// 3. 计算长度

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
