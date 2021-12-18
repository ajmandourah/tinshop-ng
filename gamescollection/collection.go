package gamescollection

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/dblk/tinshop/config"
	"github.com/dblk/tinshop/repository"
	"github.com/dblk/tinshop/utils"
)

var library map[string]map[string]interface{}
var games repository.GameType

// Load ensure that necessary data is loaded
func Load() {
	loadTitlesLibrary()

	initGamesCollection()
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

func initGamesCollection() {
	// Build games object
	games.Success = "Welcome to your own shop!"
	games.Titledb = make(map[string]map[string]interface{})
	games.Files = make([]interface{}, 0)
}

// Reset the collection of files
func Reset(src repository.Sources) {
	initGamesCollection()
}

// Library returns the titledb library
func Library() map[string]map[string]interface{} {
	return library
}

// HasGameIDInLibrary tells if we have gameID information in library
func HasGameIDInLibrary(gameID string) bool {
	return library[gameID] != nil
}

// IsBaseGame tells if the gameID is a base game or not
func IsBaseGame(gameID string) bool {
	return library[gameID]["iconUrl"] != nil
}

// Games returns the games inside the library
func Games() repository.GameType {
	return games
}

// CountGames return the number of games in collection
func CountGames() int {
	var uniqueGames int
	for _, entry := range games.Titledb {
		if entry["iconUrl"] != nil {
			uniqueGames++
		}
	}
	return uniqueGames
}

// AddNewGames increase the games available in the shop
func AddNewGames(newGames []repository.FileDesc) {
	log.Println("Add new games...")
	var gameList = make([]interface{}, 0)

	for _, file := range newGames {
		game := make(map[string]interface{})
		game["url"] = config.GetConfig().RootShop() + "/games/" + file.GameID + "#" + file.GameInfo
		game["size"] = file.Size
		gameList = append(gameList, game)

		if HasGameIDInLibrary(file.GameID) {
			// Verify already present and not update nor dlc
			if games.Titledb[file.GameID] != nil && IsBaseGame(file.GameID) {
				log.Println("Already added id!", file.GameID, file.Path)
			} else {
				games.Titledb[file.GameID] = library[file.GameID]
			}
		} else {
			log.Println("Game not found in database!", file.GameInfo, file.Path)
		}
	}
	games.Files = append(games.Files, gameList...)
	log.Printf("Added %d games in your library\n", len(gameList))
}
