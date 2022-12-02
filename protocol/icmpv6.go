package protocol

type ICMPv6Header struct {
	Type         ICMPv6Type
	Code         uint8
	Checksum     uint16
	RestOfHeader uint32
}

func (h *ICMPv6Header) Serialize() Buffer {
	var s []byte
	s = ConcatU8(s, uint8(h.Type))
	s = ConcatU8(s, h.Code)
	s = ConcatU16(s, h.Checksum)
	s = ConcatU32(s, h.RestOfHeader)
	return Buffer{
		Octet: s,
	}
}

type ICMPv6NeighborSolicitationOption struct {
	Type    NDOptionType
	Len     uint8 /* in units of 8 octets */
	Payload Buffer
}

func (opt *ICMPv6NeighborSolicitationOption) Serialize() Buffer {
	var s []byte
	s = ConcatU8(s, uint8(opt.Type))
	s = ConcatU8(s, opt.Len)
	s = ConcatBuffer(s, opt.Payload)
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
	s = ConcatU8(s, uint8(icmp.Header.Type))
	s = ConcatU8(s, icmp.Header.Code)
	s = ConcatU16(s, 0)
	s = ConcatU32(s, icmp.Header.RestOfHeader)
	s = ConcatBuffer(s, icmp.TargetAddr.Serialize())
	for i := 0; i < len(icmp.Options); i++ {
		s = ConcatBuffer(s, icmp.Options[i].Serialize())
	}

	return Buffer{
		Octet: s,
	}
}

func MakeICMPv6NeighborSolicitation(linkLocalIpv6 Ipv6Addr,
	mac EthernetAddr) *ICMPv6NeighborSolicitation {
	icmp_ns := &ICMPv6NeighborSolicitation{
		Header: ICMPv6Header{
			Type:         ICMPv6TypeNeighborSolicitation,
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

func (icmp *ICMPv6) ToIpv6Datagram(src Ipv6Addr, dst Ipv6Addr, hopLimit uint8) *Ipv6Datagram {
	payload := icmp.Serialize()
	dgrm := &Ipv6Datagram{
		Header: Ipv6Header{
			Version:    6,
			PayloadLen: payload.Length(),
			NextHeader: IPProtocolICMPV6,
			HopLimit:   hopLimit,
			Src:        src,
			Dst:        dst,
		},
		Payload: payload,
	}
	return dgrm
}

func (icmp *ICMPv6) Serialize() Buffer {
	var s []byte
	s = ConcatBuffer(s, icmp.Header.Serialize())
	s = ConcatBuffer(s, icmp.MessageBody)
	return Buffer{
		Octet: s,
	}
}

type ICMPv6 struct {
	Header      ICMPv6Header
	MessageBody Buffer
}

func ParseICMPv6(buf Buffer) (*ICMPv6, error) {
	icmpv6 := &ICMPv6{}
	parser := NetParser{
		Buffer:  buf,
		Pointer: 0,
	}
	u8, err := parser.ParseU8()
	if err != nil {
		return nil, err
	}
	icmpv6.Header.Type = ICMPv6Type(u8)
	icmpv6.Header.Code, err = parser.ParseU8()
	if err != nil {
		return nil, err
	}
	icmpv6.Header.Checksum, err = parser.ParseU16()
	if err != nil {
		return nil, err
	}
	icmpv6.Header.RestOfHeader, err = parser.ParseU32()
	if err != nil {
		return nil, err
	}
	icmpv6.MessageBody, err = parser.ParseBuffer()
	if err != nil {
		return nil, err
	}
	return icmpv6, err
}
