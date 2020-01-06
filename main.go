package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/timshannon/bolthold"
)

const (
	StaticDir    = "static/"
	MusicsDir    = "musics/"
	DatabaseName = "store"
)

var (
	store *bolthold.Store
)

func init() {

}

func main() {
	start := time.Now()

	//Get app settings from env or default
	appSet := CreateAppSettings()

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
	fmt.Println("Starting to serve the application on port " + appSet.Port + ".")
	log.Fatal(http.ListenAndServe(":"+appSet.Port, Handlers()))
}

//Handlers registers handlers for routes
func Handlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tracks", HandleGetTracks)
	staticFiles := http.FileServer(http.Dir(StaticDir))
	mux.Handle("/", staticFiles)
	return mux
}
