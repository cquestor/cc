package watcher

import (
	"sync"
	"time"
)

// Debounce 防抖
func Debounce(f func(), d time.Duration) func() {
	var lock sync.Mutex
	var timer *time.Timer
	return func() {
		lock.Lock()
		defer lock.Unlock()
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(d, f)
	}
}
