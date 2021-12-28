package directory

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/DblK/tinshop/config"
	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/nsp"
	"github.com/DblK/tinshop/repository"
	"github.com/DblK/tinshop/utils"
)

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

func addDirectoryGame(gameFiles []repository.FileDesc, extension string, size int64, path string) []repository.FileDesc {
	var newGameFiles []repository.FileDesc
	newGameFiles = append(newGameFiles, gameFiles...)

	if extension == ".nsp" || extension == ".nsz" {
		newFile := repository.FileDesc{Size: size, Path: path}
		names := utils.ExtractGameID(path)

		if names.ShortID() != "" {
			newFile.GameID = names.ShortID()
			newFile.GameInfo = names.FullID()
			newFile.HostType = repository.LocalFile

			if config.GetConfig().VerifyNSP() {
				valid, errTicket := nspCheck(newFile)
				if valid || (errTicket != nil && errTicket.Error() == "TitleDBKey for game "+newFile.GameID+" is not found") {
					newGameFiles = append(newGameFiles, newFile)
				} else {
					log.Println(errTicket)
				}
			} else {
				newGameFiles = append(newGameFiles, newFile)
			}
		} else {
			log.Println("Ignoring file because parsing failed", path)
		}
	}

	return newGameFiles
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
				extension := filepath.Ext(info.Name())
				newGameFiles = addDirectoryGame(newGameFiles, extension, info.Size(), path)
			} else if info.IsDir() {
				if path != directory {
					watchDirectory(path)
				} else {
					watchDirectory(directory)
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	gameFiles = append(gameFiles, newGameFiles...)

	// Add all files
	if len(newGameFiles) > 0 {
		collection.AddNewGames(newGameFiles)
	}

	return nil
}

func nspCheck(file repository.FileDesc) (bool, error) {
	key, err := collection.GetKey(file.GameID)
	if err != nil {
		return false, err
	}

	f, err := os.Open(file.Path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	valid, err := nsp.IsTicketValid(f, key)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("The ticket in '" + file.Path + "' is not valid!")
	}

	return valid, err
}
