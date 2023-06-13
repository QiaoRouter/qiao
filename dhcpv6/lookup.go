package dhcpv6

import (
	"qiao/protocol"
	"qiao/ripng"
)

var LookUpTable struct {
	rteList []*RouteTableEntry
}

type RouteTableEntry struct {
	Ipv6Addr protocol.Ipv6Addr // 匹配的ipv6前缀
	Len      int               // 前缀长度
	IfIndex  int               // 出端口编号
	Nexthop  protocol.Ipv6Addr // 下一跳地址
	Metric   int               // 到达该前缀的成本
}

func Update(ne *RouteTableEntry) error {
	rteLst := LookUpTable.rteList
	for i := 0; i < len(rteLst); i++ {
		e := rteLst[i]
		if e.Ipv6Addr == ne.Ipv6Addr && e.Len == ne.Len {
			e.IfIndex = ne.IfIndex
			e.Nexthop = ne.Nexthop
			e.Metric = ne.Metric
			return nil
		}
	}
	LookUpTable.rteList = append(rteLst, ne)
	return nil
}

func PrefixQuery(addr protocol.Ipv6Addr) (ans *RouteTableEntry) {
	match := func(addr protocol.Ipv6Addr, pattern protocol.Ipv6Addr, length int) bool {
		for i := 0; i < length; i++ {
			idx := i / 8
			shift := 7 - (i % 8)
			addrBit := (addr.Octet[idx] >> shift) & 1
			patternBit := (pattern.Octet[idx] >> shift) & 1
			if addrBit != patternBit {
				return false
			}
		}
		return true
	}
	length := 0
	rteLst := LookUpTable.rteList
	for i := 0; i < len(rteLst); i++ {
		e := rteLst[i]
		if e.Len > length && match(addr, e.Ipv6Addr, e.Len) {
			length = e.Len
			ans = e
		}
	}
	return ans
}

func DelRte(de *RouteTableEntry) {
	rteLst := LookUpTable.rteList
	newLst := make([]*RouteTableEntry, len(rteLst))
	for i := 0; i < len(rteLst); i++ {
		e := rteLst[i]
		if e.Ipv6Addr == de.Ipv6Addr && e.Len == de.Len {
			continue
		}
		newLst = append(newLst, e)
	}
	LookUpTable.rteList = newLst
}

func ExactQuery(addr protocol.Ipv6Addr, length int) *RouteTableEntry {
	rteLst := LookUpTable.rteList
	for i := 0; i < len(rteLst); i++ {
		e := rteLst[i]
		if e.Ipv6Addr == addr && e.Len == length {
			return e
		}
	}
	return nil
}

func (e *RouteTableEntry) ToRipngEntry(ifIdx int) *ripng.RipngRte {
	rte := &ripng.RipngRte{
		Prefix:    e.Ipv6Addr,
		PrefixLen: uint8(e.Len),
		Metric:    uint8(e.Metric),
	}
	// Poisoned Reverse
	if e.IfIndex == ifIdx {
		rte.Metric = 16
	}
	return rte
}
