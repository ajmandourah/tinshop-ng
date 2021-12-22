package sources

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/repository"
	"github.com/DblK/tinshop/utils"
	"gopkg.in/fsnotify.v1"
)

var watcherDirectories *fsnotify.Watcher

func loadGamesDirectories(directories []string, singleSource bool) {
	for _, directory := range directories {
		err := loadGamesDirectory(directory)

		if err != nil {
			if len(directories) == 1 && err.Error() == "lstat ./games: no such file or directory" && singleSource {
				log.Fatal("You must create a folder 'games' and put your games inside or use config.yml to add sources!")
			} else {
				log.Println(err)
			}
		}
	}
}

func removeGamesWatcherDirectories() {
	log.Println("Removing watcher from all directories")
	if watcherDirectories != nil {
		watcherDirectories.Close()
	}
}

func removeEntriesFromDirectory(directory string) {
	log.Println("removeEntriesFromDirectory", directory)
	for index, game := range gameFiles {
		if game.HostType == repository.LocalFile && strings.Contains(game.Path, directory) {
			// Need to remove game
			gameFiles = utils.RemoveFileDesc(gameFiles, index)

			// Stop watching of directories
			if directory == filepath.Dir(directory) {
				_ = watcherDirectories.Remove(filepath.Dir(game.Path))
			}

			// Remove entry from collection
			collection.RemoveGame(game.GameID)
		}
	}
}

func AddDirectoryGame(gameFiles []repository.FileDesc, extension string, size int64, path string) []repository.FileDesc {
	var newGameFiles []repository.FileDesc
	newGameFiles = append(newGameFiles, gameFiles...)

	if extension == ".nsp" || extension == ".nsz" {
		newFile := repository.FileDesc{Size: size, Path: path}
		names := utils.ExtractGameID(path)

		if names.ShortID() != "" {
			newFile.GameID = names.ShortID()
			newFile.GameInfo = names.FullID()
			newFile.HostType = repository.LocalFile
			newGameFiles = append(newGameFiles, newFile)
		} else {
			log.Println("Ignoring file because parsing failed", path)
		}
	}

	return newGameFiles
}

func loadGamesDirectory(directory string) error {
	log.Printf("Loading games from directory '%s'...\n", directory)

	// Add watcher for directories
	watchDirectory(directory)

	var newGameFiles []repository.FileDesc
	// Walk through games directory
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				extension := filepath.Ext(info.Name())
				newGameFiles = AddDirectoryGame(newGameFiles, extension, info.Size(), path)
			} else if info.IsDir() && path != directory {
				watchDirectory(path)
			}
			return nil
		})
	if err != nil {
		return err
	}
	AddFiles(newGameFiles)

	// Add all files
	if len(newGameFiles) > 0 {
		collection.AddNewGames(newGameFiles)
	}

	return nil
}

func downloadLocalFile(w http.ResponseWriter, r *http.Request, game, path string) {
	f, err := os.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	fi, err := f.Stat()

	if err == nil {
		http.ServeContent(w, r, game, fi.ModTime(), f)
	} else {
		http.ServeContent(w, r, game, time.Time{}, f)
	}
}

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
						newGames := AddDirectoryGame(make([]repository.FileDesc, 0), filepath.Ext(event.Name), 0, event.Name)
						AddFiles(newGames)
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
