package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"

	"github.com/dblk/tinshop/repository"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type debug struct {
	Nfs        bool
	NoSecurity bool
}

type config struct {
	rootShop         string
	ShopHost         string             `mapstructure:"host"`
	ShopProtocol     string             `mapstructure:"protocol"`
	ShopPort         int                `mapstructure:"port"`
	Debug            debug              `mapstructure:"debug"`
	AllSources       repository.Sources `mapstructure:"sources"`
	Name             string             `mapstructure:"name"`
	shopTemplateData repository.ShopTemplate
}

var serverConfig config
var allHooks []func(repository.Sources)

// LoadConfig handles viper under the hood
func LoadConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.SetDefault("sources.directories", "./games")
	viper.SetDefault("sources.nfs", "") // FIXME: Hack for issue viper

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("Config not found!")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed, update new configuration...")
		configChange()
	})
	viper.WatchConfig()

	configChange()
}

func configChange() {
	serverConfig = loadAndCompute()
	// FIXME: Hack if nfs array become empty
	if reflect.TypeOf(viper.AllSettings()["sources"].(map[string]interface{})["nfs"]).String() == "string" {
		serverConfig.AllSources.Nfs = make([]string, 0)
	}

	// Call all hooks
	for _, hook := range allHooks {
		hook(serverConfig.AllSources)
	}
}

// GetConfig returns the current configuration
func GetConfig() repository.Config {
	return &serverConfig
}

func loadAndCompute() config {
	err := viper.Unmarshal(&serverConfig)

	if err != nil {
		log.Fatalln(err)
	}
	ComputeDefaultValues(&serverConfig)

	return serverConfig
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
		var myIP = ""
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
	if config.Port() == 0 {
		rootShop += ":3000"
	} else if !(config.Port() == 443 && config.Protocol() == "https") && !(config.Port() == 80 && config.Protocol() == "http") {
		rootShop += ":" + strconv.Itoa(config.Port())
	}
	config.SetRootShop(rootShop)

	config.SetShopTemplateData(repository.ShopTemplate{
		ShopTitle: config.ShopTitle(),
	})

	return config
}

// HookOnSource Add hook function on change sources
func HookOnSource(f func(repository.Sources)) {
	allHooks = append(allHooks, f)
}

func (cfg *config) SetRootShop(root string) {
	cfg.rootShop = root
}

func (cfg *config) RootShop() string {
	return cfg.rootShop
}
func (cfg *config) Protocol() string {
	return cfg.ShopProtocol
}
func (cfg *config) Host() string {
	return cfg.ShopHost
}
func (cfg *config) Port() int {
	return cfg.ShopPort
}
func (cfg *config) DebugNfs() bool {
	return cfg.Debug.Nfs
}
func (cfg *config) DebugNoSecurity() bool {
	return cfg.Debug.NoSecurity
}
func (cfg *config) Directories() []string {
	return cfg.AllSources.Directories
}
func (cfg *config) NfsShares() []string {
	return cfg.AllSources.Nfs
}
func (cfg *config) Sources() repository.Sources {
	return cfg.AllSources
}
func (cfg *config) ShopTemplateData() repository.ShopTemplate {
	return cfg.shopTemplateData
}
func (cfg *config) SetShopTemplateData(data repository.ShopTemplate) {
	cfg.shopTemplateData = data
}
func (cfg *config) ShopTitle() string {
	return cfg.Name
}
