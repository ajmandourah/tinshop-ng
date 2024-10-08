// @title tinshop Utils

// @BasePath /utils/

// Package utils provides some cross used information
package utils

import (
	"fmt"
	"math/big"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/ajmandourah/tinshop/gameid"
	"github.com/ajmandourah/tinshop/repository"
)

// IsValidFilter returns true if the filter is handled
func IsValidFilter(filter string) bool {
	// All languages available
	// To update this list run: `jq '[.[].regions] | del(..|nulls) | flatten | unique' titles.US.en.json`
	var languageFilter = []string{
		"AR", "AT", "AU", "BE", "CA", "CL", "CN", "CO", "CZ", "DE",
		"DK", "ES", "FI", "FR", "GB", "GR", "HK", "HU", "IT", "JP",
		"KR", "MX", "NL", "NO", "NZ", "PE", "PL", "PT", "RU", "SE",
		"US", "ZA",
	}

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

// GetTitleMeta returns the BaseID of the content, as well as Update / DLC flags
func GetTitleMeta(titleID string) (string, bool, bool) {
	var lastDigit = titleID[len(titleID)-1:]
	var baseID = strings.Join([]string{titleID[:len(titleID)-3], "000"}, "")
	var update = false
	var dlc = false

	if titleID != baseID {
		update = true
	}

	if lastDigit != "0" {
		dlc = true
		update = false

		// Parse the hexadecimal string into a big integer
		intValue, success := new(big.Int).SetString(baseID, 16)
		if !success {
			return "", false, false
		}

		// Parse the subtraction value (in hexadecimal)
		subtractionValue := new(big.Int)
		subtractionValue, success = subtractionValue.SetString("1000", 16)
		if !success {
			return "", false, false
		}

		// Subtract the values
		intValue.Sub(intValue, subtractionValue)

		// Convert the resulting integer back to a hexadecimal string, left padded with 0 to 16 chars
		baseID = fmt.Sprintf("0000000000000000%X", intValue)
		baseID = baseID[len(baseID)-16:]
	}
	return baseID, update, dlc
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

// RemoveFileDesc remove a specific index inside a repository.FileDesc
func RemoveFileDesc(s []repository.FileDesc, index int) []repository.FileDesc {
	if len(s) < index+1 || index < 0 {
		return s
	}
	ret := make([]repository.FileDesc, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

// RemoveGameFile remove a specific index inside a repository.GameFileType
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

// GetIPFromRequest returns ip from the request
func GetIPFromRequest(r *http.Request) string {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if r.Header.Get("X-Forwarded-For") != "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	return ip
}
