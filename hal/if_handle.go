package hal

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/mdlayher/netx/eui64"
	"net/netip"
	"qiao/config"
	"qiao/protocol"
	"strconv"
	"strings"

	"net"
)

type HostName string
type IfNames []string

type IfHandle struct {
	IfName          string
	IfIndex         int
	PcapHandleIn    *pcap.Handle
	PcapHandleOut   *pcap.Handle
	PacketSource    *gopacket.PacketSource
	MAC             protocol.EthernetAddr
	LinkLocalIPv6   protocol.Ipv6Addr
	NeighborMacAddr net.HardwareAddr
	IPv6            protocol.Ipv6Addr
	IPv6Mask        int
}

func (h *IfHandle) Init() error {
	h.MAC = macAddr(h.IfName)
	if h.MAC.AllZero() {
		return errors.New(fmt.Sprintf("%+v mac all zero", h.IfName))
	}
	h.IfIndex = ifIndex(h.IfName)
	handleIn, err := pcap.OpenLive(h.IfName, config.BufSize, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	handleOut, err := pcap.OpenLive(h.IfName, config.BufSize, false, pcap.BlockForever)
	if err != nil {
		return err
	}
	h.PcapHandleIn = handleIn
	h.PcapHandleOut = handleOut
	//
	// Now we only support ipv6 over Ethernet,
	// please refer rfc4291 and rfc2464 for more
	// details about eui64 computing and link-local
	// address of ipv6 over Ethernet.
	//
	ip, err := eui64.ParseMAC(net.ParseIP("fe80::"), h.MAC.NetHardwareAddr())
	if err != nil {
		return err
	}
	err = h.intIpv6()
	if err != nil {
		return err
	}
	h.LinkLocalIPv6, err = protocol.ParseIpv6(ip.String())
	if err != nil {
		panic(err)
	}
	h.PacketSource = gopacket.NewPacketSource(handleIn, handleIn.LinkType())

	//if config.Experimental {
	//	// 如果是实验的话，关闭该网口的ipv6功能
	//	// 确保linux不会自行处理ipv6数据包
	//	disableIpv6(h.IfName, false)
	//}
	return nil
}

var IfHandles []*IfHandle

func (h *IfHandle) intIpv6() error {
	netIfs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for i := range netIfs {
		if netIfs[i].Name == h.IfName {
			addrs, err := netIfs[i].Addrs()
			if err != nil {
				panic(err)
			}
			for _, addr := range addrs {
				fmt.Printf("addr: %+v\n", addr)
				if strings.HasPrefix(addr.String(), "fe80") {
					continue
				}
				ss := strings.Split(addr.String(), "/")
				prefix := ss[0]
				mask, err := strconv.Atoi(ss[1])
				if err != nil {
					return err
				}
				ip, err := netip.ParseAddr(prefix)
				if err != nil {
					return err
				}
				if ip.Is6() {
					in6, err := protocol.ParseIpv6(ip.String())
					if err != nil {
						return err
					}
					h.IPv6 = in6
					h.IPv6Mask = mask
				}
			}
		}
	}
	return nil
}

func (h *IfHandle) SendIpv6(ipv6Dgm *protocol.Ipv6Datagram, macDst protocol.EthernetAddr) {
	frame := ipv6Dgm.ToEthernetFrame(h.MAC, macDst)
	err := h.PcapHandleOut.WritePacketData(frame.Serialize())
	if err != nil {
		panic(err)
	}
}
