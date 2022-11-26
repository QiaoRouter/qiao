package main

import (
	"qiao/ripng"
)

func main() {
	e := ripng.MakeRipngEngine()
	_ = e.Run()
}
