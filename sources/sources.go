// @title tinshop Sources

// @BasePath /sources/

// Package sources provides management of various sources
package sources

import (
	"log"
	"net/http"

	"github.com/DblK/tinshop/repository"
	"github.com/DblK/tinshop/utils"
)

var gameFiles []repository.FileDesc

// OnConfigUpdate from all sources
func OnConfigUpdate(cfg repository.Config) {
	log.Println("Sources loading...")
	gameFiles = make([]repository.FileDesc, 0)
	watcherDirectories = newWatcher()
	loadGamesDirectories(cfg.Directories(), len(cfg.NfsShares()) == 0)
	loadGamesNfsShares(cfg.NfsShares())
}

// BeforeConfigUpdate from all sources
func BeforeConfigUpdate(cfg repository.Config) {
	removeGamesWatcherDirectories()
}

// GetFiles returns all games files in various sources
func GetFiles() []repository.FileDesc {
	return gameFiles
}

// AddFiles add files to global sources
func AddFiles(files []repository.FileDesc) {
	gameFiles = append(gameFiles, files...)
}

// DownloadGame method provide the file based on the source storage
func DownloadGame(gameID string, w http.ResponseWriter, r *http.Request) {
	idx := utils.Search(len(GetFiles()), func(index int) bool {
		return GetFiles()[index].GameID == gameID
	})

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Game '%s' not found!", gameID)
		return
	}
	log.Println("Retrieving from location '" + GetFiles()[idx].Path + "'")
	switch GetFiles()[idx].HostType {
	case repository.LocalFile:
		downloadLocalFile(w, r, gameID, GetFiles()[idx].Path)
	case repository.NFSShare:
		downloadNfsFile(w, r, GetFiles()[idx].Path)

	default:
		w.WriteHeader(http.StatusNotImplemented)
		log.Printf("The type '%s' is not implemented to download game", GetFiles()[idx].HostType)
	}
}
