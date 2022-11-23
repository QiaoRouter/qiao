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
		mac = append(mac, etherAddr.Octet[i])
	}
	return mac
}

func (etherAddr *EthernetAddr) AllZero() bool {
	for i := 0; i < len(etherAddr.Octet); i++ {
		if etherAddr.Octet[i] != 0 {
			return false
		}
	}
	return true
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

func (ether *EthernetFrame) Serialize() []byte {
	var ret []byte
	ret = concatMac(ret, &ether.Header.DstHost)
	ret = concatMac(ret, &ether.Header.SrcHost)
	ret = concatU16(ret, uint16(ether.Header.Type))
	ret = concatBuffer(ret, ether.Payload)
	return ret
}
