package protocol

import "net"

type EthernetAddr struct {
	Octet [EthernetAddrLen]byte
}

func ParseMac(s string) (EthernetAddr, error) {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return EthernetAddr{}, err
	}
	etherAddr := EthernetAddr{}
	for i := 0; i < len(mac); i++ {
		etherAddr.Octet[i] = mac[i]
	}
	return etherAddr, nil
}

func (etherAddr *EthernetAddr) String() string {
	return etherAddr.NetHardwareAddr().String()
}

func (etherAddr *EthernetAddr) Serialize() Buffer {
	buf := Buffer{
		Octet: make([]byte, EthernetAddrLen),
	}
	for i := 0; i < EthernetAddrLen; i++ {
		buf.Octet[i] = etherAddr.Octet[i]
	}
	return buf
}

func (etherAddr *EthernetAddr) NetHardwareAddr() net.HardwareAddr {
	mac := net.HardwareAddr{}
	for i := 0; i < len(etherAddr.Octet); i++ {
		mac[i] = etherAddr.Octet[i]
	}
	return mac
}

type EthernetHeader struct {
	DstHost EthernetAddr
	SrcHost EthernetAddr
	Type    EthernetType
}

type EthernetFrame struct {
	Header  EthernetHeader
	Payload Buffer
}

func (ether *EthernetHeader) String() string {
	return ""
}

func (frame *EthernetFrame) Serialize() []byte {
	return []byte{}
}
