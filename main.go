// main.go
package main

import (
	"context"
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/DblK/tinshop/config"
	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/repository"
	"github.com/DblK/tinshop/sources"
	"github.com/DblK/tinshop/utils"
	"github.com/gorilla/mux"
)

//go:embed assets/*
var assetData embed.FS //nolint:gochecknoglobals

// TinShop holds all information about the Shop
type TinShop struct {
	Shop   repository.Shop
	Server *http.Server
}

func main() {
	shop := createShop()

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := shop.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Printf("Total of %d files in your library (%d in titledb section)\n", len(shop.Shop.Collection.Games().Files), len(shop.Shop.Collection.Games().Titledb))
	var uniqueGames = shop.Shop.Collection.CountGames()
	log.Printf("Total of %d unique games in your library\n", uniqueGames)
	log.Printf("Tinshop available at %s !\n", shop.Shop.Config.RootShop())

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15) //nolint:gomnd
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = shop.Server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0) //nolint:gocritic
}

func createShop() TinShop {
	var shop = &TinShop{}

	shop.Shop = initShop()

	r := mux.NewRouter()
	r.HandleFunc("/", shop.HomeHandler)
	r.HandleFunc("/games/{game}", shop.GamesHandler)
	r.HandleFunc("/{filter}", shop.FilteringHandler)
	r.HandleFunc("/{filter}/", shop.FilteringHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)
	r.Use(shop.TinfoilMiddleware)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:3000",

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: 0, // Installing large game can take a lot of time
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	shop.Server = srv

	return *shop
}

// ResetTinshop reset the storage for all information
// func ResetTinshop(myShop repository.Shop) {
// 	shopData = myShop
// }

func initShop() repository.Shop {
	// Init shop data
	myShop := repository.Shop{}
	myShop.Config = config.New()
	myShop.Collection = collection.New(myShop.Config)
	myShop.Sources = sources.New(myShop.Collection)
	// ResetTinshop(myShop)

	// Load collection
	myShop.Collection.Load()

	// Loading config
	myShop.Config.AddHook(myShop.Collection.OnConfigUpdate)
	myShop.Config.AddHook(myShop.Sources.OnConfigUpdate)
	myShop.Config.AddBeforeHook(myShop.Sources.BeforeConfigUpdate)
	myShop.Config.LoadConfig()

	return myShop
}

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Println("notFound")
	log.Println(r.Header)
	log.Println(r.RequestURI)
	w.WriteHeader(http.StatusNotFound)
}

func notAllowed(w http.ResponseWriter, r *http.Request) {
	log.Println("notAllowed")
	log.Println(r.Header)
	log.Println(r.RequestURI)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func serveCollection(w http.ResponseWriter, tinfoilCollection interface{}) {
	jsonResponse, jsonError := json.Marshal(tinfoilCollection)

	if jsonError != nil {
		log.Println("Unable to encode JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

// HomeHandler handles list of games
func (s *TinShop) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if s.Shop.Collection == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	serveCollection(w, s.Shop.Collection.Games())
}

// GamesHandler handles downloading games
func (s *TinShop) GamesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("Requesting game", vars["game"])

	s.Shop.Sources.DownloadGame(vars["game"], w, r)
}

// FilteringHandler handles filtering games collection
func (s *TinShop) FilteringHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if !utils.IsValidFilter(vars["filter"]) {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if s.Shop.Collection == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	serveCollection(w, s.Shop.Collection.Filter(vars["filter"]))
}
