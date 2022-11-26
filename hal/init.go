package hal

import (
	"fmt"
	"net"
	"sync"
)

var once = sync.Once{}

func Init() {
	once.Do(func() {
		displayInterfaces()
		ifs, err := net.Interfaces()
		if err != nil {
			panic(err)
		}
		for i := range ifs {
			// disableIpv6(ifs[i].Name, true)
			ifHandle := &IfHandle{
				IfName: ifs[i].Name,
			}
			err = ifHandle.Init()
			if err != nil {
				fmt.Printf("ifHandle.Init() fail, err is %+v\n", err)
				continue
			}
			IfHandles = append(IfHandles, ifHandle)
		}
	})
}
