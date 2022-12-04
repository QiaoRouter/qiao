package hal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

		ifNameToIpv6 := make(map[string]*IfHandle)
		path := "./if_name_to_ipv6"
		j, err := ioutil.ReadFile(path)
		if err == nil {
			_ = json.Unmarshal(j, &ifNameToIpv6)
		}
		for i := range ifs {
			disableIpv6(ifs[i].Name, true)
			ifHandle := &IfHandle{
				IfName: ifs[i].Name,
			}
			err = ifHandle.Init()
			if err != nil {
				fmt.Printf("ifHandle.Init() fail, err is %+v\n", err)
				continue
			}
			if ifHandle.IPv6.AllZero() {
				// 如果找不到ipv6地址，也许是因为之前关闭过
				// 可以去持久化文件里找
				ifHandle.IPv6 = ifNameToIpv6[ifs[i].Name].IPv6
				ifHandle.IPv6Mask = ifNameToIpv6[ifs[i].Name].IPv6Mask
			} else {
				ifNameToIpv6[ifHandle.IfName] = ifHandle
			}
			fmt.Printf("interface %+v, ipv6 addr is %+v\n", ifHandle.IfName, ifHandle.IPv6.String())

			IfHandles = append(IfHandles, ifHandle)
			disableIpv6(ifs[i].Name, false)
		}
		j, _ = json.Marshal(ifNameToIpv6)
		err = ioutil.WriteFile(path, j, 0777)
		if err != nil {
			panic(err)
		}
		go ticker()
	})
}

func Close() {
	for _, h := range IfHandles {
		disableIpv6(h.IfName, true)
	}
}
