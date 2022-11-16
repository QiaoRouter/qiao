package main

import (
	"fmt"
	"qiao/hal"
)

func main() {
	hal.Init()
	fmt.Println("Hello, world!")
	for i := 0; i < len(hal.IfHandles); i++ {
		hal.IfHandles[i].GetNeighborMacAddr()
	}
}
