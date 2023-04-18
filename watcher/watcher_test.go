package watcher_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cquestor/cc/watcher"
)

func TestWatch(t *testing.T) {
	path, _ := os.Getwd()
	watch, err := watcher.NewWatcher(path)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	watch.AddEvent(watcher.WRITE)
	watch.AddExcludes("watcher_test.go")
	go watch.Watch()
	go func() {
		for {
			select {
			case event := <-watch.Events:
				fmt.Println(event)
			case err := <-watch.Errs:
				panic(err)
			}
		}
	}()
	if err := watch.AddWatch("."); err != nil {
		t.Fatal(err)
	}
	<-done
}
