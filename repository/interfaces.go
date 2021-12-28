// @title tinshop Utils

// @BasePath /repository/

// Package repository holds all interfaces and shared struct
package repository

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// GameID interface
type GameID interface {
	FullID() string
	ShortID() string
	Extension() string
}

// Sources describe all sources type handled
type Sources struct {
	Directories []string `mapstructure:"directories"`
	Nfs         []string `mapstructure:"nfs"`
}

// Config interface
type Config interface {
	RootShop() string
	SetRootShop(string)
	Host() string
	Protocol() string
	Port() int

	DebugNfs() bool
	DebugNoSecurity() bool
	DebugTicket() bool

	Sources() Sources
	Directories() []string
	NfsShares() []string
	ShopTitle() string
	ShopTemplateData() ShopTemplate
	SetShopTemplateData(ShopTemplate)

	IsBlacklisted(string) bool
	IsWhitelisted(string) bool
	IsBannedTheme(string) bool
	BannedTheme() []string

	CustomDB() map[string]TitleDBEntry
	VerifyNSP() bool
}

// ShopTemplate contains all variables used for shop template
type ShopTemplate struct {
	ShopTitle string
}

// HostType new typed string
type HostType string

const (
	// LocalFile Describe local directory file
	LocalFile HostType = "localFile"
	// NFSShare Describe nfs directory file
	NFSShare HostType = "NFS"
)

// FileDesc structure
type FileDesc struct {
	GameID   string
	Size     int64
	GameInfo string
	Path     string
	HostType HostType
}

// GameType structure
type GameType struct {
	Success        string                  `json:"success"`
	Titledb        map[string]TitleDBEntry `json:"titledb"`
	Files          []GameFileType          `json:"files"`
	ThemeBlackList []string                `json:"themeBlackList,omitempty"`
}

// GameFileType stores the fields needed for game files
type GameFileType struct {
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

// NInt is a nullable int
type NInt int

// NString is a nullable string
type NString string

// UnmarshalJSON handles unmarshalling null string
func (n *NString) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*string)(n))
}

// UnmarshalJSON handles unmarshalling null int
func (n *NInt) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	// Handle bad data in json file
	if strings.Contains(string(b), "\"") {
		newNumber, _ := strconv.Atoi(string(b)[1 : len(string(b))-1])
		bs := []byte(strconv.Itoa(newNumber))
		return json.Unmarshal(bs, (*int)(n))
	}
	return json.Unmarshal(b, (*int)(n))
}

// TitleDBEntry describe the various fields for entries
type TitleDBEntry struct {
	ID              string   `mapstructure:"id" json:"id"`
	RightsID        string   `mapstructure:"rightsId" json:"rightsId,omitempty"`
	Name            string   `mapstructure:"name" json:"name,omitempty"`
	Version         NInt     `mapstructure:"version" json:"version,omitempty"`
	Key             NString  `mapstructure:"key" json:"key,omitempty"`
	IsDemo          bool     `mapstructure:"isDemo" json:"isDemo,omitempty"`
	Region          string   `mapstructure:"region" json:"region,omitempty"`
	Regions         []string `mapstructure:"regions" json:"regions,omitempty"`
	ReleaseDate     int      `mapstructure:"releaseDate" json:"releaseDate,omitempty"`
	NsuID           uint64   `mapstructure:"nsuId" json:"nsuId,omitempty"`
	Category        []string `mapstructure:"category" json:"category,omitempty"`
	RatingContent   []string `mapstructure:"ratingContent" json:"ratingContent,omitempty"`
	NumberOfPlayers int      `mapstructure:"numberOfPlayers" json:"numberOfPlayers,omitempty"`
	Publisher       string   `mapstructure:"publisher" json:"publisher,omitempty"`
	Rating          int      `mapstructure:"rating" json:"rating,omitempty"`
	IconURL         string   `mapstructure:"iconUrl" json:"iconUrl,omitempty"`
	BannerURL       string   `mapstructure:"bannerUrl" json:"bannerUrl,omitempty"`
	Intro           string   `mapstructure:"intro" json:"intro,omitempty"`
	Description     string   `mapstructure:"description" json:"description,omitempty"`
	Languages       []string `mapstructure:"languages" json:"languages,omitempty"`
	Size            int      `mapstructure:"size" json:"size,omitempty"`
	Rank            int      `mapstructure:"rank" json:"rank,omitempty"`
}

// Source describes the common functions for Sources
type Source interface {
	Load([]string, bool)
	Download(http.ResponseWriter, *http.Request, string, string)
	UnWatchAll()
	Reset()
	GetFiles() []FileDesc
}
