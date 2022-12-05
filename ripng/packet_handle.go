package ripng

import (
	"fmt"
	"qiao/hal"
	"qiao/protocol"
)

func (e *Engine) receivePacketAndHandleIt(h *hal.IfHandle) {
	for true {
		ipv6Dgrm, ether, err := h.NextPacket()
		if err != nil {
			continue
		}
		go e.HandleIpv6(h, ipv6Dgrm, ether)
	}
}

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
		if dgrm.Header.NextHeader == protocol.IPProtocolUdp {
			ripngPacket, err := ParseRipngPacket(dgrm)
			if err != nil {
				return
			}
			for _, e := range ripngPacket.Entries {
				rte := RouteTableEntry{
					Ipv6Addr: e.Prefix,
					Len:      int(e.PrefixLen),
					IfIndex:  h.IfIndex,
					Nexthop:  dgrm.Header.Src,
					Metric:   int(e.Metric + 1),
				}
				if rte.Metric == 0xf {
					continue
				}
				oldE := ExactQuery(e.Prefix, int(e.PrefixLen))
				update := func() {
					if err = Update(&rte); err != nil {
						fmt.Printf("Update(&rte) fail, err: %+v\n", err)
					}
				}
				if oldE != nil {
					if oldE.Nexthop == rte.Nexthop && oldE.IfIndex == rte.IfIndex {
						update()
					} else {
						if rte.Metric < oldE.Metric {
							update()
						}
					}
				} else {
					update()
				}
			}
			return
		}
		if dgrm.Header.NextHeader == protocol.IPProtocolICMPV6 {
			icmp, err := protocol.ParseICMPv6(dgrm.Payload)
			if err != nil {
				fmt.Printf("protocol.ParseIcmpv6 fail, err: %+v\n", err)
			}
			if icmp.Header.Type == protocol.ICMPv6TypeEchoRequest {
				icmp.Header.Type = protocol.ICMPv6TypeEchoReply
				icmp.Header.Checksum = 0
				ipv6Reply := icmp.ToIpv6Datagram(dgrm.Header.Dst, dgrm.Header.Src, 64)
				hal.SendIpv6(h.IfIndex, ipv6Reply, ether.Header.SrcHost)
			}
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

func ParseRipngPacket(dgrm *protocol.Ipv6Datagram) (*RipngPacket, error) {
	if dgrm.Payload.Length() < protocol.UdpHeaderLen {
		return nil, ErrLength
	}
	if dgrm.Header.PayloadLen != dgrm.Payload.Length() {
		return nil, ErrLength
	}
	udp, err := protocol.ParseUdp(dgrm.Payload)
	if err != nil {
		return nil, err
	}
	if udp.DstPort != RipngUdpPort || udp.SrcPort != RipngUdpPort {
		return nil, ErrPortNotRipng
	}
	ripngEntriesLen := udp.Len - protocol.UdpHeaderLen - RipngHeaderLength
	if ripngEntriesLen%RipngEntriesLength != 0 {
		return nil, ErrLength
	}
	ripngPacket := &RipngPacket{}
	parser := protocol.NetParser{
		Buffer: udp.Payload,
	}
	ripngPacket.Command, err = parser.ParseU8()
	if err != nil {
		return nil, err
	}
	version, err := parser.ParseU8()
	if version != RipngVersion {
		return nil, ErrVersion
	}
	mustBeZero, err := parser.ParseU16()
	if err != nil {
		return nil, err
	}
	if mustBeZero != 0 {
		return nil, ErrBadZero
	}
	ripngPacket.NumEntries = uint32(ripngEntriesLen / RipngEntriesLength)
	for i := 0; i < int(ripngPacket.NumEntries); i++ {
		e := &RipngRte{}
		e.Prefix, err = parser.ParseIpv6Addr()
		if err != nil {
			return nil, err
		}
		e.RouteTag, err = parser.ParseU16()
		if err != nil {
			return nil, err
		}
		e.PrefixLen, err = parser.ParseU8()
		if err != nil {
			return nil, err
		}
		e.Metric, err = parser.ParseU8()
		if err != nil {
			return nil, err
		}
		if e.Metric >= 0xf {
			continue
		}
		if e.PrefixLen > 128 {
			return nil, ErrBadPrefixLen
		}
		routeAddr := e.Prefix.ToRouteAddr(int(e.PrefixLen))
		if !routeAddr.Equals(e.Prefix) {
			return nil, ErrInconsistentPrefix
		}
		ripngPacket.Entries = append(ripngPacket.Entries, e)
	}

	return ripngPacket, nil
}
