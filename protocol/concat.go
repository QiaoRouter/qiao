package protocol

func concatMac(s []byte, mac *EthernetAddr) []byte {
	for i := 0; i < len(mac.Octet); i++ {
		s = append(s, mac.Octet[i])
	}
	return s
}

func concatU16(s []byte, u16 uint16) []byte {
	s = append(s, uint8(u16>>8))
	s = append(s, uint8(u16))
	return s
}

func concatBuffer(s []byte, buf Buffer) []byte {
	TODO()
	return s
}
