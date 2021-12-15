package main

import (
	"log"
	"strings"
)

func AddNewGames(newGames []FileDesc) {
	log.Printf("\n\nAdd new games...\n")
	var gameList = make([]interface{}, 0)
	for _, file := range newGames {
		game := make(map[string]interface{})
		game["url"] = rootShop + "/games/" + file.url + "#" + file.gameInfo
		game["size"] = file.size
		gameList = append(gameList, game)

		gameID := strings.ToUpper(file.url)
		if library[gameID] != nil {
			// Verify already present and not update nor dlc
			if Games["titledb"].(map[string]interface{})[gameID] != nil && library[gameID].(map[string]interface{})["iconUrl"] != nil {
				log.Println("Already added id!", gameID, file.path)
			} else {
				Games["titledb"].(map[string]interface{})[gameID] = library[gameID]
			}
		} else {
			log.Println("Game not found in database!", file.gameInfo, file.path)
		}
	}
	Games["files"] = append(Games["files"].([]interface{}), gameList...)
	log.Printf("Added %d games in your library\n", len(gameList))
}
