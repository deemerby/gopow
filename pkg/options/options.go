package options

import "time"

type AppOptions struct {
	MaxNumber    int64
	ZeroCnt      int
	MaxIteration int
	ConDeadline  time.Duration
}
