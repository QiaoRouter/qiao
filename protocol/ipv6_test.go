package protocol

import "testing"

func isEqual(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestIpv6Datagram_ToEthernetFrame(t *testing.T) {
	ipv6Dgrm := Ipv6Datagram{}
	src, _ := ParseMac("32:ca:2c:70:64:fc")
	dst, _ := ParseMac("6a:de:01:a2:3c:3f")
	etherFrame := ipv6Dgrm.ToEthernetFrame(src, dst)
	ans := []byte{
		0x32, 0xca, 0x2c, 0x70, 0x64, 0xfc, 0x6a, 0xde, 0x01, 0xa2, 0x3c, 0x3f, 0x86, 0xdd,
	}
	if !isEqual(ans, etherFrame.Serialize()) {
		t.Fatal("ToEthernetFrame check Serialize fail")
	}
}
