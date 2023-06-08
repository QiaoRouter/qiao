package dhcpv6

import (
	"fmt"
	"qiao/hal"
	"qiao/protocol"
	"qiao/ripng"
)

type Engine struct {
}

func MakeDhcpv6Engine() *Engine {
	return &Engine{}
}

func (e *Engine) Run() error {
	hal.Init()
	defer hal.Close()

	// 插入直连路由
	for _, h := range hal.IfHandles {
		fmt.Printf("插入%+v\n", h.IPv6.String())
		e := &ripng.RouteTableEntry{
			Ipv6Addr: h.IPv6.ToRouteAddr(h.IPv6Mask),
			Len:      h.IPv6Mask,
			IfIndex:  h.IfIndex,
			Nexthop:  protocol.Ipv6Addr{}, // all zero, link-local route
			Metric:   1,
		}
		fmt.Printf("e:%+v\n", e)
		err := ripng.Update(e)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
