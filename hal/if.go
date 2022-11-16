package hal

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"net/netip"

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
	MAC             net.HardwareAddr
	LinkLocalIPv6   netip.Addr
	NeighborMacAddr net.HardwareAddr
	IPv6            []netip.Addr
}

var IfHandles []*IfHandle
