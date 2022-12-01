package protocol

const (
	UdpHeaderLen = 8
)

type UdpPacket struct {
	DstPort  uint16
	SrcPort  uint16
	Len      uint16
	Checksum uint16

	Payload Buffer
}

func (udp *UdpPacket) Serialize() Buffer {
	var s []byte
	s = ConcatU16(s, udp.SrcPort)
	s = ConcatU16(s, udp.DstPort)
	s = ConcatU16(s, udp.Len)
	s = ConcatU16(s, udp.Checksum)
	s = ConcatBuffer(s, udp.Payload)
	return Buffer{
		Octet: s,
	}
}

func ParseUdp(buf Buffer) (*UdpPacket, error) {
	var err error
	parser := NetParser{
		Buffer:  buf,
		Pointer: 0,
	}
	udp := &UdpPacket{}
	udp.SrcPort, err = parser.ParseU16()
	if err != nil {
		return nil, err
	}
	udp.DstPort, err = parser.ParseU16()
	if err != nil {
		return nil, err
	}
	udp.Len, err = parser.ParseU16()
	if err != nil {
		return nil, err
	}
	udp.Checksum, err = parser.ParseU16()
	if err != nil {
		return nil, err
	}
	udp.Payload, err = parser.ParseBuffer()
	if err != nil {
		return nil, err
	}
	return udp, nil
}
