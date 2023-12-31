package ripng

import (
	"qiao/hal"
	"qiao/protocol"
)

const (
	RipngTypeRequest  = 1
	RipngTypeResponse = 2
)

// 向所有interface发送完整路由表
func (e *Engine) broadcast() {
	ifHandles := hal.IfHandles
	for i := 0; i < len(ifHandles); i++ {
		ipv6Dst, _ := protocol.ParseIpv6("ff02::9")
		e.sendRipngs(i, ipv6Dst)
	}
}

func (e *Engine) sendRipngs(handleIdx int, ipv6Dst protocol.Ipv6Addr) {
	h := hal.IfHandles[handleIdx]
	macDst := ipv6Dst.MulticastMac()
	rtes := LookUpTable.rteList

	for i := 0; i < len(rtes)/RipngMaxRte+1; i++ {
		ripngPacket := &RipngPacket{}
		lo := i * RipngMaxRte
		hi := 0
		if lo+RipngMaxRte < len(rtes) {
			hi = lo + RipngMaxRte
		} else {
			hi = len(rtes)
		}
		ripngPacket.NumEntries = uint32(hi - lo)
		ripngPacket.Command = RipngTypeResponse
		for j := lo; j < hi; j++ {
			ripngPacket.Entries = append(ripngPacket.Entries,
				rtes[j].ToRipngEntry(h.IfIndex))
		}
		ipv6Dgrm := ripngPacket.ToIpv6UdpPacket(h.LinkLocalIPv6, ipv6Dst)

		go hal.SendIpv6(h.IfIndex, ipv6Dgrm, macDst)
	}
}
