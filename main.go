package main

import (
	"context"
	"embed"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

var library map[string]interface{}
var Games map[string]interface{}
var gameFiles []FileDesc
var rootShop string

//go:embed assets/*
var assetData embed.FS

type HostType string

const (
	LocalFile HostType = "localFile"
	NFSShare  HostType = "NFS"
)

type FileDesc struct {
	url      string
	size     int64
	gameInfo string
	path     string
	hostType HostType
}

type GameId struct {
	FullId    string
	ShortId   string
	Extension string
}

func main() {
	initServer()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/games/{game}", GamesHandler)
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
	log.Printf("Total of %d files in your library (%d in titledb section)\n", len(Games["files"].([]interface{})), len(Games["titledb"].(map[string]interface{})))
	var uniqueGames int
	for _, entry := range Games["titledb"].(map[string]interface{}) {
		if entry.(map[string]interface{})["iconUrl"] != nil {
			uniqueGames += 1
		}
	}
	log.Printf("Total of %d unique games in your library\n", uniqueGames)
	log.Printf("Tinshop available at %s !\n", rootShop)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func initServer() {
	// Loading config
	loadConfig()

	// Load JSON library
	loadTitlesLibrary()

	// Load Games
	gameFiles = make([]FileDesc, 0)
	initGamesCollection()
	loadGamesDirectories(len(nfsShares) == 0)
	loadGamesNfsShares()
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

func loadTitlesLibrary() {
	// Open our jsonFile
	jsonFile, err := os.Open("titles.US.en.json")
	if err != nil {

		if err.Error() == "open titles.US.en.json: no such file or directory" {
			log.Println("Missing 'titles.US.en.json'! Start downloading it.")
			downloadErr := DownloadFile("https://github.com/AdamK2003/titledb/releases/download/latest/titles.US.en.json", "titles.US.en.json")
			if downloadErr != nil {
				log.Fatalln(err, downloadErr)
			} else {
				jsonFile, err = os.Open("titles.US.en.json")
				if err != nil {
					log.Fatalln("Error while parsing downloaded json file.\nPlease remove the file and start again the program.\n", err)
				}
			}
		} else {
			log.Fatalln(err)
		}

	}
	log.Println("Successfully Opened titles library")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal([]byte(byteValue), &library)
	if err != nil {
		log.Println("Error while loading titles library", err)
	} else {
		log.Println("Successfully Loaded titles library")
	}
}

func initGamesCollection() {
	// Build games object
	Games = make(map[string]interface{})
	Games["success"] = "Welcome to your own shop!"
	Games["titledb"] = make(map[string]interface{})
	Games["files"] = make([]interface{}, 0)
}

// Handle list of games
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse, jsonError := json.Marshal(Games)

	if jsonError != nil {
		log.Println("Unable to encode JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

// Handle downloading games
func GamesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("Requesting game", vars["game"])

	idx := Search(len(gameFiles), func(index int) bool {
		return gameFiles[index].url == vars["game"]
	})

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Game '%s' not found!", vars["game"])
		return
	} else {
		log.Println(gameFiles[idx].path)
		switch gameFiles[idx].hostType {
		case LocalFile:
			downloadLocalFile(w, r, vars["game"], gameFiles[idx].path)
		case NFSShare:
			downloadNfsFile(w, r, gameFiles[idx].path)

		default:
			w.WriteHeader(http.StatusNotImplemented)
			log.Printf("The type '%s' is not implemented to download game", gameFiles[idx].hostType)
		}
	}
}
