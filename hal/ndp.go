package hal

import (
	"errors"
	"qiao/protocol"
	"sync"
	"time"
)

var (
	NDPErrIpNotFound   = errors.New("ip not found")
	NDPRecordTimeout   = time.Second * 14400
	NDPRequestInterval = time.Second
)

var NdpTable struct {
	sync.Mutex
	m map[protocol.Ipv6Addr]*NDPRecord
}

type NDPRecord struct {
	Mac        protocol.EthernetAddr
	ExpireTime time.Time
}

var ndpTimer = map[protocol.Ipv6Addr]time.Time{}

func GetNeighborMacAddr(ifIdx int, ip protocol.Ipv6Addr) (protocol.EthernetAddr, error) {
	etherAddr := protocol.EthernetAddr{}
	for _, h := range IfHandles {
		if h.IfIndex == ifIdx {
			return h.getNeighborMacAddr(ip)
		}
	}
	return etherAddr, nil
}

func (h *IfHandle) getNeighborMacAddr(ip protocol.Ipv6Addr) (protocol.EthernetAddr, error) {
	addr := protocol.EthernetAddr{}
	if NdpTable.m[ip] != nil {
		return NdpTable.m[ip].Mac, nil
	}
	if !ndpTimer[ip].Add(NDPRequestInterval).Before(time.Now()) {
		return addr, NDPErrIpNotFound
	}
	snm := ip.SolicitedNodeMulticast()
	icmp_ns := protocol.MakeICMPv6NeighborSolicitation(ip, h.MAC)
	var err error
	ipv6_datagram := icmp_ns.ToIpv6Datagram(h.IPv6, snm)
	ether := ipv6_datagram.ToEthernetFrame(h.MAC, snm.MulticastMac())
	err = h.PcapHandleOut.WritePacketData(ether.Serialize())
	if err != nil {
		return addr, err
	}
	return addr, NDPErrIpNotFound
}
