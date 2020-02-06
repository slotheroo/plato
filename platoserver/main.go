package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/husobee/vestigo"
	"github.com/timshannon/bolthold"
)

const (
	StaticDir    = "static"
	MusicsDir    = "/musics"
	DatabaseName = "store"
)

var (
	store         *bolthold.Store
	executableDir string
)

func init() {

}

func main() {
	var (
		err error
	)
	start := time.Now()

	//Get current executable directory
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	executableDir = filepath.Dir(executablePath)
	err = os.Chdir(executableDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Set working directory to:", executableDir)

	//Get app settings from env or default
	appSettings := CreateAppSettings()

	tracks, err := ScanAll()
	//_, err := ScanAll()
	if err != nil {
		fmt.Println(err)
	}

	store, err = bolthold.Open(DatabaseName, 0600, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer store.Close()

	for _, track := range tracks {
		err = store.Upsert(track.ID, track)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Startup and load time:", time.Since(start))

	//Serve
	fmt.Println("Starting to serve the application on port " + appSettings.Port + ".")
	log.Fatal(http.ListenAndServe(":"+appSettings.Port, Handlers()))
}

//Handlers registers handlers for routes
func Handlers() *vestigo.Router {
	router := vestigo.NewRouter()
	router.Get("/api/tracks", HandleGetTracks)
	staticFileHandler := http.FileServer(http.Dir(StaticDir))
	//router.Handle("/", staticFileHandler)
	router.Get("/*", staticFileHandler.ServeHTTP)
	return router
}

func Handlers2() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tracks", HandleGetTracks)
	staticFiles := http.FileServer(http.Dir(StaticDir))
	mux.Handle("/", staticFiles)
	return mux
}
