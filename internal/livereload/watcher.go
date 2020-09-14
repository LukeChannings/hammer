package livereload

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

func InitWatcher() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	watcher = w
}

func CloseWatcher() {
	fmt.Println("Closing watcher")
	watcher.Close()
}

func AddWatcher(file string) error {
	return watcher.Add(file)
}
