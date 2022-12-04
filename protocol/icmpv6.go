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

type ICMPv6NDOption struct {
	Type          NDOptionType
	Len           uint8 /* in units of 8 octets */
	LinkLayerAddr EthernetAddr
}

func (opt *ICMPv6NDOption) Serialize() Buffer {
	var s []byte
	s = ConcatU8(s, uint8(opt.Type))
	s = ConcatU8(s, opt.Len)
	s = ConcatMac(s, &opt.LinkLayerAddr)
	return Buffer{
		Octet: s,
	}
}

type ICMPv6NeighborSolicitation struct {
	Header     ICMPv6Header
	TargetAddr Ipv6Addr
	Options    []ICMPv6NDOption
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
		Options: []ICMPv6NDOption{
			{
				Type:          NDOptionSourceLinkAddr,
				Len:           1,
				LinkLayerAddr: mac,
			},
		},
	}
	return icmp_ns
}

func MakeICMPv6NA(target Ipv6Addr, mac EthernetAddr) *ICMPv6NA {
	icmp_ns := &ICMPv6NA{
		Header: ICMPv6Header{
			Type:         ICMPv6TypeNeighborAdvertisement,
			Code:         0,
			RestOfHeader: 0x60000000,
		},
		Target: target,
		Options: []ICMPv6NDOption{
			{
				Type:          NDOptionTargetLinkAddr,
				Len:           1,
				LinkLayerAddr: mac,
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

func (p *NetParser) ParseICMPv6Header() (ICMPv6Header, error) {
	header := ICMPv6Header{}
	u8, err := p.ParseU8()
	if err != nil {
		return header, err
	}
	header.Type = ICMPv6Type(u8)
	header.Code, err = p.ParseU8()
	if err != nil {
		return header, err
	}
	header.Checksum, err = p.ParseU16()
	if err != nil {
		return header, err
	}
	header.RestOfHeader, err = p.ParseU32()
	if err != nil {
		return header, err
	}
	return header, nil
}

func ParseICMPv6(buf Buffer) (*ICMPv6, error) {
	icmpv6 := &ICMPv6{}
	p := NetParser{
		Buffer:  buf,
		Pointer: 0,
	}
	var err error
	icmpv6.Header, err = p.ParseICMPv6Header()
	if err != nil {
		return nil, err
	}
	icmpv6.MessageBody, err = p.ParseBuffer()
	if err != nil {
		return nil, err
	}
	return icmpv6, err
}

func MakeICMPv6Packet(icmpTye ICMPv6Type, code uint8, payload Buffer) *ICMPv6 {
	res := &ICMPv6{
		Header: ICMPv6Header{
			Type:         icmpTye,
			Code:         code,
			Checksum:     0,
			RestOfHeader: 0,
		},
		MessageBody: payload,
	}
	return res
}

type ICMPv6NA struct {
	Header  ICMPv6Header
	Target  Ipv6Addr
	Options []ICMPv6NDOption
}

func (na *ICMPv6NA) ToIpv6Datagram(src Ipv6Addr, dst Ipv6Addr) *Ipv6Datagram {
	payload := na.Serialize()
	return &Ipv6Datagram{
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
}

func (na *ICMPv6NA) Serialize() Buffer {
	var s []byte
	s = ConcatBuffer(s, na.Header.Serialize())
	s = ConcatIpv6Addr(s, &na.Target)
	for _, opt := range na.Options {
		s = ConcatU8(s, uint8(opt.Type))
		s = ConcatU8(s, opt.Len)
		s = ConcatMac(s, &opt.LinkLayerAddr)
	}
	return Buffer{
		Octet: s,
	}
}

func ParseICMPv6NeighborAdvert(buf Buffer) (*ICMPv6NA, error) {
	na := &ICMPv6NA{}
	var err error
	p := NetParser{
		Buffer:  buf,
		Pointer: 0,
	}
	na.Header, err = p.ParseICMPv6Header()
	if err != nil {
		return nil, err
	}
	na.Target, err = p.ParseIpv6Addr()
	if err != nil {
		return nil, err
	}
	return na, nil
}

func ParseICMPv6NeighborSolicitation(buf Buffer) (*ICMPv6NeighborSolicitation, error) {
	ns := &ICMPv6NeighborSolicitation{}
	var err error
	p := NetParser{
		Buffer: buf,
	}
	ns.Header, err = p.ParseICMPv6Header()
	if err != nil {
		return nil, err
	}
	ns.TargetAddr, err = p.ParseIpv6Addr()
	if err != nil {
		return nil, err
	}
	return ns, nil
}
