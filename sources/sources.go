// @title tinshop Sources

// @BasePath /sources/

// Package sources provides management of various sources
package sources

import (
	"log"
	"net/http"

	"github.com/DblK/tinshop/repository"
	"github.com/DblK/tinshop/sources/directory"
	"github.com/DblK/tinshop/sources/nfs"
	"github.com/DblK/tinshop/utils"
)

// SourceProvider stores all sources available
type SourceProvider struct {
	Directory repository.Source
	NFS       repository.Source
}

var sourcesProvider SourceProvider

// OnConfigUpdate from all sources
func OnConfigUpdate(cfg repository.Config) {
	log.Println("Sources loading...")

	// Directories
	srcDirectories := directory.New()
	srcDirectories.Reset()
	srcDirectories.Load(cfg.Directories(), len(cfg.NfsShares()) == 0)
	sourcesProvider.Directory = srcDirectories

	// NFS
	srcNFS := nfs.New()
	srcNFS.Reset()
	srcNFS.Load(cfg.NfsShares(), false)
	sourcesProvider.NFS = srcNFS
}

// BeforeConfigUpdate from all sources
func BeforeConfigUpdate(cfg repository.Config) {
	if sourcesProvider.Directory != nil {
		sourcesProvider.Directory.UnWatchAll()
	}
	if sourcesProvider.NFS != nil {
		sourcesProvider.NFS.UnWatchAll()
	}
}

// GetFiles returns all games files in various sources
func GetFiles() []repository.FileDesc {
	mergedGameFiles := make([]repository.FileDesc, 0)
	if sourcesProvider.Directory != nil {
		mergedGameFiles = append(mergedGameFiles, sourcesProvider.Directory.GetFiles()...)
	}
	if sourcesProvider.NFS != nil {
		mergedGameFiles = append(mergedGameFiles, sourcesProvider.NFS.GetFiles()...)
	}
	return mergedGameFiles
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
		sourcesProvider.Directory.Download(w, r, gameID, GetFiles()[idx].Path)
	case repository.NFSShare:
		sourcesProvider.NFS.Download(w, r, gameID, GetFiles()[idx].Path)

	default:
		w.WriteHeader(http.StatusNotImplemented)
		log.Printf("The type '%s' is not implemented to download game", GetFiles()[idx].HostType)
	}
}
