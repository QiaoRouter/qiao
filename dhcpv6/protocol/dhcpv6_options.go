package protocol

import (
	"qiao/protocol"
)

// NdRouterAdvert router advertisement
type NdRouterAdvert struct {
	NdRaHdr        protocol.ICMPv6
	NdRaReachable  uint32 /* reachable time */
	NdRaRetransmit uint32 /* retransmit timer */
	/* could be followed by options */
}

type NdOptMtu struct {
	NdOptMtuType     uint8
	NdOptMtuLen      uint8
	NdOptMtuReserved uint16
	NdOptMtuMtu      uint16
}

// OptionsHdr dhcpv6 options field header
type OptionsHdr struct {
	OptionCode uint16
	OptionLen  uint16
}

// SrcLinkLocalAddr option1 的字段内容
type SrcLinkLocalAddr struct {
	Type          uint8
	Len           uint8
	LinkLocalAddr protocol.EthernetAddr
}

type NdRouterAdvertWithOption struct {
	Ip6              protocol.Ipv6Header
	Icmp6            NdRouterAdvert
	SrcLinkLocalAddr SrcLinkLocalAddr
	NdOptMtu         NdOptMtu
}

// DuidLlt DUID Based on Link-Layer Address Plus Time (DUID-LLT) (in 11.2.)
// 2 + 2 + 4 + 6 = 14B
type DuidLlt struct {
	DuidType      uint16
	HdType        uint16
	Time          uint32
	LinkLayerAddr protocol.EthernetAddr
}

// IaNa Identity Association for Non-temporary Addresses Option (in 21.4.)
type IaNa struct {
	IaId uint32
	T1   uint32
	T2   uint32
}

/*
   uint16_t option_iaaddr;
   uint16_t option_len;
   struct in6_addr ip6addr;
   uint32_t preferred_lifetime;
   uint32_t valid_lifetime;
*/

// IaNaOptions 21.6.  IA Address Option
// 4 + 4 + 16 + 8 + 8 = 40
type IaNaOptions struct {
	OptionIaAddr      uint16
	OptionLen         uint16
	Ip6Addr           protocol.Ipv6Addr
	PreferredLifetime uint32
	ValidLifetime     uint32
}
