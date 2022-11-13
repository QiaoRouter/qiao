package hal

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"net"
	"os"
	"qiao/config"
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

func Init() {
	once.Do(func() {
		displayInterfaces()
		ifNames := experimentInterfaces[Host]
		for i := range ifNames {
			// if exist in net interfaces
			if !isInNetInterfaces(ifNames[i]) {
				continue
			}
			if config.Experimental {
				disableIpv6(ifNames[i])
			}
		}

		if !config.Experimental {
			ifs, err := net.Interfaces()
			if err != nil {
				panic(err)
			}
			for i := range ifs {
				ifNames = append(ifNames, ifs[i].Name)
			}
		}

		for i := range ifNames {
			handle, err := pcap.OpenLive(ifNames[i], config.BufSize, true, pcap.BlockForever)
			if err != nil {
				fmt.Printf("pcap.OpenLive %+v fail, err: %+v\n", ifNames[i], err)
				continue
			}

			ifHandle := &IfHandle{
				IfName:       ifNames[i],
				PcapHandle:   handle,
				PacketSource: gopacket.NewPacketSource(handle, handle.LinkType()),
				MAC:          macAddr(ifNames[i]),
			}
			fmt.Printf("%+v mac is %+v\n", ifHandle.IfName, ifHandle.MAC)
			IfHandles = append(IfHandles, ifHandle)
			fmt.Printf("hal: pcap capture on interface %+v\n", ifNames[i])
		}

	})
}
