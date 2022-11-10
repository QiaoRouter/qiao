package hal

import (
	"fmt"
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

func isInNetInterfaces(netIf string) bool {
	netIfs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for i := range netIfs {
		if netIfs[i].Name == netIf {
			return true
		}
	}
	return false
}

func disableIpv6(netIf string) {
	path := fmt.Sprintf("/proc/sys/net/ipv6/conf/%s/disable_ipv6", netIf)
	f, err := os.OpenFile(path, os.O_CREATE, 0)
	if err != nil {
		panic(err)
	}
	n, err := f.Write([]byte{1})
	if n != 1 || err != nil {
		panic(fmt.Sprintf("n != 1 || err != nil, err: %+v", err))
	}
}

func Init() {
	once.Do(func() {
		displayInterfaces()
		ifNames := experimentInterfaces[host]
		for i := range ifNames {
			// if exist in net interfaces
			if !isInNetInterfaces(ifNames[i]) {
				continue
			}
			if config.Experimental {
				disableIpv6(ifNames[i])
			}
		}
	})
}
