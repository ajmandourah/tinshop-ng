// @title tinshop Config

// @BasePath /config/

// Package config provides everything related to configuration
package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/ajmandourah/tinshop-ng/repository"
	"github.com/ajmandourah/tinshop-ng/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type debug struct {
	Nfs        bool
	NoSecurity bool
	Ticket     bool
}

type security struct {
	Whitelist   []string `mapstructure:"whitelist"`
	Blacklist   []string `mapstructure:"blacklist"`
	BannedTheme []string `mapstructure:"bannedTheme"`
	ForwardAuth string   `mapstructure:"forwardAuth"`
	Hauth string         `mapstructure:"hauth"`	
}

type nsp struct {
	CheckVerified bool `mapstructure:"checkVerified"`
}

// Configuration holds all config information
type Configuration struct {
	rootShop             string
	ShopHost             string                             `mapstructure:"host"`
	ShopProtocol         string                             `mapstructure:"protocol"`
	Keys		     string				`mapstructure:"keys"`
	RenameFiles          bool                               `mapstructure:"renameFiles"`
	ShopWelcomeMessage   string                             `mapstructure:"welcomeMessage"`
	ShopNoWelcomeMessage bool                               `mapstructure:"noWelcomeMessage"`
	ShopPort             int                                `mapstructure:"port"`
	Debug                debug                              `mapstructure:"debug"`
	Proxy                bool                               `mapstructure:"reverseProxy"`
	AllSources           repository.ConfigSources           `mapstructure:"sources"`
	Name                 string                             `mapstructure:"name"`
	Security             security                           `mapstructure:"security"`
	CustomTitleDB        map[string]repository.TitleDBEntry `mapstructure:"customTitledb"`
	NSP                  nsp                                `mapstructure:"nsp"`
	shopTemplateData     repository.ShopTemplate

	allHooks       []func(repository.Config)
	beforeAllHooks []func(repository.Config)
}

// New returns a new configuration
func New() repository.Config {
	return &Configuration{}
}

// LoadConfig handles viper under the hood
func (cfg *Configuration) LoadConfig() {
	viper.SetConfigName("config")     // name of config file (without extension)
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/data")      // for the docker image
	viper.AddConfigPath(".")          // optionally look for config in the working directory
	viper.SetTypeByDefaultValue(true) // Allows []string to be parsed from Env Vars

	viper.SetDefault("host", "")
	viper.SetDefault("protocol", "http")
	viper.SetDefault("keys","prod.keys")
	viper.SetDefault("renameFiles","false")
	viper.SetDefault("name", "TinShop")
	viper.SetDefault("reverseProxy", false)
	viper.SetDefault("welcomeMessage", "Welcome to your own TinShop!")
	viper.SetDefault("noWelcomeMessage", false)

	viper.SetDefault("debug.nfs", false)
	viper.SetDefault("debug.noSecurity", false)
	viper.SetDefault("debug.ticket", false)

	viper.SetDefault("nsp.checkVerified", false)

	viper.SetDefault("sources.directories", []string{"./games"})
	viper.SetDefault("sources.nfs", []string{})

	viper.SetDefault("security.bannedTheme", []string{})
	viper.SetDefault("security.whitelist", []string{})
	viper.SetDefault("security.blacklist", []string{})
	viper.SetDefault("security.forwardAuth", "")
	viper.SetDefault("security.hauth", "")
	

	viper.SetEnvPrefix("TINSHOP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("Config not found!")
			err := viper.WriteConfigAs("/data/config.yaml")
			if err != nil {
				log.Println("Error while creating config.yaml in /data, most likely you are not in a container.")
			}
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed, update new configuration...")
		cfg.configChange()
	})
	viper.WatchConfig()

	cfg.configChange()
}

func (cfg *Configuration) configChange() {
	// Call all before hooks
	for _, hook := range cfg.beforeAllHooks {
		hook(cfg)
	}

	newConfig := loadAndCompute()
	cfg.rootShop = newConfig.rootShop
	cfg.ShopHost = newConfig.ShopHost
	cfg.ShopProtocol = newConfig.ShopProtocol
	cfg.Keys = newConfig.Keys
	cfg.RenameFiles = newConfig.RenameFiles
	cfg.ShopPort = newConfig.ShopPort
	if newConfig.ShopWelcomeMessage != "" {
		cfg.ShopWelcomeMessage = newConfig.ShopWelcomeMessage
	} else {
		cfg.ShopWelcomeMessage = "Welcome to your own TinShop!"
	}
	cfg.ShopNoWelcomeMessage = newConfig.ShopNoWelcomeMessage
	cfg.Proxy = newConfig.Proxy
	cfg.Debug = newConfig.Debug
	cfg.AllSources = newConfig.AllSources
	cfg.Name = newConfig.Name
	cfg.Security = newConfig.Security
	cfg.CustomTitleDB = newConfig.CustomTitleDB
	cfg.NSP = newConfig.NSP
	cfg.shopTemplateData = newConfig.shopTemplateData

	// Call all hooks
	for _, hook := range cfg.allHooks {
		hook(cfg)
	}
}

func loadAndCompute() *Configuration {
	var loadedConfig = &Configuration{}
	err := viper.Unmarshal(&loadedConfig)

	if err != nil {
		log.Fatalln(err)
	}
	ComputeDefaultValues(loadedConfig)

	return loadedConfig
}

// ComputeDefaultValues change the value taken from the config file
func ComputeDefaultValues(config repository.Config) repository.Config {
	// ----------------------------------------------------------
	// Compute rootShop url
	// ----------------------------------------------------------
	var rootShop string
	if config.Protocol() == "" {
		rootShop = "http"
	} else {
		rootShop = config.Protocol()
	}
	rootShop += "://"
	if config.Host() == "" {
		// Retrieve current IP
		host, _ := os.Hostname()
		addrs, _ := net.LookupIP(host)
		var myIP = "0.0.0.0"
		for _, addr := range addrs {
			if ipv4 := addr.To4(); ipv4 != nil {
				if myIP == "" {
					myIP = ipv4.String()
				}
			}
		}
		rootShop += myIP
	} else {
		rootShop += config.Host()
	}

	if !config.ReverseProxy() {
		if config.Port() == 0 {
			rootShop += ":3000"
		} else if !(config.Port() == 443 && config.Protocol() == "https") && !(config.Port() == 80 && config.Protocol() == "http") {
			rootShop += ":" + strconv.Itoa(config.Port())
		}
	}
	log.Println(rootShop)
	config.SetRootShop(rootShop)

	config.SetShopTemplateData(repository.ShopTemplate{
		ShopTitle: config.ShopTitle(),
	})

	return config
}

// AddHook Add hook function on change config
func (cfg *Configuration) AddHook(f func(repository.Config)) {
	cfg.allHooks = append(cfg.allHooks, f)
}

// AddBeforeHook Add hook function before on change config
func (cfg *Configuration) AddBeforeHook(f func(repository.Config)) {
	cfg.beforeAllHooks = append(cfg.beforeAllHooks, f)
}

// SetRootShop allow to change the root url of the shop
func (cfg *Configuration) SetRootShop(root string) {
	cfg.rootShop = root
}

// RootShop returns the RootShop url
func (cfg *Configuration) RootShop() string {
	return cfg.rootShop
}

// ReverseProxy returns the ReverseProxy setting
func (cfg *Configuration) ReverseProxy() bool {
	return cfg.Proxy
}

// WelcomeMessage returns the WelcomeMessage
func (cfg *Configuration) WelcomeMessage() string {
	return cfg.ShopWelcomeMessage
}

// NoWelcomeMessage returns the NoWelcomeMessage
func (cfg *Configuration) NoWelcomeMessage() bool {
	return cfg.ShopNoWelcomeMessage
}

// Protocol returns the protocol scheme (http or https)
func (cfg *Configuration) Protocol() string {
	return cfg.ShopProtocol
}

// Protocol returns the path/keys file
func (cfg *Configuration) ProdKeys() string {
	return cfg.Keys
}

// Tells if files should be renamed after decrypting
func (cfg *Configuration) Rename() bool {
	return cfg.RenameFiles
}
// Host returns the host of the shop
func (cfg *Configuration) Host() string {
	return cfg.ShopHost
}

// Port returns the port number for outside access
func (cfg *Configuration) Port() int {
	return cfg.ShopPort
}

// DebugTicket tells if we should display additional log for ticket verification
func (cfg *Configuration) DebugTicket() bool {
	return cfg.Debug.Ticket
}

// DebugNfs tells if we should display additional log for nfs
func (cfg *Configuration) DebugNfs() bool {
	return cfg.Debug.Nfs
}

// DebugNoSecurity returns if we should disable security or not
func (cfg *Configuration) DebugNoSecurity() bool {
	return cfg.Debug.NoSecurity
}

// Directories returns the list of directories sources
func (cfg *Configuration) Directories() []string {
	return cfg.AllSources.Directories
}

// CustomDB returns the list of custom title db
func (cfg *Configuration) CustomDB() map[string]repository.TitleDBEntry {
	return cfg.CustomTitleDB
}

// NfsShares returns the list of nfs sources
func (cfg *Configuration) NfsShares() []string {
	return cfg.AllSources.Nfs
}

// Sources returns all available sources
func (cfg *Configuration) Sources() repository.ConfigSources {
	return cfg.AllSources
}

// ShopTemplateData returns the data needed to render template
func (cfg *Configuration) ShopTemplateData() repository.ShopTemplate {
	return cfg.shopTemplateData
}

// SetShopTemplateData sets the data for template
func (cfg *Configuration) SetShopTemplateData(data repository.ShopTemplate) {
	cfg.shopTemplateData = data
}

// ShopTitle returns the name of the shop
func (cfg *Configuration) ShopTitle() string {
	return cfg.Name
}

// VerifyNSP tells if we need to verify NSP
func (cfg *Configuration) VerifyNSP() bool {
	return cfg.NSP.CheckVerified
}

// ForwardAuthURL returns the url of the forward auth
func (cfg *Configuration) ForwardAuthURL() string {
	return cfg.Security.ForwardAuth
}

//get Hauth code
func (cfg *Configuration) Get_Hauth() string {
	return cfg.Security.Hauth
}


// IsBlacklisted tells if the uid is blacklisted or not
func (cfg *Configuration) IsBlacklisted(uid string) bool {
	if len(cfg.Security.Whitelist) != 0 {
		return !cfg.isInWhiteList(uid)
	}
	return cfg.isInBlackList(uid)
}

// IsWhitelisted tells if the uid is whitelisted or not
func (cfg *Configuration) IsWhitelisted(uid string) bool {
	if len(cfg.Security.Whitelist) == 0 {
		return !cfg.isInBlackList(uid)
	}
	return cfg.isInWhiteList(uid)
}

func (cfg *Configuration) isInBlackList(uid string) bool {
	idxBlackList := utils.Search(len(cfg.Security.Blacklist), func(index int) bool {
		return cfg.Security.Blacklist[index] == uid
	})
	return idxBlackList != -1
}
func (cfg *Configuration) isInWhiteList(uid string) bool {
	idxWhiteList := utils.Search(len(cfg.Security.Whitelist), func(index int) bool {
		return cfg.Security.Whitelist[index] == uid
	})
	return idxWhiteList != -1
}

// IsBannedTheme tells if the theme is banned or not
func (cfg *Configuration) IsBannedTheme(theme string) bool {
	idxBannedTheme := utils.Search(len(cfg.Security.BannedTheme), func(index int) bool {
		return cfg.Security.BannedTheme[index] == theme
	})
	return idxBannedTheme != -1
}

// BannedTheme returns all banned theme
func (cfg *Configuration) BannedTheme() []string {
	return cfg.Security.BannedTheme
}
