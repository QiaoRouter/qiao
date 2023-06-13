package dhcpv6

import "qiao/protocol"

// @asu fixme 这里 ipv6 地址似乎不用硬编码些这么长，好像可以使用 ParseIp 还是什么来着
var DefaultGateway = protocol.Ipv6Addr{
	Octet: [protocol.IPv6AddrLen]byte{
		0xfd, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x02,
	},
}

// 2402:f000:1:801::8:28
var DnsA = protocol.Ipv6Addr{
	Octet: [protocol.IPv6AddrLen]byte{
		0x24, 0x02, 0xf0, 0x00,
		0x00, 0x01, 0x08, 0x01,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x08, 0x02, 0x08,
	},
}

// 和 2402:f000:1:801::8:29
var DnsB = protocol.Ipv6Addr{
	Octet: [protocol.IPv6AddrLen]byte{
		0x24, 0x02, 0xf0, 0x00,
		0x00, 0x01, 0x08, 0x01,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x08, 0x02, 0x09,
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
	DHCP_REPLY     = 7
)

// dhcp option code
const (
	DHCP_OPTION_CLIENTID    = 1
	DHCP_OPTION_SERVERID    = 2
	DHCP_OPTION_IA_NA       = 3
	DHCP_OPTION_IAADDR      = 5
	DHCP_OPTION_DNS_SERVERS = 23
)

// dhcp column length
const (
	DHCP_MSG_TYPE_LEN = 1
	DHCP_TXN_CODE_LEN = 3
)

// dhcp options length
const (
	DHCP_OPTIONS_HDR_LEN         = 4
	DHCP_OPTIONS_DUID_LEN        = 14
	DHCP_OPTIONS_IANA_LEN        = 12
	DHCP_OPTIONS_IANA_OPTION_LEN = 24
)

// IA_NA options parameters
const (
	Preferredlifetime = 54000
	Validlifetime     = 86400
)
