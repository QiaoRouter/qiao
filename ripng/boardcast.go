package ripng

import (
	"fmt"
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
		e.sendRipngs(i, ifHandles[i].IPv6)
	}
}

func (e *Engine) sendRipngs(ifhIdx int, ipv6Dst protocol.Ipv6Addr) {
	h := hal.IfHandles[ifhIdx]
	fmt.Printf("%v sendRipngs\n", h.IfName)
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
		fmt.Printf("lo: %d, hi: %d\n", lo, hi)
		ripngPacket.NumEntries = uint32(hi - lo)
		ripngPacket.Command = RipngTypeResponse
		for j := lo; j < hi; j++ {
			ripngPacket.Entries = append(ripngPacket.Entries,
				rtes[j].ToRipngEntry(ipv6Dst))
		}
		ipv6Dgrm := ripngPacket.ToIpv6UdpPacket(h.LinkLocalIPv6, ipv6Dst)

		go h.SendIpv6(ipv6Dgrm, macDst)
	}
}
