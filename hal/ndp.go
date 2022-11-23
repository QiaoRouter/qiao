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

func (h *IfHandle) GetNeighborMacAddr(ip protocol.Ipv6Addr) (net.HardwareAddr, error) {
	snm := ip.SolicitedNodeMulticast()
	icmp_ns := protocol.MakeICMPv6NeighborSolicitation(ip, h.MAC)
	var err error
	for i := 0; i < len(h.IPv6); i++ {
		ipv6_datagram := icmp_ns.ToIpv6Datagram(h.IPv6[i], snm)
		//fmt.Printf("src_ip is %s, dst_ip is %s\n", h.IPv6[0].String(), snm)
		snm_mac := snm.MulticastMac()
		fmt.Printf("src_ether is %s, dst_ether is %s\n", h.MAC.String(), snm_mac.String())
		ether := ipv6_datagram.ToEthernetFrame(h.MAC, snm.MulticastMac())
		fmt.Printf("icmp_ns: %+v\n", icmp_ns)
		fmt.Printf("ipv6_datagram: %+v\n", ipv6_datagram)
		err = h.PcapHandleOut.WritePacketData(ether.Serialize())
	}

	return nil, err
}
