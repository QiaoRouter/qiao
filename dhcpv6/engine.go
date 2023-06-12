package dhcpv6

import (
	"fmt"
	"qiao/hal"
	"qiao/protocol"
	"qiao/ripng"
	"time"
)

// 解决 Invalid receiver type '*ripng.Engine' ('ripng.Engine' is a non-local type 问题
// 解决方案链接
// https://stackoverflow.com/questions/44406077/golang-invalid-receiver-type-in-method-func

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

	// 插入默认路由
	// fixme @asu: 这里用指针会好一些吗
	entry := &ripng.RouteTableEntry{
		Ipv6Addr: protocol.Ipv6Addr{}, // all zero
		Len:      0,
		IfIndex:  1,
		Nexthop:  DefaultGateway,
	}
	_ = ripng.Update(entry)

	go e.ticker()
	for _, h := range hal.IfHandles {
		go e.receivePacketAndHandleIt(h)
	}
	for true {
		time.Sleep(time.Hour)
	}
	return nil
}
