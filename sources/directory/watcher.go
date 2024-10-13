package directory

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/ajmandourah/tinshop-ng/repository"
	"gopkg.in/fsnotify.v1"
)

func (src *directorySource) newWatcher() *fsnotify.Watcher {
	watcherDirectories, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	return watcherDirectories
}

func (src *directorySource) watchDirectory(directory string) {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		defer src.watcherDirectories.Close()

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-src.watcherDirectories.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()
						return
					}

					if event.Op&fsnotify.Create != 0 {
						newGames := src.addDirectoryGame(make([]repository.FileDesc, 0), filepath.Ext(event.Name), 0, event.Name)
						src.gameFiles = append(src.gameFiles, newGames...)
						src.collection.AddNewGames(newGames)
					} else if event.Op&fsnotify.Remove != 0 {
						src.removeEntriesFromDirectory(event.Name)
					} else if event.Op&fsnotify.Rename == fsnotify.Rename {
						src.removeEntriesFromDirectory(event.Name)
					}

				case err, ok := <-src.watcherDirectories.Errors:
					if ok { // 'Errors' channel is not closed
						log.Printf("watcher error: %v\n", err)
					}
					eventsWG.Done()
					return
				}
			}
		}()
		errWatcher := src.watcherDirectories.Add(directory)
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
		if errWatcher != nil {
			eventsWG.Done()
		}
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
}
