package protocol

const (
	EthernetAddrLen = 6  /* in unit of 1 octet */
	IPv6AddrLen     = 16 /* in unit of 1 octet */
	IPv6Version     = 6 << 4
)

type EthernetType uint16

type ICMPv6Type uint8

type NDOptionType uint8

const (
	NDOptionSourceLinkAddr = NDOptionType(1)
	NDOptionTargetLinkAddr = NDOptionType(2)
)

const (
	IPProtocolICMPV6 = 56
)
