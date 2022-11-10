package hal

import "github.com/google/gopacket/pcap"

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

var IfHandles = make(map[string]*pcap.Handle) // interface name -> pcap.handle
