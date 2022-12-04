package hal

import (
	"fmt"
	"net"
	"os"
	"qiao/protocol"
)

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

func disableIpv6(ifName string, open bool) {
	path := fmt.Sprintf("/proc/sys/net/ipv6/conf/%s/disable_ipv6", ifName)
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		panic(err)
	}

	var b []byte
	if open {
		b = []byte{'0'}
	} else {
		b = []byte{'1'}
	}
	n, err := f.Write(b)
	if n != 1 || err != nil {
		panic(fmt.Sprintf("n != 1 || err != nil, err: %+v", err))
	}
	if open {
		fmt.Printf("open %s interface ipv6\n", ifName)
	} else {
		fmt.Printf("disable %s interface ipv6\n", ifName)
	}
}

func macAddr(ifName string) protocol.EthernetAddr {
	netIfs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for i := range netIfs {
		if netIfs[i].Name == ifName {
			mac, err := protocol.ParseMac(netIfs[i].HardwareAddr.String())
			if err != nil {
				return protocol.EthernetAddr{}
			}
			return mac
		}
	}
	return protocol.EthernetAddr{}
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
