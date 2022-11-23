package protocol

type ICMPv6Header struct {
	Type         ICMPv6Type
	Code         uint8
	Checksum     uint16
	RestOfHeader uint32
}

type ICMPv6NeighborSolicitationOption struct {
	Type    NDOptionType
	Len     uint8 /* in units of 8 octets */
	Payload Buffer
}

func (opt *ICMPv6NeighborSolicitationOption) Serialize() Buffer {
	var s []byte
	s = concatU8(s, uint8(opt.Type))
	s = concatU8(s, opt.Len)
	s = concatBuffer(s, opt.Payload)
	return Buffer{
		Octet: s,
	}
}

type ICMPv6NeighborSolicitation struct {
	Header     ICMPv6Header
	TargetAddr Ipv6Addr
	Options    []ICMPv6NeighborSolicitationOption
}

func (icmp *ICMPv6NeighborSolicitation) Serialize() Buffer {
	var s []byte
	s = concatU8(s, uint8(icmp.Header.Type))
	s = concatU8(s, icmp.Header.Code)
	s = concatU16(s, 0)
	s = concatU32(s, icmp.Header.RestOfHeader)
	s = concatBuffer(s, icmp.TargetAddr.Serialize())
	for i := 0; i < len(icmp.Options); i++ {
		s = concatBuffer(s, icmp.Options[i].Serialize())
	}

	return Buffer{
		Octet: s,
	}
}

func MakeICMPv6NeighborSolicitation(linkLocalIpv6 Ipv6Addr,
	mac EthernetAddr) *ICMPv6NeighborSolicitation {
	icmp_ns := &ICMPv6NeighborSolicitation{
		Header: ICMPv6Header{
			Type:         135,
			Code:         0,
			RestOfHeader: 0,
		},
		TargetAddr: linkLocalIpv6,
		Options: []ICMPv6NeighborSolicitationOption{
			{
				Type:    NDOptionSourceLinkAddr,
				Len:     1,
				Payload: mac.Serialize(),
			},
		},
	}
	return icmp_ns
}

func (icmp *ICMPv6NeighborSolicitation) ToIpv6Datagram(src Ipv6Addr, dst Ipv6Addr) *Ipv6Datagram {
	payload := icmp.Serialize()
	dgrm := &Ipv6Datagram{
		Header: Ipv6Header{
			Version:    6,
			FlowLabel:  0,
			PayloadLen: payload.Length(),
			NextHeader: IPProtocolICMPV6,
			HopLimit:   255,
			Src:        src,
			Dst:        dst,
		},
		Payload: payload,
	}
	return dgrm
}
