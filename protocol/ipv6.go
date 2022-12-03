package protocol

import (
	"encoding/binary"
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
	s = ConcatU32(s, ctrl)

	s = ConcatU16(s, header.PayloadLen)
	s = ConcatU8(s, header.NextHeader)
	s = ConcatU8(s, header.HopLimit)
	s = ConcatIpv6Addr(s, &header.Src)
	s = ConcatIpv6Addr(s, &header.Dst)
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
	s = ConcatBuffer(s, dgrm.Header.Serialize())
	s = ConcatBuffer(s, dgrm.Payload)
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

func (dgrm *Ipv6Datagram) ChecksumValid() bool {
	checksum := Checksum32{}
	checksum.AddBuffer(dgrm.Header.Src.Serialize())
	checksum.AddBuffer(dgrm.Header.Dst.Serialize())
	checksum.AddU16(dgrm.Header.PayloadLen)
	checksum.AddU8(dgrm.Header.NextHeader)
	checksum.AddBuffer(dgrm.Payload)
	return checksum.U16() == 0xffff
}

func (dgrm *Ipv6Datagram) FillChecksum() {
	if dgrm.Header.NextHeader == IPProtocolICMPV6 {
		binary.BigEndian.PutUint16(dgrm.Payload.Octet[2:4], 0)
	} else if dgrm.Header.NextHeader == IPProtocolUdp {
		binary.BigEndian.PutUint16(dgrm.Payload.Octet[6:8], 0)
	}
	checksum := Checksum32{}
	checksum.AddBuffer(dgrm.Header.Src.Serialize())
	checksum.AddBuffer(dgrm.Header.Dst.Serialize())
	checksum.AddU16(dgrm.Header.PayloadLen)
	checksum.AddU8(dgrm.Header.NextHeader)
	checksum.AddBuffer(dgrm.Payload)
	if dgrm.Header.NextHeader == IPProtocolICMPV6 {
		u16 := 0xffff - checksum.U16()
		binary.BigEndian.PutUint16(dgrm.Payload.Octet[2:4], u16)
	}
	if dgrm.Header.NextHeader == IPProtocolUdp {
		u16 := 0xffff - checksum.U16()
		if u16 == 0 {
			u16 = 0xffff
		}
		binary.BigEndian.PutUint16(dgrm.Payload.Octet[6:8], u16)
	}
}

func LenToMaskAddr(maskLen int) Ipv6Addr {
	ret := Ipv6Addr{}
	if maskLen < 0 || maskLen > 128 {
		return ret
	}
	idx := 0
	for maskLen >= 8 {
		ret.Octet[idx] = 0xff
		maskLen -= 8
		idx++
	}
	if maskLen > 0 {
		x := byte(0b10000000)
		for i := 0; i < maskLen; i++ {
			ret.Octet[idx] += x >> idx
		}
	}
	return ret
}

func (addr *Ipv6Addr) ToRouteAddr(mask int) Ipv6Addr {
	ret := Ipv6Addr{}
	maskAddr := LenToMaskAddr(mask)
	for i := 0; i < len(addr.Octet); i++ {
		ret.Octet[i] = addr.Octet[i] & maskAddr.Octet[i]
	}
	return ret
}

func (addr *Ipv6Addr) Equals(ipv6 Ipv6Addr) bool {
	for i := 0; i < IPv6AddrLen; i++ {
		if addr.Octet[i] != ipv6.Octet[i] {
			return false
		}
	}
	return true
}

func ParseIpv6Datagram(buf Buffer) (*Ipv6Datagram, error) {
	parser := NetParser{
		Buffer:  buf,
		Pointer: 0,
	}
	dgrm := &Ipv6Datagram{}
	ctrl, err := parser.ParseU32()
	if err != nil {
		return nil, err
	}
	if uint8(ctrl>>28) == 6 {
		dgrm.Header.Version = 6
	} else {
		return nil, ParseErr
	}
	dgrm.Header.TrafficClass = uint8((ctrl >> 20) & 0xff)
	dgrm.Header.FlowLabel = ctrl & 0xfffff
	dgrm.Header.PayloadLen, err = parser.ParseU16()
	if err != nil {
		return nil, err
	}
	dgrm.Header.NextHeader, err = parser.ParseU8()
	if err != nil {
		return nil, err
	}
	dgrm.Header.HopLimit, err = parser.ParseU8()
	if err != nil {
		return nil, err
	}
	dgrm.Header.Src, err = parser.ParseIpv6Addr()
	if err != nil {
		return nil, err
	}
	dgrm.Header.Dst, err = parser.ParseIpv6Addr()
	if err != nil {
		return nil, err
	}
	dgrm.Payload, err = parser.ParseBuffer()
	if err != nil {
		return nil, err
	}
	return dgrm, nil
}

func (addr *Ipv6Addr) AllZero() bool {
	for _, v := range addr.Octet {
		if v != 0 {
			return false
		}
	}
	return true
}
