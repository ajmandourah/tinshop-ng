package sources

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	collection "github.com/dblk/tinshop/gamescollection"
	"github.com/dblk/tinshop/repository"
	"github.com/dblk/tinshop/utils"
)

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

func loadGamesDirectory(directory string) error {
	log.Printf("Loading games from directory '%s'...\n", directory)
	var newGameFiles []repository.FileDesc
	// Walk through games directory
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				newFile := repository.FileDesc{Size: info.Size(), Path: path}
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
