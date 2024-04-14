package utils

import "time"

// NewTimer produce a new async ticker.
//
// d is for the duration between two ticks
// and fn is the functions called at each tick: it takes a chan in parameter, and you can put anything here to disable the
// ticker
func NewTimer(d time.Duration, fn func(chan struct{})) {
	ticker := time.NewTicker(d)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fn(quit)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
