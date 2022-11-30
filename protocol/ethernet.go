package protocol

import (
	"fmt"
	"net"
)

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

func (etherAddr *EthernetAddr) Equals(addr EthernetAddr) bool {
	for i := 0; i < EthernetAddrLen; i++ {
		if etherAddr.Octet[i] != addr.Octet[i] {
			return false
		}
	}
	return true
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
	ret = ConcatMac(ret, &ether.Header.DstHost)
	ret = ConcatMac(ret, &ether.Header.SrcHost)
	ret = ConcatU16(ret, uint16(ether.Header.Type))
	ret = ConcatBuffer(ret, ether.Payload)
	return ret
}

func (ether *EthernetFrame) Display() {
	fmt.Printf("src: %+v, dst: %+v, type: %x\n",
		ether.Header.SrcHost.String(),
		ether.Header.DstHost.String(),
		ether.Header.Type)
	fmt.Printf("payload: ")
	for i := 0; i < len(ether.Payload.Octet); i++ {
		fmt.Printf("%x ", ether.Payload.Octet[i])
	}
	fmt.Printf("\n")
}

func ParseEtherFrame(data []byte) (*EthernetFrame, error) {
	var err error
	parser := NetParser{
		Buffer: Buffer{
			Octet: data,
		},
		Pointer: 0,
	}
	ether := &EthernetFrame{}
	if ether.Header.DstHost, err = parser.ParseEthernetAddr(); err != nil {
		return nil, err
	}
	if ether.Header.SrcHost, err = parser.ParseEthernetAddr(); err != nil {
		return nil, err
	}
	u16, _ := parser.ParseU16()
	ether.Header.Type = EthernetType(u16)
	if ether.Payload, err = parser.ParseBuffer(); err != nil {
		return nil, err
	}
	return ether, nil
}
