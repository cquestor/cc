package watch_test

import (
	"testing"

	"github.com/cquestor/cc/watch"
)

func TestWatch(t *testing.T) {
	if w, err := watch.NewWatcher(); err != nil {
		t.Fatal(err)
	} else {
		done := make(chan bool)
		go w.Start()
		w.Add(".")
		<-done
	}
}
