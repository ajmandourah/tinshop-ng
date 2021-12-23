// @title tinshop Games Collection

// @BasePath /gamescollection/

// Package gamescollection provides and stores all information about collection
package gamescollection

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/DblK/tinshop/config"
	"github.com/DblK/tinshop/repository"
	"github.com/DblK/tinshop/utils"
)

var library map[string]repository.TitleDBEntry
var mergedLibrary map[string]repository.TitleDBEntry
var games repository.GameType

// Load ensure that necessary data is loaded
func Load() {
	loadTitlesLibrary()

	ResetGamesCollection()
}

func loadTitlesLibrary() {
	// Open our jsonFile
	jsonFile, err := os.Open("titles.US.en.json")

	if err != nil {
		if err.Error() == "open titles.US.en.json: no such file or directory" {
			log.Println("Missing 'titles.US.en.json'! Start downloading it.")
			downloadErr := utils.DownloadFile("https://github.com/AdamK2003/titledb/releases/download/latest/titles.US.en.json", "titles.US.en.json")
			if downloadErr != nil {
				log.Fatalln(err, downloadErr)
			} else {
				jsonFile, err = os.Open("titles.US.en.json")
				if err != nil {
					log.Fatalln("Error while parsing downloaded json file.\nPlease remove the file and start again the program.\n", err)
				}
			}
		} else {
			log.Fatalln(err)
		}
	}
	log.Println("Successfully Opened titles library")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &library)
	if err != nil {
		log.Println("Error while loading titles library", err)
	} else {
		log.Println("Successfully Loaded titles library")
	}
}

// ResetGamesCollection reset the game collection
func ResetGamesCollection() {
	// Build games object
	games.Success = "Welcome to your own shop!"
	games.Titledb = make(map[string]repository.TitleDBEntry)
	games.Files = make([]repository.GameFileType, 0)
	games.ThemeBlackList = nil
}

// OnConfigUpdate the collection of files
func OnConfigUpdate(cfg repository.Config) {
	ResetGamesCollection()

	// Create merged library
	mergedLibrary = make(map[string]repository.TitleDBEntry)

	// Copy library
	for key, entry := range library {
		gameID := strings.ToUpper(key)

		mergedLibrary[gameID] = entry
	}

	// Copy CustomDB
	for key, entry := range cfg.CustomDB() {
		gameID := strings.ToUpper(key)
		if _, ok := mergedLibrary[gameID]; ok {
			log.Println("Duplicate customDB entry from official titledb (consider removing from configuration)", gameID)
		} else {
			mergedLibrary[gameID] = entry
		}
	}

	// Check if blacklist entries
	if len(cfg.BannedTheme()) != 0 {
		games.ThemeBlackList = cfg.BannedTheme()
	} else {
		games.ThemeBlackList = nil
	}
}

// Library returns the titledb library
func Library() map[string]repository.TitleDBEntry {
	return mergedLibrary
}

// HasGameIDInLibrary tells if we have gameID information in library
func HasGameIDInLibrary(gameID string) bool {
	_, ok := Library()[gameID]
	return ok
}

// IsBaseGame tells if the gameID is a base game or not
func IsBaseGame(gameID string) bool {
	return Library()[gameID].IconURL != ""
}

// Games returns the games inside the library
func Games() repository.GameType {
	return games
}

// Filter returns the games inside the library after filtering
func Filter(filter string) repository.GameType {
	var filteredGames repository.GameType
	filteredGames.Success = games.Success
	filteredGames.ThemeBlackList = games.ThemeBlackList
	upperFilter := strings.ToUpper(filter)

	newTitleDB := make(map[string]repository.TitleDBEntry)
	newFiles := make([]repository.GameFileType, 0)
	for ID, entry := range games.Titledb {
		entryFiltered := false
		if upperFilter == "WORLD" {
			entryFiltered = true
		} else if upperFilter == "MULTI" {
			numberPlayers := entry.NumberOfPlayers

			if numberPlayers > 1 {
				entryFiltered = true
			}
		} else if utils.IsValidFilter(upperFilter) {
			languages := entry.Languages

			if utils.Contains(languages, upperFilter) || utils.Contains(languages, strings.ToLower(upperFilter)) {
				entryFiltered = true
			}
		}

		if entryFiltered {
			newTitleDB[ID] = entry
			idx := utils.Search(len(games.Files), func(index int) bool {
				return strings.Contains(games.Files[index].URL, ID)
			})

			if idx != -1 {
				newFiles = append(newFiles, games.Files[idx])
			}
		}
	}
	filteredGames.Titledb = newTitleDB
	filteredGames.Files = newFiles

	return filteredGames
}

// RemoveGame remove ID from the collection
func RemoveGame(ID string) {
	gameID := strings.ToUpper(ID)
	log.Println("Removing game", gameID)

	// Remove from Files entry
	idx := utils.Search(len(games.Files), func(index int) bool {
		return strings.Contains(games.Files[index].URL, gameID)
	})

	if idx != -1 {
		games.Files = utils.RemoveGameFile(games.Files, idx)
	}

	// Remove from titledb entry
	delete(games.Titledb, gameID)
}

// CountGames return the number of games in collection
func CountGames() int {
	var uniqueGames int
	for _, entry := range games.Titledb {
		if entry.IconURL != "" {
			uniqueGames++
		}
	}
	return uniqueGames
}

// AddNewGames increase the games available in the shop
func AddNewGames(newGames []repository.FileDesc) {
	log.Println("Add new games...")
	var gameList = make([]repository.GameFileType, 0)

	for _, file := range newGames {
		game := repository.GameFileType{
			URL:  config.GetConfig().RootShop() + "/games/" + file.GameID + "#" + file.GameInfo,
			Size: file.Size,
		}
		gameList = append(gameList, game)

		if HasGameIDInLibrary(file.GameID) {
			// Verify already present and not update nor dlc
			if _, ok := games.Titledb[file.GameID]; ok && IsBaseGame(file.GameID) {
				log.Println("Already added id!", file.GameID, file.Path)
			} else {
				games.Titledb[file.GameID] = Library()[file.GameID]
			}
		} else {
			log.Println("Game not found in database!", file.GameInfo, file.Path)
		}
	}
	games.Files = append(games.Files, gameList...)
	log.Printf("Added %d games in your library\n", len(gameList))
}
