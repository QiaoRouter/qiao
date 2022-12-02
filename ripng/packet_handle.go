package ripng

import (
	"fmt"
	"qiao/hal"
	"qiao/protocol"
)

func (e *Engine) receivePacketAndHandleIt(h *hal.IfHandle) {
	for true {
		p, err := h.PacketSource.NextPacket()
		if err != nil {
			continue
		}
		ether, err := protocol.ParseEtherFrame(p.Data())
		if err != nil {
			continue
		}
		if ether.Header.SrcHost.Equals(h.MAC) {
			continue
		}
		if ether.Header.Type != protocol.EthernetProtocolIPv6 {
			continue
		}
		ipv6Dgrm, err := protocol.ParseIpv6Datagram(ether.Payload)
		if err != nil {
			continue
		}
		if !ipv6Dgrm.ChecksumValid() {
			continue
		}
		go e.HandleIpv6(h, ipv6Dgrm, ether.Header.SrcHost)
	}
}

func (e *Engine) HandleIpv6(h *hal.IfHandle, dgrm *protocol.Ipv6Datagram, srcMac protocol.EthernetAddr) {
	dstIsMe := false
	if dgrm.Header.Dst.Equals(h.IPv6) {
		dstIsMe = true
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
				if rte.Metric == 0xff {
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
				h.SendIpv6(ipv6Reply, srcMac)
			}
		}
	} else {
		// forwarding

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
		if e.Metric >= 0xff {
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
