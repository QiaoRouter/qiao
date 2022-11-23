package protocol

import (
	"fmt"
	"github.com/mdlayher/ndp"
	"net/netip"
)

type Ipv6Addr struct {
	Octet [IPv6AddrLen]byte
}

func (addr *Ipv6Addr) String() string {
	return netip.AddrFrom16(addr.Octet).String()
}

func ParseIpv6(s string) (Ipv6Addr, error) {
	ip, err := netip.ParseAddr(s)
	if err != nil {
		return Ipv6Addr{}, err
	}
	return Ipv6Addr{
		ip.As16(),
	}, nil
}

type Ipv6Header struct {
	Version      uint8    // protocol version
	TrafficClass uint8    // traffic class
	FlowLabel    uint32   // flow label
	PayloadLen   uint16   // payload length
	NextHeader   uint8    // next header
	HopLimit     uint8    // hop limit
	Src          Ipv6Addr // source address
	Dst          Ipv6Addr // destination address
}

type Ipv6Datagram struct {
	Header  Ipv6Header
	Payload Buffer
}

func (addr *Ipv6Addr) NetIP() netip.Addr {
	return netip.AddrFrom16(addr.Octet)
}

func (addr *Ipv6Addr) SolicitedNodeMulticast() Ipv6Addr {
	snm, err := ndp.SolicitedNodeMulticast(addr.NetIP())
	if err != nil {
		panic(err)
	}
	res, err := ParseIpv6(snm.String())
	return res
}

func (addr *Ipv6Addr) MulticastMac() EthernetAddr {
	mac := EthernetAddr{}
	mac.Octet[0] = 0x33
	mac.Octet[1] = 0x33
	for i := 0; i < 4; i++ {
		mac.Octet[2+i] = addr.Octet[12+i]
	}
	return mac
}

func (dgrm *Ipv6Datagram) ToEthernetFrame(srcMac EthernetAddr, dstMac EthernetAddr) *EthernetFrame {
	dgrm.FillChecksum()

	return &EthernetFrame{
		Header: EthernetHeader{
			DstHost: dstMac,
			SrcHost: srcMac,
			Type:    EthernetProtocolIPv6,
		},
		Payload: dgrm.Serialize(),
	}
}

func (header *Ipv6Header) Serialize() Buffer {
	var s []byte
	ctrl := uint32(0)
	if header.Version == 6 {
		ctrl += uint32(IPv6Version) << 24
	} else {
		panic("not ipv6")
	}

	ctrl += uint32(header.TrafficClass) << 20
	ctrl += header.FlowLabel
	s = concatU32(s, ctrl)

	s = concatU16(s, header.PayloadLen)
	s = concatU8(s, header.NextHeader)
	s = concatU8(s, header.HopLimit)
	s = concatIpv6Addr(s, &header.Src)
	s = concatIpv6Addr(s, &header.Dst)
	return Buffer{
		Octet: s,
	}
}

func (addr *Ipv6Addr) Serialize() Buffer {
	s := make([]byte, 16)
	for i := 0; i < len(addr.Octet); i++ {
		s[i] = addr.Octet[i]
	}
	return Buffer{
		Octet: s,
	}
}

func (dgrm *Ipv6Datagram) Serialize() Buffer {
	var s []byte
	s = concatBuffer(s, dgrm.Header.Serialize())
	s = concatBuffer(s, dgrm.Payload)
	return Buffer{
		Octet: s,
	}
}

func (dgrm *Ipv6Datagram) String() string {
	s := ""
	s += fmt.Sprintf("header: %+v, ", dgrm.Header)
	s += fmt.Sprintf("src: %+v, ", dgrm.Header.Src.String())
	s += fmt.Sprintf("dst: %+v, ", dgrm.Header.Dst.String())
	return s
}
