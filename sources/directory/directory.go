package directory

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DblK/tinshop/repository"
)

var gameFiles []repository.FileDesc

type directorySource struct {
}

// New create a directory source
func New() repository.Source {
	gameFiles = make([]repository.FileDesc, 0)
	return &directorySource{}
}

func (src *directorySource) Download(w http.ResponseWriter, r *http.Request, game, path string) {
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

func (src *directorySource) Load(directories []string, unique bool) {
	for _, directory := range directories {
		err := loadGamesDirectory(directory)

		if err != nil {
			if len(directories) == 1 && err.Error() == "lstat ./games: no such file or directory" && unique {
				log.Fatal("You must create a folder 'games' and put your games inside or use config.yml to add sources!")
			} else {
				log.Println(err)
			}
		}
	}
}

func (src *directorySource) Reset() {
	watcherDirectories = newWatcher()
	gameFiles = make([]repository.FileDesc, 0)
}

func (src *directorySource) UnWatchAll() {
	removeGamesWatcherDirectories()
}

func (src *directorySource) GetFiles() []repository.FileDesc {
	return gameFiles
}
