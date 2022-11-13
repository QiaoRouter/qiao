package hal

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	"net"
)

type HostName string
type IfNames []string

var Host = HostName("R2")

// for experiments
var experimentInterfaces = map[HostName]IfNames{
	"PC1": {"pc1r1", "eth2", "eth3", "eth4"},
	"R1":  {"r1pc1", "r1r2", "eth3", "eth4"},
	"R2":  {"r2r1", "r2r3", "eth3", "eth4", "en0"},
	"R3":  {"r3r2", "r3pc2", "eth3", "eth4"},
}

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
