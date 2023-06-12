package dhcpv6

import "qiao/protocol"

type DHCPv6 struct {
	Header      DHCPv6Header
	MessageBody protocol.Buffer
}

type DHCPv6Header struct {
}
