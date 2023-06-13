package dhcpv6

import (
	dhcpv6 "qiao/dhcpv6/protocol"
	"qiao/protocol"
)

type DHCPv6Packet struct {
	Header      DHCPv6Header
	MessageBody protocol.Buffer
}

type DHCPv6Header struct {
	MsgType         uint8
	TransactionIdHi uint8
	TransactionIdLo uint16
}

// 构造 dhcp 报文
func MakeDhcpv6(replyDhcp *DHCPv6Packet, clientIdPacket *dhcpv6.OptionClientId, IaNaPacket *dhcpv6.OptionIaNa, localMac protocol.EthernetAddr) {
	// 构造响应对方的 dhcp 报文

	// 构造 options
	// @asu fixme: 这个计算长度原本在 cpp 中用于确认该报文的长度，填入 udp 中的 len
	// @asu fixme: 但是在这里似乎可以直接用 len 获取这个 Buffer 的长度，但还是先放在这里了
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

	// 2. 计算长度
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

	// 3. 头部
	replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, DHCP_OPTION_IA_NA)
	replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, DHCP_OPTIONS_IANA_LEN+4+DHCP_OPTIONS_IANA_OPTION_LEN)

	// 3. 身体
	replyDhcp.MessageBody.Octet = protocol.ConcatU32(replyDhcp.MessageBody.Octet, IaNaPacket.IaId)
	replyDhcp.MessageBody.Octet = protocol.ConcatU32(replyDhcp.MessageBody.Octet, 0)
	replyDhcp.MessageBody.Octet = protocol.ConcatU32(replyDhcp.MessageBody.Octet, 0)

	// 3. 计算长度
	curLen = DHCP_OPTIONS_HDR_LEN + DHCP_OPTIONS_IANA_LEN // 4 + 12
	replyDhcpLen += curLen

	// 3. sub 头部
	replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, DHCP_OPTION_IAADDR)
	replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, DHCP_OPTIONS_IANA_OPTION_LEN)

	// 3. sub 身体
	// 16 + 4 + 4 = 24
	replyDhcp.MessageBody.Octet = protocol.ConcatIpv6Addr(replyDhcp.MessageBody.Octet, &DefaultGateway)
	replyDhcp.MessageBody.Octet = protocol.ConcatU32(replyDhcp.MessageBody.Octet, Preferredlifetime)
	replyDhcp.MessageBody.Octet = protocol.ConcatU32(replyDhcp.MessageBody.Octet, Validlifetime)

	// 3. sub 计算长度
	curLen = DHCP_OPTIONS_HDR_LEN + DHCP_OPTIONS_IANA_OPTION_LEN // 4 + 12
	replyDhcpLen += curLen

	// 4. DNS recursive name server：包括两个 DNS 服务器地址
	// 2402:f000:1:801::8:28 和 2402:f000:1:801::8:29。
	//    - https://www.rfc-editor.org/rfc/rfc3646#section-3
	//    - Option Code: 23
	//    - Option Length: 32
	//    - DNS: 2402:f000:1:801::8:28
	//    - DNS: 2402:f000:1:801::8:29

	// 4. 头部
	replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, DHCP_OPTION_DNS_SERVERS)
	replyDhcp.MessageBody.Octet = protocol.ConcatU16(replyDhcp.MessageBody.Octet, 16<<3) // 16 * 8

	// 4. 身体
	replyDhcp.MessageBody.Octet = protocol.ConcatIpv6Addr(replyDhcp.MessageBody.Octet, &DnsA)
	replyDhcp.MessageBody.Octet = protocol.ConcatIpv6Addr(replyDhcp.MessageBody.Octet, &DnsB)

	// 4. 计算长度
	curLen = DHCP_OPTIONS_HDR_LEN + (protocol.IPv6AddrLen << 1) // size of two ipv6 address
	replyDhcpLen += curLen
}

// 将 dhcp 报文序列化成 byte stream
func (packet *DHCPv6Packet) Serialize() protocol.Buffer {
	var s []byte
	s = protocol.ConcatU8(s, packet.Header.MsgType)
	s = protocol.ConcatU8(s, packet.Header.TransactionIdHi)
	s = protocol.ConcatU16(s, packet.Header.TransactionIdLo)
	s = protocol.ConcatBuffer(s, packet.MessageBody)

	return protocol.Buffer{
		Octet: s,
	}
}

// 构造 udp 报文，将 hdcp 打包进去，加上 udp 头部
func (packet *DHCPv6Packet) ToUdpPacket() *protocol.UdpPacket {
	buf := packet.Serialize()
	return &protocol.UdpPacket{
		DstPort: DHCP_CLIENT,
		SrcPort: DHCP_SERVER,
		Len:     buf.Length() + protocol.UdpHeaderLen,
		Payload: buf,
	}
}

func (packet *DHCPv6Packet) ToIpv6UdpPacket(src protocol.Ipv6Addr, dst protocol.Ipv6Addr) *protocol.Ipv6Datagram {
	udp := packet.ToUdpPacket()
	payload := udp.Serialize()
	ipv6Dgrm := &protocol.Ipv6Datagram{
		Header: protocol.Ipv6Header{
			Version:    6,
			PayloadLen: payload.Length(),
			FlowLabel:  0,
			NextHeader: protocol.IPProtocolUdp,
			HopLimit:   255,
			Src:        src,
			Dst:        dst,
		},
		Payload: payload,
	}
	return ipv6Dgrm
}
