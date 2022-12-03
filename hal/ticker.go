package hal

import (
	"qiao/protocol"
	"time"
)

func ticker() {
	for true {
		// NDP 缓存
		m := NdpTable.m
		nm := make(map[protocol.Ipv6Addr]*NDPRecord)

		for addr := range m {
			record := m[addr]
			if record.ExpireTime.After(time.Now()) {
				nm[addr] = record
			}
		}
		NdpTable.Lock()
		NdpTable.m = nm
		NdpTable.Unlock()

		// 记录ndp请求
		nT := make(map[protocol.Ipv6Addr]time.Time)
		for addr, t := range ndpTimer {
			if t.Add(NDPRequestInterval).After(time.Now()) {
				nT[addr] = t
			}
		}
		ndpTimer = nT
		time.Sleep(time.Second * 10)
	}
}
