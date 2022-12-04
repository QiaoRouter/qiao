package protocol

import (
	"encoding/binary"
	"errors"
)

var ParseErr = errors.New("net parse err")

type NetParser struct {
	Buffer  Buffer
	Pointer int
}

func (p *NetParser) ParseEthernetAddr() (EthernetAddr, error) {
	etherAddr := EthernetAddr{}
	if (p.Pointer + EthernetAddrLen) > int(p.Buffer.Length()) {
		return etherAddr, ParseErr
	}
	for i := 0; i < EthernetAddrLen; i++ {
		etherAddr.Octet[i] = p.Buffer.Octet[p.Pointer+i]
	}
	p.Pointer += EthernetAddrLen
	return etherAddr, nil
}

func (p *NetParser) ParseU16() (uint16, error) {
	if p.Pointer+2 > int(p.Buffer.Length()) {
		return 0, ParseErr
	}
	u16 := binary.BigEndian.Uint16(p.Buffer.Octet[p.Pointer : p.Pointer+2])
	p.Pointer += 2
	return u16, nil
}

func (p *NetParser) ParseU32() (uint32, error) {
	if p.Pointer+4 > int(p.Buffer.Length()) {
		return 0, ParseErr
	}
	u32 := binary.BigEndian.Uint32(p.Buffer.Octet[p.Pointer : p.Pointer+4])
	p.Pointer += 4
	return u32, nil
}

func (p *NetParser) ParseBuffer() (Buffer, error) {
	buf := Buffer{}
	buf.Octet = p.Buffer.Octet[p.Pointer:]
	p.Pointer = int(p.Buffer.Length())
	return buf, nil
}

func (p *NetParser) ParseU8() (uint8, error) {
	if p.Pointer+1 > int(p.Buffer.Length()) {
		return 0, ParseErr
	}
	u8 := p.Buffer.Octet[p.Pointer]
	p.Pointer++
	return u8, nil
}

func (p *NetParser) ParseIpv6Addr() (Ipv6Addr, error) {
	addr := Ipv6Addr{}
	if p.Pointer+IPv6AddrLen > int(p.Buffer.Length()) {
		return addr, ParseErr
	}
	for i := 0; i < IPv6AddrLen; i++ {
		addr.Octet[i] = p.Buffer.Octet[p.Pointer+i]
	}
	p.Pointer += IPv6AddrLen
	return addr, nil
}

func (p *NetParser) Eof() bool {
	return p.Pointer == int(p.Buffer.Length())
}
