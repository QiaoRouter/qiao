package hal

import (
	"fmt"
	"github.com/mdlayher/ndp"
	"net"
	"net/netip"
)

func (h *IfHandle) GetNeighborMacAddr() (net.HardwareAddr, error) {
	target, err := netip.ParseAddr(h.LinkLocalIPv6.String())
	if err != nil {
		panic(err)
	}
	snm, err := ndp.SolicitedNodeMulticast(target)
	if err != nil {
		panic(err)
	}
	fmt.Printf("if_%+v link-local ip is %+v, snm is %+v\n",
		h.IfName, h.LinkLocalIPv6, snm)
	p, err := h.PacketSource.NextPacket()
	if err != nil {
		panic(err)
	}

	fmt.Printf("收到一个数据包， 原封不动发回去, p: %+v\n",
		p.String())
	err = h.PcapHandleOut.WritePacketData(p.Data())
	if err != nil {
		panic(err)
	}

	return nil, nil
}
