package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func loadConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("Config not found!")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		// TODO: Reload config on change
	})
	viper.WatchConfig()

	// ----------------------------------------------------------
	// General config
	// ----------------------------------------------------------
	host := viper.Get("host")
	protocol := viper.Get("protocol")
	port := viper.Get("port")

	if protocol == nil {
		rootShop = "http"
	} else {
		rootShop = protocol.(string)
	}
	rootShop = rootShop + "://"
	if host == nil {
		// Retrieve current IP
		host, _ := os.Hostname()
		addrs, _ := net.LookupIP(host)
		var myIp = ""
		for _, addr := range addrs {
			if ipv4 := addr.To4(); ipv4 != nil {
				if myIp == "" {
					myIp = ipv4.String()
				}
			}
		}
		rootShop = rootShop + myIp
	} else {
		rootShop = rootShop + host.(string)
	}
	if port == nil {
		rootShop = rootShop + ":3000"
	} else {
		rootShop = rootShop + ":" + strconv.Itoa(port.(int))
	}

	// ----------------------------------------------------------
	// Debug
	// ----------------------------------------------------------
	debugNfs = viper.GetBool("debug.nfs")

	// ----------------------------------------------------------
	// Sources
	// ----------------------------------------------------------
	// Directories
	cfgDirectories := viper.GetStringSlice("sources.directories")
	if cfgDirectories == nil {
		// Default search
		directories = make([]string, 0)
		directories = append(directories, "./games")
	} else {
		directories = cfgDirectories
	}

	// NFS
	cfgNfs := viper.GetStringSlice("sources.nfs")
	if cfgNfs != nil {
		nfsShares = cfgNfs
	}

	// ----------------------------------------------------------
	// Shop Template
	// ----------------------------------------------------------
	shopTitle := viper.GetString("name")
	if shopTitle == "" {
		shopTitle = "TinShop"
	}
	shopTemplateData = ShopTemplate{
		ShopTitle: shopTitle,
	}
}
