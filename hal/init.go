package hal

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/mdlayher/netx/eui64"
	"net"
	"net/netip"
	"os"
	"qiao/config"
	"strings"
	"sync"
)

var once = sync.Once{}

func displayInterfaces() {
	ifs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	fmt.Println("---display interfaces------")
	for i := range ifs {
		fmt.Printf("| interface_name: %+v index: %+v, flags: %+v, mtu: %+v, mac: %+v\n",
			ifs[i].Name, ifs[i].Index, ifs[i].Flags, ifs[i].MTU, ifs[i].HardwareAddr.String())
	}
	fmt.Println("---display interfaces end---")
}

func isInNetInterfaces(ifName string) bool {
	netIfs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for i := range netIfs {
		if netIfs[i].Name == ifName {
			return true
		}
	}
	return false
}

func disableIpv6(ifName string) {
	path := fmt.Sprintf("/proc/sys/net/ipv6/conf/%s/disable_ipv6", ifName)
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		panic(err)
	}
	n, err := f.Write([]byte{'1'})
	if n != 1 || err != nil {
		panic(fmt.Sprintf("n != 1 || err != nil, err: %+v", err))
	}
	fmt.Printf("disable %s interface ipv6\n", ifName)
}

func macAddr(ifName string) net.HardwareAddr {
	netIfs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for i := range netIfs {
		if netIfs[i].Name == ifName {
			return netIfs[i].HardwareAddr
		}
	}
	return nil
}

func ifIndex(ifName string) int {
	netIfs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for i := range netIfs {
		if netIfs[i].Name == ifName {
			return netIfs[i].Index
		}
	}
	return -1
}

func prefix(ipv6 string) string {
	ss := strings.Split(ipv6, "/")
	return ss[0]
}

func ipv6(ifName string) []netip.Addr {
	var ipv6Addr []netip.Addr
	netIfs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for i := range netIfs {
		if netIfs[i].Name == ifName {
			addrs, err := netIfs[i].Addrs()
			if err != nil {
				panic(err)
			}
			for _, addr := range addrs {
				ip, err := netip.ParseAddr(prefix(addr.String()))
				if err != nil {
					continue
				}
				if ip.Is6() {
					ipv6Addr = append(ipv6Addr, ip)
				}
			}
			return ipv6Addr
		}
	}
	return nil
}

func Init() {
	once.Do(func() {
		displayInterfaces()
		ifNames := make([]string, 0)

		ifs, err := net.Interfaces()
		if err != nil {
			panic(err)
		}
		for i := range ifs {
			ifNames = append(ifNames, ifs[i].Name)
		}

		for i := range ifNames {
			handleIn, err := pcap.OpenLive(ifNames[i], config.BufSize, true, pcap.BlockForever)
			if err != nil {
				fmt.Printf("pcap.OpenLive %+v fail, err: %+v\n", ifNames[i], err)
				continue
			}
			handleOut, err := pcap.OpenLive(ifNames[i], config.BufSize, false, pcap.BlockForever)
			if err != nil {
				fmt.Printf("pcap.OpenLive %+v fail, err: %+v\n", ifNames[i], err)
				continue
			}

			ifHandle := &IfHandle{
				IfName:        ifNames[i],
				MAC:           macAddr(ifNames[i]),
				IfIndex:       ifIndex(ifNames[i]),
				PcapHandleIn:  handleIn,
				PcapHandleOut: handleOut,
			}
			//
			// Now we only support ipv6 over Ethernet,
			// please refer rfc4291 and rfc2464 for more details
			// about eui64 computing and link-local address of ipv6 over Ethernet.
			//
			ip, err := eui64.ParseMAC(net.ParseIP("fe80::"), ifHandle.MAC)
			if err != nil {
				fmt.Printf("if %+v: eui64.ParseMAC %+v fail, err is %+v\n\n", ifHandle.IfName,
					ifHandle.MAC, err)
				continue
			}
			ifHandle.IPv6 = ipv6(ifNames[i])
			fmt.Printf("%+v ipv6 addrs are %+v\n", ifHandle.IfName, ifHandle.IPv6)
			ifHandle.LinkLocalIPv6, err = netip.ParseAddr(ip.String())
			if err != nil {
				panic(err)
			}
			ifHandle.PacketSource = gopacket.NewPacketSource(handleIn, handleIn.LinkType())
			if config.Experimental {
				disableIpv6(ifNames[i])
			}
			fmt.Printf("%+v mac is %+v\n", ifHandle.IfName, ifHandle.MAC)
			fmt.Printf("%+v link-local addr is %+v\n", ifHandle.IfName, ifHandle.LinkLocalIPv6)
			IfHandles = append(IfHandles, ifHandle)
			fmt.Printf("hal: pcap capture on interface %+v\n\n", ifNames[i])
		}
	})
}
