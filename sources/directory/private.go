package directory

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	collection "github.com/ajmandourah/tinshop/gamescollection"
	"github.com/ajmandourah/tinshop/nsp"
	"github.com/ajmandourah/tinshop/repository"
	"github.com/ajmandourah/tinshop/utils"
	"github.com/charlievieth/fastwalk"
)

func (src *directorySource) removeGamesWatcherDirectories() {
	log.Println("Removing watcher from all directories")
	if src.watcherDirectories != nil {
		src.watcherDirectories.Close()
	}
}

func (src *directorySource) removeEntriesFromDirectory(directory string) {
	log.Println("removeEntriesFromDirectory", directory)
	for index, game := range src.gameFiles {
		if game.HostType == repository.LocalFile && strings.Contains(game.Path, directory) {
			// Need to remove game
			src.gameFiles = utils.RemoveFileDesc(src.gameFiles, index)

			// Stop watching of directories
			if directory == filepath.Dir(directory) {
				_ = src.watcherDirectories.Remove(filepath.Dir(game.Path))
			}

			// Remove entry from collection
			src.collection.RemoveGame(game.GameID)
		}
	}
}

func (src *directorySource) addDirectoryGame(gameFiles []repository.FileDesc, extension string, size int64, path string) []repository.FileDesc {
	var newGameFiles []repository.FileDesc
	newGameFiles = append(newGameFiles, gameFiles...)

	if extension == ".nsp" || extension == ".nsz" || extension == ".xci" {
		newFile := repository.FileDesc{Size: size, Path: path}
		names, decrypted := utils.ExtractGameID(path)
		//Rename the file if decrypted and option is enabled
		if decrypted {
			if collection.Rename {
				title := src.collection.GenTitle(names.ShortID())
				newName := filepath.Join(filepath.Dir(path), title+extension)

				os.Rename(path, newName)
				log.Println("renamed: ", filepath.Base(path), " to ", filepath.Base(newName))
			}
		}
		if names.ShortID() != "" {
			newFile.GameID = names.ShortID()
			newFile.GameInfo = names.FullID()
			newFile.HostType = repository.LocalFile
			newFile.Extension = names.Extension()

			if src.config.VerifyNSP() {
				valid, errTicket := src.nspCheck(newFile)
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

func (src *directorySource) loadGamesDirectory(directory string) error {
	log.Printf("Loading games from directory '%s'...\n", directory)

	var newGameFiles []repository.FileDesc

	conf := fastwalk.Config{
		Sort: fastwalk.SortFilesFirst,
	}

	// Walk through games directory
	err := fastwalk.Walk(&conf, directory,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				extension := filepath.Ext(info.Name())
				fileInfo, err := info.Info()
				if err != nil {
					return err
				}

				newGameFiles = src.addDirectoryGame(newGameFiles, extension, fileInfo.Size(), path)

			} else if info.IsDir() {
				if path != directory {
					src.watchDirectory(path)
					return fastwalk.SkipDir
				} else {
					src.watchDirectory(directory)
				}
			}
			return nil
		})

	if err != nil {
		return err
	}
	src.gameFiles = append(src.gameFiles, newGameFiles...)

	// Add all files
	if len(newGameFiles) > 0 {
		src.collection.AddNewGames(newGameFiles)
	}

	return nil
}

func (src *directorySource) nspCheck(file repository.FileDesc) (bool, error) {
	key, err := src.collection.GetKey(file.GameID)
	if err != nil {
		if src.config.DebugTicket() && err.Error() == "TitleDBKey for game "+file.GameID+" is not found" {
			log.Println(err)
		}
		return false, err
	}

	f, err := os.Open(file.Path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	valid, err := nsp.IsTicketValid(f, key, src.config.DebugTicket())
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("The ticket in '" + file.Path + "' is not valid!")
	}

	return valid, err
}
