package hal

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"qiao/protocol"

	"net"
)

type HostName string
type IfNames []string

type IfHandle struct {
	IfName          string
	IfIndex         int
	PcapHandleIn    *pcap.Handle
	PcapHandleOut   *pcap.Handle
	PacketSource    *gopacket.PacketSource
	MAC             protocol.EthernetAddr
	LinkLocalIPv6   protocol.Ipv6Addr
	NeighborMacAddr net.HardwareAddr
	IPv6            []protocol.Ipv6Addr
}

var IfHandles []*IfHandle
