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

// ConfigSources describe all sources type handled
type ConfigSources struct {
	Directories []string `mapstructure:"directories"`
	Nfs         []string `mapstructure:"nfs"`
}

// Config interface
type Config interface {
	RootShop() string
	SetRootShop(string)
	Host() string
	Protocol() string
	ProdKeys() string
	Rename() bool
	Port() int
	ReverseProxy() bool
	WelcomeMessage() string
	NoWelcomeMessage() bool

	DebugNfs() bool
	DebugNoSecurity() bool
	DebugTicket() bool

	Sources() ConfigSources
	Directories() []string
	NfsShares() []string
	ShopTitle() string
	ShopTemplateData() ShopTemplate
	SetShopTemplateData(ShopTemplate)

	ForwardAuthURL() string
	Get_Hauth() string
	Get_Httpauth() []string
	IsBlacklisted(string) bool
	IsWhitelisted(string) bool
	IsBannedTheme(string) bool
	BannedTheme() []string

	CustomDB() map[string]TitleDBEntry
	VerifyNSP() bool

	AddHook(f func(Config))
	AddBeforeHook(f func(Config))
	LoadConfig()
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
	GameID    string
	Size      int64
	GameInfo  string
	Path      string
	Extension string
	HostType  HostType
}

// GameType structure
type GameType struct {
	Success        string         `json:"success,omitempty"`
	Files          []GameFileType `json:"files"`
	ThemeBlackList []string       `json:"themeBlackList,omitempty"`
	// Removing the titledb for the resulted json.
	Titledb map[string]TitleDBEntry `json:"-"`
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
	Rating          int      `mapstructure:"rating" json:"rating,omitempty"`
	Developer       string   `mapstructure:"developer" json:"developer,omitempty"`
	Publisher       string   `mapstructure:"publisher" json:"publisher,omitempty"`
	FrontBoxArt     string   `mapstructure:"frontBoxArt" json:"frontBoxArt,omitempty"`
	IconURL         string   `mapstructure:"iconUrl" json:"iconUrl,omitempty"`
	Screenshots     []string `mapstructure:"screenshots" json:"screenshots,omitempty"`
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

// Sources describes all function to handle all sources
type Sources interface {
	OnConfigUpdate(Config)
	BeforeConfigUpdate(Config)
	GetFiles() []FileDesc
	HasGame(string) bool
	DownloadGame(string, http.ResponseWriter, *http.Request)
}

// Collection describes all information about collection
type Collection interface {
	Load()
	OnConfigUpdate(Config)
	Filter(string) GameType
	RemoveGame(string)
	CountGames() int
	AddNewGames([]FileDesc)
	Library() map[string]TitleDBEntry
	HasGameIDInLibrary(string) bool
	IsBaseGame(string) bool
	Games() GameType
	GetKey(string) (string, error)
	ResetGamesCollection()
	GenTitle(string) (string, bool)
}

// Switch holds all information about the switch
type Switch struct {
	IP       string
	UID      string
	Theme    string
	Version  string
	Language string
}

// StatsSummary holds all information about tinshop
type StatsSummary struct {
	Visit           uint64                 `json:"visit,omitempty"`
	UniqueSwitch    uint64                 `json:"uniqueSwitch,omitempty"`
	VisitPerSwitch  map[string]interface{} `json:"visitPerSwitch,omitempty"`
	DownloadAsked   uint64                 `json:"downloadAsked,omitempty"`
	DownloadDetails map[string]interface{} `json:"downloadDetails,omitempty"`
}

// Stats holds all information about statistics
type Stats interface {
	Load()
	Close() error
	ListVisit(*Switch) error
	DownloadAsked(string, string) error
	Summary() (StatsSummary, error)
}

// Shop holds all tinshop information
type Shop struct {
	Collection Collection
	Sources    Sources
	Config     Config
	Stats      Stats
	API        API
}

// API holds all function for api
type API interface {
	Stats(http.ResponseWriter, StatsSummary)
}
