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

type allSources struct {
	sourcesProvider SourceProvider
	collection      repository.Collection
}

// New create a new collection
func New(collection repository.Collection) repository.Sources {
	return &allSources{
		collection: collection,
	}
}

// OnConfigUpdate from all sources
func (s *allSources) OnConfigUpdate(cfg repository.Config) {
	log.Println("Sources loading...")

	// Directories
	srcDirectories := directory.New(s.collection, cfg)
	srcDirectories.Reset()
	srcDirectories.Load(cfg.Directories(), len(cfg.NfsShares()) == 0)
	s.sourcesProvider.Directory = srcDirectories

	// NFS
	srcNFS := nfs.New(s.collection, cfg)
	srcNFS.Reset()
	srcNFS.Load(cfg.NfsShares(), false)
	s.sourcesProvider.NFS = srcNFS
}

// BeforeConfigUpdate from all sources
func (s *allSources) BeforeConfigUpdate(cfg repository.Config) {
	if s.sourcesProvider.Directory != nil {
		s.sourcesProvider.Directory.UnWatchAll()
	}
	if s.sourcesProvider.NFS != nil {
		s.sourcesProvider.NFS.UnWatchAll()
	}
}

// GetFiles returns all games files in various sources
func (s *allSources) GetFiles() []repository.FileDesc {
	mergedGameFiles := make([]repository.FileDesc, 0)
	if s.sourcesProvider.Directory != nil {
		mergedGameFiles = append(mergedGameFiles, s.sourcesProvider.Directory.GetFiles()...)
	}
	if s.sourcesProvider.NFS != nil {
		mergedGameFiles = append(mergedGameFiles, s.sourcesProvider.NFS.GetFiles()...)
	}
	return mergedGameFiles
}

func (s *allSources) HasGame(gameID string) bool {
	idx := utils.Search(len(s.GetFiles()), func(index int) bool {
		return s.GetFiles()[index].GameID == gameID
	})
	return idx != -1
}

// DownloadGame method provide the file based on the source storage
func (s *allSources) DownloadGame(gameID string, w http.ResponseWriter, r *http.Request) {
	idx := utils.Search(len(s.GetFiles()), func(index int) bool {
		return s.GetFiles()[index].GameID == gameID
	})

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Game '%s' not found!", gameID)
		return
	}
	log.Println("Retrieving from location '" + s.GetFiles()[idx].Path + "'")
	switch s.GetFiles()[idx].HostType {
	case repository.LocalFile:
		s.sourcesProvider.Directory.Download(w, r, gameID, s.GetFiles()[idx].Path)
	case repository.NFSShare:
		s.sourcesProvider.NFS.Download(w, r, gameID, s.GetFiles()[idx].Path)

	default:
		w.WriteHeader(http.StatusNotImplemented)
		log.Printf("The type '%s' is not implemented to download game", s.GetFiles()[idx].HostType)
	}
}
