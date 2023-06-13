package dhcpv6

import "qiao/protocol"

type DHCPv6Packet struct {
	Header      DHCPv6Header
	MessageBody protocol.Buffer
}

type DHCPv6Header struct {
	MsgType         uint8
	TransactionIdHi uint8
	TransactionIdLo uint8
}
