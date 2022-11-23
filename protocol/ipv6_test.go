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
	srcIpv6, _ := ParseIpv6("fd00::5:1")
	dstIpv6, _ := ParseIpv6("fd00::1:2")
	ipv6Dgrm := Ipv6Datagram{
		Header: Ipv6Header{
			Version:      6,
			TrafficClass: 0,
			FlowLabel:    0,
			PayloadLen:   0,
			NextHeader:   IPProtocolICMPV6,
			HopLimit:     255,
			Src:          srcIpv6,
			Dst:          dstIpv6,
		},
		Payload: Buffer{},
	}
	srcMac, _ := ParseMac("32:ca:2c:70:64:fc")
	dstMac, _ := ParseMac("6a:de:01:a2:3c:3f")
	etherFrame := ipv6Dgrm.ToEthernetFrame(srcMac, dstMac)
	TODO()
	ans := []byte{
		0x32, 0xca, 0x2c, 0x70, 0x64, 0xfc, 0x6a, 0xde, 0x01, 0xa2, 0x3c, 0x3f, 0x86, 0xdd,
	}
	if !isEqual(ans, etherFrame.Serialize()) {
		t.Fatal("ToEthernetFrame check Serialize fail")
	}
}
