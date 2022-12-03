package protocol

const (
	EthernetAddrLen   = 6  /* in unit of 1 octet */
	IPv6AddrLen       = 16 /* in unit of 1 octet */
	IPv6Version       = uint8(6 << 4)
	EthernetHeaderLen = 14
	Ipv6HeaderLen     = 40
	ICMPv6HeaderLen   = 8
	Ipv6MinimumMTU    = 1280
)

type EthernetType uint16

type ICMPv6Type uint8

const (
	ICMPv6TypeDestinationUnreachable = ICMPv6Type(3)
	ICMPv6TypeNeighborSolicitation   = ICMPv6Type(135)
	ICMPv6TypeNeighborAdvertisement  = ICMPv6Type(136)
	ICMPv6TypeEchoRequest            = ICMPv6Type(128)
	ICMPv6TypeEchoReply              = ICMPv6Type(129)
)

const (
	ICMPv6CodeDestinationNetworkUnreachable = 0
)

type NDOptionType uint8

const (
	NDOptionSourceLinkAddr = NDOptionType(1)
	NDOptionTargetLinkAddr = NDOptionType(2)
)

const (
	IPProtocolICMPV6 = uint8(58)
	IPProtocolUdp    = uint8(17)
)

const (
	EthernetProtocolIPv6 = EthernetType(0x86DD)
)

func TODO() {
	panic("TODO")
}
