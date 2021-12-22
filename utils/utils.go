// @title tinshop Utils

// @BasePath /utils/

// Package utils provides some cross used information
package utils

import (
	"regexp"
	"strings"

	"github.com/DblK/tinshop/gameid"
	"github.com/DblK/tinshop/repository"
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

func RemoveFileDesc(s []repository.FileDesc, index int) []repository.FileDesc {
	if len(s) < index+1 || index < 0 {
		return s
	}
	ret := make([]repository.FileDesc, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func RemoveGameFile(s []repository.GameFileType, index int) []repository.GameFileType {
	if len(s) < index+1 || index < 0 {
		return s
	}
	ret := make([]repository.GameFileType, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}
