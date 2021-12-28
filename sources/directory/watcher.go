package directory

import (
	"log"
	"path/filepath"
	"sync"

	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/repository"
	"gopkg.in/fsnotify.v1"
)

var watcherDirectories *fsnotify.Watcher

func newWatcher() *fsnotify.Watcher {
	watcherDirectories, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	return watcherDirectories
}

func watchDirectory(directory string) {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		defer watcherDirectories.Close()

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcherDirectories.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()
						return
					}

					if event.Op&fsnotify.Create != 0 {
						newGames := addDirectoryGame(make([]repository.FileDesc, 0), filepath.Ext(event.Name), 0, event.Name)
						gameFiles = append(gameFiles, newGames...)
						collection.AddNewGames(newGames)
					} else if event.Op&fsnotify.Remove != 0 {
						removeEntriesFromDirectory(event.Name)
					} else if event.Op&fsnotify.Rename == fsnotify.Rename {
						removeEntriesFromDirectory(event.Name)
					}

				case err, ok := <-watcherDirectories.Errors:
					if ok { // 'Errors' channel is not closed
						log.Printf("watcher error: %v\n", err)
					}
					eventsWG.Done()
					return
				}
			}
		}()
		errWatcher := watcherDirectories.Add(directory)
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
		if errWatcher != nil {
			eventsWG.Done()
		}
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
}
