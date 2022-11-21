package protocol

import (
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
	Version      int      // protocol version
	TrafficClass int      // traffic class
	FlowLabel    int      // flow label
	PayloadLen   int      // payload length
	NextHeader   uint8    // next header
	HopLimit     int      // hop limit
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
	return &EthernetFrame{
		Header: EthernetHeader{
			DstHost: dstMac,
			SrcHost: srcMac,
			Type:    EthernetProtocolIPv6,
		},
		Payload: dgrm.Serialize(),
	}
}

func (dgrm *Ipv6Datagram) Serialize() Buffer {
	return Buffer{}
}
