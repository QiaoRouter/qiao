package ripng

import (
	"fmt"
	"qiao/hal"
	"qiao/protocol"
)

func (e *Engine) receivePacketAndHandleIt(h *hal.IfHandle) {
	for true {
		p, err := h.PacketSource.NextPacket()
		if err != nil {
			fmt.Printf("h.PacketSource.NextPacket fail, err: %+v\n", err)
			continue
		}
		ether, err := protocol.ParseEtherFrame(p.Data())
		if err != nil {
			continue
		}
		if ether.Header.SrcHost.Equals(h.MAC) {
			continue
		}
		if ether.Header.Type != protocol.EthernetProtocolIPv6 {
			continue
		}
		ipv6Dgrm, err := protocol.ParseIpv6Datagram(ether.Payload)
		if err != nil {
			continue
		}
		if !ipv6Dgrm.ChecksumValid() {
			continue
		}
		go e.HandleIpv6(ipv6Dgrm)
	}
}

func (e *Engine) HandleIpv6(dgrm *protocol.Ipv6Datagram) {

}
