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

	CustomDB() map[string]CustomDBEntry
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
	Success        string                 `json:"success"`
	Titledb        map[string]interface{} `json:"titledb"`
	Files          []GameFileType         `json:"files"`
	ThemeBlackList []string               `json:"themeBlackList,omitempty"`
}

type GameFileType struct {
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

// CustomDBEntry describe the various fields for entries
type CustomDBEntry struct {
	ID          string `mapstructure:"id" json:"id"`
	Name        string `mapstructure:"name" json:"name"`
	Region      string `mapstructure:"region" json:"region"`
	Size        int    `mapstructure:"size" json:"size"`
	ReleaseDate int    `mapstructure:"releaseDate" json:"releaseDate"`
	Description string `mapstructure:"description" json:"description"`
	IconURL     string `mapstructure:"iconUrl" json:"iconUrl"`
}

type WatcherDirectory struct {
	Watcher *fsnotify.Watcher
}
