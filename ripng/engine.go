package ripng

import (
	"qiao/hal"
	"qiao/protocol"
	"time"
)

type Engine struct {
}

func MakeRipngEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Run() error {
	hal.Init()
	defer hal.Close()

	// 插入直连路由
	for i := range hal.IfHandles {
		h := hal.IfHandles[i]
		e := &RouteTableEntry{
			Ipv6Addr: h.IPv6.ToRouteAddr(h.IPv6Mask),
			Len:      h.IPv6Mask,
			IfIndex:  h.IfIndex,
			Nexthop:  protocol.Ipv6Addr{}, // all zero, link-local route
			Metric:   1,
		}
		err := Update(e)
		if err != nil {
			panic(err)
		}
	}
	go e.ticker()
	for _, h := range hal.IfHandles {
		go e.receivePacketAndHandleIt(h)
	}
	for true {
		time.Sleep(time.Hour)
	}
	return nil
}
