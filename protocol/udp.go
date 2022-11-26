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
