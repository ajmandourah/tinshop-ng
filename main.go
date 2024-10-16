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
	"strconv"
	"time"

	"github.com/ajmandourah/tinshop-ng/api"
	"github.com/ajmandourah/tinshop-ng/config"
	collection "github.com/ajmandourah/tinshop-ng/gamescollection"
	"github.com/ajmandourah/tinshop-ng/keys"
	"github.com/ajmandourah/tinshop-ng/repository"
	"github.com/ajmandourah/tinshop-ng/sources"
	"github.com/ajmandourah/tinshop-ng/stats"
	"github.com/ajmandourah/tinshop-ng/utils"
	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
)

//go:embed assets/*
var assetData embed.FS //nolint:gochecknoglobals
// TinShop holds all information about the Shop
type TinShop struct {
	Shop   repository.Shop
	Server *http.Server
}

var creds []string
func main() {

	// this is dirty. will leave it for now untill implemented correctly as there are some conflicts around the shop init

	config := config.New()
	config.LoadConfig()

	prodkeys, _ := keys.InitSwitchKeys(config.ProdKeys())
	if prodkeys == nil || prodkeys.GetKey("header_key") == "" {
		log.Println("!!NOTE!!: keys file was not found, deep scan is disabled, library will be based on file tags.")
		keys.UseKey = false
	}
	collection.Rename = config.Rename()

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
	creds = shop.Shop.Config.Get_Httpauth()

	authOpts := httpauth.AuthOptions{
		Realm: "Tinfoil",
		AuthFunc: HttpAuthCheck,
	}

	r := mux.NewRouter()

	authRoute := r.Methods(http.MethodGet).Subrouter()
	authRoute.HandleFunc("/", shop.HomeHandler)
	authRoute.HandleFunc("/{filter}", shop.FilteringHandler)
	authRoute.HandleFunc("/{filter}/", shop.FilteringHandler)
	authRoute.HandleFunc("/api/{endpoint}", shop.APIHandler)

	r.HandleFunc("/games/{game}", shop.GamesHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)
	

	if len(shop.Shop.Config.Get_Httpauth()) != 0 {
		authRoute.Use(httpauth.BasicAuth(authOpts))
	}

	// r.Use(shop.StatsMiddleware)
	r.Use(shop.TinfoilMiddleware)
	r.Use(shop.CORSMiddleware)
	http.Handle("/", r)

	var port = 3000
	if shop.Shop.Config.Port() != 0 {
		port = shop.Shop.Config.Port()
	}

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:" + strconv.Itoa(port),

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
	myShop.Stats = stats.New()
	myShop.API = api.New()

	// Load collection
	myShop.Collection.Load()

	// Loading config
	myShop.Config.AddHook(myShop.Collection.OnConfigUpdate)
	myShop.Config.AddHook(myShop.Sources.OnConfigUpdate)
	myShop.Config.AddBeforeHook(myShop.Sources.BeforeConfigUpdate)
	myShop.Config.LoadConfig()

	// Loading stats
	myShop.Stats.Load()

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
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(jsonResponse)
	if err != nil {
		log.Println(err)
	}
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

// APIHandler handles api calls
func (s *TinShop) APIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["endpoint"] == "stats" {
		summary, err := s.Shop.Stats.Summary()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		s.Shop.API.Stats(w, summary)
		return
	}
	// Everything not existing
	w.WriteHeader(http.StatusBadRequest)
}

// StatsMiddleware is a middleware to collect statistics
func (s *TinShop) StatsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/" || utils.IsValidFilter(cleanPath(r.RequestURI)) {
			console := &repository.Switch{
				IP:       utils.GetIPFromRequest(r),
				UID:      r.Header.Get("Uid"),
				Theme:    r.Header.Get("Theme"),
				Version:  r.Header.Get("Version"),
				Language: r.Header.Get("Language"),
			}
			_ = s.Shop.Stats.ListVisit(console)
		} else if r.RequestURI[0:7] == "/games/" {
			vars := mux.Vars(r)
			if s.Shop.Sources.HasGame(vars["game"]) {
				_ = s.Shop.Stats.DownloadAsked(utils.GetIPFromRequest(r), vars["game"])
			}
		}
		next.ServeHTTP(w, r)
	})
}

