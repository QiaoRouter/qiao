package hal

import (
	"fmt"
	"net"
	"qiao/protocol"
	"sync"
	"time"
)

var NdpTable struct {
	sync.Mutex
	m map[protocol.Ipv6Addr]NDPRecord
}

type NDPRecord struct {
	Mac        protocol.EthernetAddr
	ExpireTime time.Time
}

func (h *IfHandle) GetNeighborMacAddr() (net.HardwareAddr, error) {
	target := h.LinkLocalIPv6
	snm := target.SolicitedNodeMulticast()
	fmt.Printf("if_%+v link-local ip is %+v, snm is %+v\n",
		h.IfName, target.String(), snm.String())

	icmp_ns := protocol.MakeICMPv6NeighborSolicitation(target, h.MAC)
	ipv6_datagram := icmp_ns.ToIpv6Datagram(h.IPv6[0], snm)
	snm_mac := snm.MulticastMac()
	fmt.Printf("src_ether is %s, dst_ether is %s\n", h.MAC.String(), snm_mac.String())
	ether := ipv6_datagram.ToEthernetFrame(h.MAC, snm.MulticastMac())
	fmt.Printf("icmp_ns: %+v\n", icmp_ns)
	fmt.Printf("ipv6_datagram: %+v\n", ipv6_datagram)
	fmt.Printf("ether: %+v\n\n", ether)
	fmt.Printf("ether Serialize: %+v\n\n", ether.Serialize())
	err := h.PcapHandleOut.WritePacketData(ether.Serialize())
	return nil, err
}
