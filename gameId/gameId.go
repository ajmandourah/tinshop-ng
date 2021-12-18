package gameid

import (
	"github.com/dblk/tinshop/repository"
)

type gameID struct {
	fullID    string
	shortID   string
	extension string
}

// New create a new GameID
func New(shortID, fullID, extension string) repository.GameID {
	return &gameID{
		fullID:    fullID,
		shortID:   shortID,
		extension: extension,
	}
}

func (game *gameID) FullID() string {
	return game.fullID
}
func (game *gameID) ShortID() string {
	return game.shortID
}
func (game *gameID) Extension() string {
	return game.extension
}
