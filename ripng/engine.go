package ripng

import (
	"fmt"
	"qiao/hal"
	"qiao/protocol"
)

type Engine struct {
}

func MakeRipngEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Run() error {
	hal.Init()

	fmt.Printf("we have %d if_handles\n", len(hal.IfHandles))
	// 插入直连路由
	for i := range hal.IfHandles {
		h := hal.IfHandles[i]
		fmt.Printf("ipv6: %+v\n", h.IPv6)
		e := &RouteTableEntry{
			Ipv6Addr: h.IPv6,
			Len:      h.IPv6Mask,
			IfIndex:  h.IfIndex,
			Nexthop:  protocol.Ipv6Addr{}, // all zero, link-local route
			Metric:   1,
		}
		fmt.Printf("e: %+v\n", e)
		err := AddRte(e)
		if err != nil {
			panic(err)
		}
	}
	go e.ticker()
	for true {

	}
	return nil
}
