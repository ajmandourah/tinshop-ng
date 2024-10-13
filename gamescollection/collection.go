// @title tinshop Games Collection

// @BasePath /gamescollection/

// Package gamescollection provides and stores all information about collection
package gamescollection

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ajmandourah/tinshop-ng/repository"
	"github.com/ajmandourah/tinshop-ng/utils"
)

var Rename bool = false

type collect struct {
	games         repository.GameType
	library       map[string]repository.TitleDBEntry
	mergedLibrary map[string]repository.TitleDBEntry
	config        repository.Config
}

// New create a new collection
func New(config repository.Config) repository.Collection {
	return &collect{
		config: config,
	}
}

// Load ensure that necessary data is loaded
func (c *collect) Load() {
	c.loadTitlesLibrary()

	c.ResetGamesCollection()
}

func (c *collect) loadTitlesLibrary() {
	//determine path 
	var jsonPath string
	if _, err := os.Stat("/data/config.yaml"); !os.IsNotExist(err) {
		jsonPath = "/data/titles.US.en.json" 
	}else{
		jsonPath = "titles.US.en.json"
	}

	// Open our jsonFile
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		log.Println("Missing 'titles.US.en.json'! Start downloading it.")
		downloadErr := utils.DownloadFile("https://tinfoil.media/repo/db/titles.json", jsonPath )
		if downloadErr != nil {
			log.Fatalln(err, downloadErr)
	}
	}
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		log.Fatalln("Error while parsing downloaded json file.\nPlease remove the file and start again the program.\n", err)
	}

	log.Println("Successfully Opened titles library")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &c.library)
	if err != nil {
		log.Println("Error while loading titles library", err)
	} else {
		log.Println("Successfully Loaded titles library")
	}
}

// ResetGamesCollection reset the game collection
func (c *collect) ResetGamesCollection() {
	// Build games object
	if c.config.NoWelcomeMessage() {
		c.games.Success = ""
	} else {
		c.games.Success = c.config.WelcomeMessage()
	}
	c.games.Titledb = make(map[string]repository.TitleDBEntry)
	c.games.Files = make([]repository.GameFileType, 0)
	c.games.ThemeBlackList = nil
}

// OnConfigUpdate the collection of files
func (c *collect) OnConfigUpdate(cfg repository.Config) {
	c.config = cfg
	c.ResetGamesCollection()

	// Create merged library
	c.mergedLibrary = make(map[string]repository.TitleDBEntry)

	// Copy library
	for key, entry := range c.library {
		gameID := strings.ToUpper(key)

		c.mergedLibrary[gameID] = entry
	}

	// Copy CustomDB
	for key, entry := range c.config.CustomDB() {
		gameID := strings.ToUpper(key)
		if _, ok := c.mergedLibrary[gameID]; ok {
			log.Println("Duplicate customDB entry from official titledb (consider removing from configuration)", gameID)
		} else {
			c.mergedLibrary[gameID] = entry
		}
	}

	// Check if blacklist entries
	if len(c.config.BannedTheme()) != 0 {
		c.games.ThemeBlackList = c.config.BannedTheme()
	} else {
		c.games.ThemeBlackList = nil
	}
}

// Library returns the titledb library
func (c *collect) Library() map[string]repository.TitleDBEntry {
	return c.mergedLibrary
}

// HasGameIDInLibrary tells if we have gameID information in library
func (c *collect) HasGameIDInLibrary(gameID string) bool {
	_, ok := c.Library()[gameID]
	return ok
}

// IsBaseGame tells if the gameID is a base game or not
func (c *collect) IsBaseGame(gameID string) bool {
	return c.Library()[gameID].IconURL != ""
}

// Games returns the games inside the library
func (c *collect) Games() repository.GameType {
	return c.games
}

// Filter returns the games inside the library after filtering
func (c *collect) Filter(filter string) repository.GameType {
	var filteredGames repository.GameType
	if !c.config.NoWelcomeMessage() {
		filteredGames.Success = c.games.Success
	}
	filteredGames.ThemeBlackList = c.games.ThemeBlackList
	upperFilter := strings.ToUpper(filter)

	newTitleDB := make(map[string]repository.TitleDBEntry)
	newFiles := make([]repository.GameFileType, 0)
	for ID, entry := range c.games.Titledb {
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
			idx := utils.Search(len(c.games.Files), func(index int) bool {
				return strings.Contains(c.games.Files[index].URL, ID)
			})

			if idx != -1 {
				newFiles = append(newFiles, c.games.Files[idx])
			}
		}
	}
	filteredGames.Titledb = newTitleDB
	filteredGames.Files = newFiles

	return filteredGames
}

// RemoveGame remove ID from the collection
func (c *collect) RemoveGame(ID string) {
	gameID := strings.ToUpper(ID)
	log.Println("Removing game", gameID)

	// Remove from Files entry
	idx := utils.Search(len(c.games.Files), func(index int) bool {
		return strings.Contains(c.games.Files[index].URL, gameID)
	})

	if idx != -1 {
		c.games.Files = utils.RemoveGameFile(c.games.Files, idx)
	}

	// Remove from titledb entry
	delete(c.games.Titledb, gameID)
}

// CountGames return the number of games in collection
func (c *collect) CountGames() int {
	var uniqueGames int
	for _, entry := range c.games.Titledb {
		if entry.IconURL != "" || entry.BannerURL != "" {
			uniqueGames++
		}
	}
	return uniqueGames
}

// AddNewGames increase the games available in the shop
func (c *collect) AddNewGames(newGames []repository.FileDesc) {
	log.Println("Add new games...")
	var gameList = make([]repository.GameFileType, 0)

	for _, file := range newGames {
		game := repository.GameFileType{
			URL:  c.config.RootShop() + "/games/" + file.GameID + "#" + c.getFriendlyName(file),
			Size: file.Size,
		}

		// Handle duplicate entry for file
		idx := utils.Search(len(gameList), func(index int) bool {
			return strings.Contains(gameList[index].URL, file.GameID)
		})
		if idx == -1 {
			gameList = append(gameList, game)
		} else {
			log.Println("Duplicate Game", file.GameID, file.Path)
		}

		if c.HasGameIDInLibrary(file.GameID) {
			// Verify already present and not update nor dlc
			if _, ok := c.games.Titledb[file.GameID]; ok && c.IsBaseGame(file.GameID) {
				log.Println("Already added id!", file.GameID, file.Path)
			} else {
				c.games.Titledb[file.GameID] = c.Library()[file.GameID]
			}
		} else {
			log.Println("Game not found in database!", file.GameInfo, file.Path)
		}
	}
	c.games.Files = append(c.games.Files, gameList...)
	log.Printf("Added %d games in your library\n", len(gameList))
}

// GetKey return the key from the titledb
func (c *collect) GetKey(gameID string) (string, error) {
	var key = c.Library()[gameID].Key
	if key == "" {
		return "", errors.New("TitleDBKey for game " + gameID + " is not found")
	}
	return string(key), nil
}

func (c *collect) getFriendlyName(file repository.FileDesc) string {
	baseID, update, dlc := utils.GetTitleMeta(file.GameID)
	baseTitle := c.Library()[baseID]
	title := c.Library()[file.GameID]

	// Default extra for Base title
	var extra = " [BASE]"

	// Append DLC Name and tag when dlc
	if dlc {
		extra = " - " + title.Name + " [DLC]"
	}

	// Append version when update
	if update {
		extra = fmt.Sprintf(" [v%d][UPD]", title.Version)
	}

	name := ""
	if baseTitle.Name != "" {
		name = " " + baseTitle.Name
	}

	region := ""
	if baseTitle.Region != "" {
		region = " (" + baseTitle.Region + ")"
	}

	// Build the friendly name for Tinfoil
	reg := []string{"[" + file.GameID + "]", name, region, extra, "." + file.Extension}
	return strings.Join(reg[:], "")
}

func (c *collect) GenTitle(gameID string) string {
	baseID, update, dlc := utils.GetTitleMeta(gameID)
	baseTitle := c.Library()[baseID]
	title := c.Library()[gameID]

	// Default extra for Base title
	var extra = " [BASE]"

	// Append DLC Name and tag when dlc
	if dlc {
		extra = " - " + title.Name + " [DLC]"
	}

	// Append version when update
	if update {
		extra = fmt.Sprintf("[UPD]")
	}
	version := fmt.Sprintf("[v%d]", title.Version)
	name := ""
	if baseTitle.Name != "" {
		name = baseTitle.Name
	} else {
		name = title.Name
	}

	region := ""
	if baseTitle.Region != "" {
		region = " (" + baseTitle.Region + ")"
	}

	// Build the friendly name for Tinfoil
	reg := []string{name, region, extra, "[" + gameID + "]", version}
	return strings.Join(reg[:], "")
}
