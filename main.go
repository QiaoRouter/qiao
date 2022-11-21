package main

import (
	"fmt"
	"qiao/hal"
)

func main() {
	hal.Init()
	fmt.Printf("we have %d if_handles\n", len(hal.IfHandles))

	for i := 0; i < len(hal.IfHandles); i++ {
		hal.IfHandles[i].GetNeighborMacAddr()
	}
}
