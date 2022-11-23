package main

import (
	"fmt"
	"qiao/hal"
)

func main() {
	hal.Init()
	fmt.Printf("we have %d if_handles\n", len(hal.IfHandles))
}
