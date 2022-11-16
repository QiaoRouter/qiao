package hal

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	"net"
)

type HostName string
type IfNames []string

type IfHandle struct {
	IfName          string
	IfIndex         int
	PcapHandle      *pcap.Handle
	PacketSource    *gopacket.PacketSource
	MAC             net.HardwareAddr
	LinkLocalIPv6   net.IP
	NeighborMacAddr net.HardwareAddr
}

var IfHandles []*IfHandle
