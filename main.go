package main

import (
	"fmt"
	"os"
	"qiao/dhcpv6"
	"qiao/ripng"
)

func main() {
	// 获取命令行参数
	args := os.Args

	if len(args) < 2 {
		fmt.Println("usage: ./qiao [options]\n" +
			"Options:\n" +
			"--ripng          run ripng server\n" +
			"--dhcpv6         run dhcpv6 server")
		os.Exit(1)
	} else {
		if args[1] == "--ripng" {
			e := ripng.MakeRipngEngine()
			_ = e.Run()
		} else if args[1] == "--dhcpv6" {
			fmt.Println("dhcpv6 server is running!")
			e := dhcpv6.MakeDhcpv6Engine()
			_ = e.Run()
		}
	}

}
