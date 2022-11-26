package protocol

import (
	"encoding/binary"
	"fmt"
)

type Checksum32 struct {
	Val uint32
}

func (check *Checksum32) AddBuffer(buf Buffer) {
	u32 := check.Val
	_len := len(buf.Octet)
	index := 0
	for _len > 1 {
		u16 := binary.BigEndian.Uint16(buf.Octet[index : index+2])
		u32 = u32 + uint32(u16)
		_len -= 2
		index += 2
	}
	if _len > 0 {
		u32 += uint32(buf.Octet[index]) << 8
	}
	check.Val = u32
}

func (check *Checksum32) AddU16(u16 uint16) {
	check.Val = check.Val + uint32(u16)
}

func (check *Checksum32) AddU8(u8 uint8) {
	check.Val = check.Val + uint32(u8)
}

func (check *Checksum32) U16() uint16 {
	u32 := check.Val
	for u32 > 0xffff {
		u32 = (u32 & 0xffff) + (u32 >> 16)
	}
	return uint16(u32)
}

func (dgrm *Ipv6Datagram) FillChecksum() {
	checksum := Checksum32{}
	checksum.AddBuffer(dgrm.Header.Src.Serialize())
	checksum.AddBuffer(dgrm.Header.Dst.Serialize())
	checksum.AddU16(dgrm.Header.PayloadLen)
	checksum.AddU8(dgrm.Header.NextHeader)
	checksum.AddBuffer(dgrm.Payload)
	u16 := 0xffff - checksum.U16()
	if dgrm.Header.NextHeader == IPProtocolICMPV6 {
		binary.BigEndian.PutUint16(dgrm.Payload.Octet[2:4], u16)
	}
	if dgrm.Header.NextHeader == IPProtocolUdp {
		fmt.Printf("next header is udp\n")
		if u16 == 0 {
			u16 = 0xffff
		}
		fmt.Printf("u16: %d\n", u16)
		binary.BigEndian.PutUint16(dgrm.Payload.Octet[6:8], u16)
	}
}
