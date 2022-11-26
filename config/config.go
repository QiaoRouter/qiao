package config

import "time"

const (
	DEBUG        = true
	Experimental = true

	BufSize = 8192
	Backoff = time.Second * 5
)
