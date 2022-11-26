package protocol

const (
	EthernetAddrLen   = 6  /* in unit of 1 octet */
	IPv6AddrLen       = 16 /* in unit of 1 octet */
	IPv6Version       = uint8(6 << 4)
	EthernetHeaderLen = 14
)

type EthernetType uint16

type ICMPv6Type uint8

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
