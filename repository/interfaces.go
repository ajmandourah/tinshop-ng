// @title tinshop Utils

// @BasePath /repository/

// Package repository holds all interfaces and shared struct
package repository

import (
	"gopkg.in/fsnotify.v1"
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

type GameFileType struct {
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

// TitleDBEntry describe the various fields for entries
type TitleDBEntry struct {
	ID              string   `mapstructure:"id" json:"id"`
	RightsID        string   `mapstructure:"rightsId" json:"rightsId,omitempty"`
	Name            string   `mapstructure:"name" json:"name,omitempty"`
	Version         uint     `mapstructure:"version" json:"version,omitempty"`
	Key             string   `mapstructure:"key" json:"key,omitempty"`
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

type WatcherDirectory struct {
	Watcher *fsnotify.Watcher
}
