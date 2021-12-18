package utils

import (
	"regexp"
	"strings"

	"github.com/dblk/tinshop/gameid"
	"github.com/dblk/tinshop/repository"
)

// ExtractGameID from fileName the id of game and version
func ExtractGameID(fileName string) repository.GameID {
	ext := strings.Split(fileName, ".")
	re := regexp.MustCompile(`\[(\w{16})\].*\[(v\d+)\]`)
	matches := re.FindStringSubmatch(fileName)

	if len(matches) != 3 {
		return gameid.New("", "", "")
	}

	return gameid.New(strings.ToUpper(matches[1]), "["+strings.ToUpper(matches[1])+"]["+matches[2]+"]."+ext[len(ext)-1], ext[len(ext)-1])
}

// Search returns the index in an object
func Search(length int, f func(index int) bool) int {
	for index := 0; index < length; index++ {
		if f(index) {
			return index
		}
	}
	return -1
}
