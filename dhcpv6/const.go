package dhcpv6

import "qiao/protocol"

var DefaultGateway = protocol.Ipv6Addr{
	Octet: [protocol.IPv6AddrLen]byte{
		0xfd, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x02,
	},
}

// dhcp's udp port
const (
	DHCP_SERVER = 547
	DHCP_CLIENT = 546
)

// dhcp mst-type
const (
	DHCP_SOLICIT   = 1
	DHCP_ADVERTISE = 2
	DHCP_REQUEST   = 3
)

// dhcp option code
const (
	DHCP_OPTION_CLIENTID    = 1
	DHCP_OPTION_SERVERID    = 2
	DHCP_OPTION_IA_NA       = 3
	DHCP_OPTION_IAADDR      = 5
	DHCP_OPTION_DNS_SERVERS = 23
)
