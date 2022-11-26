package ripng

import "qiao/protocol"

const (
	RipngMaxRte  = 72
	RipngUdpPort = 521
	RipngVersion = 1
)

func MakeRipngPacket() *RipngPacket {
	return &RipngPacket{}
}

type RipngPacket struct {
	NumEntries uint32
	Command    uint8
	Entries    []*RipngRte
}

//
// RipngRte -RIPng entry
// 定义
// https://datatracker.ietf.org/doc/html/rfc2080#page-6
// "Route Table Entry (RTE) has the following format:"
//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// ~                        IPv6 prefix (16)                       ~
// |                                                               |
// +---------------------------------------------------------------+
// |         route tag (2)         | prefix len (1)|  metric (1)   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
type RipngRte struct {
	Prefix    protocol.Ipv6Addr
	RouteTag  uint16
	PrefixLen uint8
	Metric    uint8
}

func (e *RipngRte) Serialize() protocol.Buffer {
	var s []byte
	s = protocol.ConcatBuffer(s, e.Prefix.Serialize())
	s = protocol.ConcatU16(s, e.RouteTag)
	s = protocol.ConcatU8(s, e.PrefixLen)
	s = protocol.ConcatU8(s, e.Metric)
	return protocol.Buffer{
		Octet: s,
	}
}

func (packet *RipngPacket) Serialize() protocol.Buffer {
	var s []byte
	s = protocol.ConcatU8(s, packet.Command)
	s = protocol.ConcatU8(s, RipngVersion)
	s = protocol.ConcatU16(s, 0)
	for _, e := range packet.Entries {
		s = protocol.ConcatBuffer(s, e.Serialize())
	}

	return protocol.Buffer{
		Octet: s,
	}
}

func (packet *RipngPacket) ToUdpPacket() *protocol.UdpPacket {
	buf := packet.Serialize()
	return &protocol.UdpPacket{
		DstPort: RipngUdpPort,
		SrcPort: RipngUdpPort,
		Len:     buf.Length() + protocol.UdpHeaderLen,
		Payload: buf,
	}
}

func (packet *RipngPacket) ToIpv6UdpPacket(src protocol.Ipv6Addr,
	dst protocol.Ipv6Addr) *protocol.Ipv6Datagram {
	udp := packet.ToUdpPacket()
	payload := udp.Serialize()
	ipv6Dgrm := &protocol.Ipv6Datagram{
		Header: protocol.Ipv6Header{
			Version:    6,
			PayloadLen: payload.Length(),
			NextHeader: protocol.IPProtocolUdp,
			HopLimit:   255,
			Src:        src,
			Dst:        dst,
		},
		Payload: payload,
	}
	return ipv6Dgrm
}
