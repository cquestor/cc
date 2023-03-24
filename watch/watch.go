package watch

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

// Watcher 监听器
type Watcher struct {
	watcher *fsnotify.Watcher
}

// NewWatcher 构造监听器
func NewWatcher() (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		watcher: watcher,
	}, nil
}

func (watcher *Watcher) Add(name string) error {
	return watcher.watcher.Add(name)
}

func (watcher *Watcher) Start() {
	for {
		select {
		case event := <-watcher.watcher.Events:
			fmt.Println("event:", event)
		case err := <-watcher.watcher.Errors:
			fmt.Println("error:", err)
		}
	}
}
