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
	"github.com/DblK/tinshop/sources"
	"github.com/DblK/tinshop/utils"
	"github.com/gorilla/mux"
)

//go:embed assets/*
var assetData embed.FS

func main() {
	initServer()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/games/{game}", GamesHandler)
	r.HandleFunc("/{filter}", FilteringHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)
	r.Use(tinfoilMiddleware)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:3000",

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: 0, // Installing large game can take a lot of time
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Printf("Total of %d files in your library (%d in titledb section)\n", len(collection.Games().Files), len(collection.Games().Titledb))
	var uniqueGames = collection.CountGames()
	log.Printf("Total of %d unique games in your library\n", uniqueGames)
	log.Printf("Tinshop available at %s !\n", config.GetConfig().RootShop())

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
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0) //nolint:gocritic
}

func initServer() {
	// Load collection
	collection.Load()

	// Loading config
	config.AddHook(collection.OnConfigUpdate)
	config.AddHook(sources.OnConfigUpdate)
	config.AddBeforeHook(sources.BeforeConfigUpdate)
	config.LoadConfig()
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
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	serveCollection(w, collection.Games())
}

// GamesHandler handles downloading games
func GamesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("Requesting game", vars["game"])

	sources.DownloadGame(vars["game"], w, r)
}

// FilteringHandler handles filtering games collection
func FilteringHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if !utils.IsValidFilter(vars["filter"]) {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	serveCollection(w, collection.Filter(vars["filter"]))
}
