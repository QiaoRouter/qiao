package protocol

import "encoding/binary"

func ConcatMac(s []byte, mac *EthernetAddr) []byte {
	for i := 0; i < len(mac.Octet); i++ {
		s = append(s, mac.Octet[i])
	}
	return s
}

func ConcatIpv6Addr(s []byte, ipv6Addr *Ipv6Addr) []byte {
	for i := 0; i < len(ipv6Addr.Octet); i++ {
		s = append(s, ipv6Addr.Octet[i])
	}
	return s
}

func ConcatNBytes(s []byte, N int, bytes []byte) []byte {
	for i := 0; i < N; i++ {
		s = append(s, bytes[i])
	}
	return s
}

func ConcatU16(s []byte, u16 uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, u16)
	s = append(s, b...)
	return s
}

func ConcatBuffer(s []byte, buf Buffer) []byte {
	for i := 0; i < len(buf.Octet); i++ {
		s = append(s, buf.Octet[i])
	}
	return s
}

func ConcatU32(s []byte, u32 uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, u32)
	s = append(s, b...)
	return s
}

func ConcatU8(s []byte, u8 uint8) []byte {
	s = append(s, u8)
	return s
}
