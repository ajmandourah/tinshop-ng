// @title tinshop Utils

// @BasePath /utils/

// Package utils provides some cross used information
package utils

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/DblK/tinshop/gameid"
	"github.com/DblK/tinshop/repository"
)

// All languages available
// To update this list run: `jq '[.[].regions] | del(..|nulls) | flatten | unique' titles.US.en.json`
var languageFilter = []string{
	"AR", "AT", "AU", "BE", "CA", "CL", "CN", "CO", "CZ", "DE",
	"DK", "ES", "FI", "FR", "GB", "GR", "HK", "HU", "IT", "JP",
	"KR", "MX", "NL", "NO", "NZ", "PE", "PL", "PT", "RU", "SE",
	"US", "ZA",
}

// IsValidFilter returns true if the filter is handled
func IsValidFilter(filter string) bool {
	upperFilter := strings.ToUpper(filter)
	return upperFilter == "MULTI" || upperFilter == "WORLD" || Contains(languageFilter, upperFilter)
}

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

// Contains returns true if an element is present in a slice
func Contains(list interface{}, elem interface{}) bool {
	listV := reflect.ValueOf(list)

	if listV.Kind() == reflect.Slice {
		for i := 0; i < listV.Len(); i++ {
			item := listV.Index(i).Interface()

			target := reflect.ValueOf(elem).Convert(reflect.TypeOf(item)).Interface()
			if ok := reflect.DeepEqual(item, target); ok {
				return true
			}
		}
	}
	return false
}
