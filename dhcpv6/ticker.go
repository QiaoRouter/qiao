package dhcpv6

import (
	"qiao/config"
	"time"
)

func (e *Engine) ticker() {
	t := time.Now()
	for true {
		t1 := time.Now()
		if t1.Sub(t) > config.Backoff {
			t = t1
			go e.broadcast()
		}
		time.Sleep(time.Millisecond * 100)
	}
}
