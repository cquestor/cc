package watcher

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// Watcher 文件监控
type Watcher struct {
	watcher        *fsnotify.Watcher
	BasePath       string
	interestEvents []fsnotify.Op
	excludes       []string
	includes       []string
	Events         chan fsnotify.Event
	Errs           chan error
}

const (
	WRITE  fsnotify.Op = fsnotify.Write
	RENAME fsnotify.Op = fsnotify.Rename
	REMOVE fsnotify.Op = fsnotify.Remove
	CHMOD  fsnotify.Op = fsnotify.Chmod
	CREATE fsnotify.Op = fsnotify.Create
)

// NewWactcher 构造文件监控
func NewWatcher(basePath string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		Events:   make(chan fsnotify.Event),
		Errs:     make(chan error),
		BasePath: basePath,
		watcher:  watcher,
		includes: make([]string, 0),
		excludes: make([]string, 0),
	}, nil
}

// Init 初始化监听
func (watcher *Watcher) Init() error {
	return filepath.Walk(watcher.BasePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if watcher.inIncludes(path) {
				return watcher.AddWatch(path)
			}
			if strings.HasPrefix(filepath.Base(path), ".") {
				return nil
			}
			if watcher.inExcludes(path) {
				return nil
			}
			return watcher.AddWatch(path)
		}
		return nil
	})
}

// Watch 开始监听
func (watcher *Watcher) Watch() {
	for {
		select {
		case event := <-watcher.watcher.Events:
			if watcher.isInterested(event) {
				watcher.Events <- event
			}
		case err := <-watcher.watcher.Errors:
			watcher.Errs <- err
		}
	}
}

// AddEvent 添加监听事件
func (watcher *Watcher) AddEvent(events ...fsnotify.Op) {
	watcher.interestEvents = append(watcher.interestEvents, events...)
}

// AddWatch 添加监听文件
func (watcher *Watcher) AddWatch(name string) error {
	stat, err := os.Stat(name)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return nil
	}
	if watcher.inIncludes(name) {
		return watcher.watcher.Add(name)
	}
	if strings.HasPrefix(filepath.Base(name), ".") {
		return watcher.watcher.Add(name)
	}
	if watcher.inExcludes(name) {
		return nil
	}
	return watcher.watcher.Add(name)
}

// RemoveWatch 删除监听
func (watcher *Watcher) RemoveWatch(name string) error {
	if watcher.inIncludes(name) {
		return watcher.watcher.Remove(name)
	}
	if strings.HasPrefix(filepath.Base(name), ".") {
		return watcher.watcher.Remove(name)
	}
	if watcher.inExcludes(name) {
		return nil
	}
	return watcher.watcher.Remove(name)
}

// AddIncludes 添加包含文件
func (watcher *Watcher) AddIncludes(names ...string) {
	for _, path := range names {
		watcher.includes = append(watcher.includes, filepath.Join(watcher.BasePath, path))
	}
}

// AddExcludes 添加排除文件
func (watcher *Watcher) AddExcludes(names ...string) {
	for _, path := range names {
		watcher.excludes = append(watcher.excludes, filepath.Join(watcher.BasePath, path))
	}
}

// Close 关闭文件监控
func (watcher *Watcher) Close() {
	watcher.watcher.Close()
}

// isInterested 是否是需要上报的事件
func (watcher *Watcher) isInterested(event fsnotify.Event) bool {
	if !watcher.inEvents(event) {
		return false
	}
	filename := filepath.Base(event.Name)
	if watcher.inIncludes(event.Name) {
		return true
	}
	if strings.HasPrefix(filename, ".") {
		return false
	}
	if watcher.inExcludes(event.Name) {
		return false
	}
	return true
}

// inEvents 是否是要监听的事件
func (watcher *Watcher) inEvents(event fsnotify.Event) bool {
	for _, op := range watcher.interestEvents {
		if op == event.Op {
			return true
		}
	}
	return false
}

// inIncludes 是否是声明包含的文件
func (watcher *Watcher) inIncludes(name string) bool {
	for _, path := range watcher.includes {
		if path == name || strings.HasPrefix(name, path) {
			return true
		}
	}
	return false
}

// inExcludes 是否是声明不包含的文件
func (watcher *Watcher) inExcludes(name string) bool {
	for _, path := range watcher.excludes {
		if path == name || strings.HasPrefix(name, path) {
			return true
		}
	}
	return false
}
