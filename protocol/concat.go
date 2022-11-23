package protocol

import "encoding/binary"

func concatMac(s []byte, mac *EthernetAddr) []byte {
	for i := 0; i < len(mac.Octet); i++ {
		s = append(s, mac.Octet[i])
	}
	return s
}

func concatIpv6Addr(s []byte, ipv6Addr *Ipv6Addr) []byte {
	for i := 0; i < len(ipv6Addr.Octet); i++ {
		s = append(s, ipv6Addr.Octet[i])
	}
	return s
}

func concatU16(s []byte, u16 uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, u16)
	s = append(s, b...)
	return s
}

func concatBuffer(s []byte, buf Buffer) []byte {
	for i := 0; i < len(buf.Octet); i++ {
		s = append(s, buf.Octet[i])
	}
	return s
}

func concatU32(s []byte, u32 uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, u32)
	s = append(s, b...)
	return s
}

func concatU8(s []byte, u8 uint8) []byte {
	s = append(s, u8)
	return s
}
