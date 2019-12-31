package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/timshannon/bolthold"
)

const (
	StaticDir    = "static/"
	MusicsDir    = "static/musics/"
	DatabaseName = "store"
)

func init() {

}

func main() {
	//Get app settings from env or default
	appSet := CreateAppSettings()

	tracks, err := ScanAll()
	if err != nil {
		fmt.Println(err)
	}

	store, err := bolthold.Open(DatabaseName, 0600, nil)
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

	returnData := []Track{}

	store.Find(&returnData, bolthold.Where("Artist").Eq("The Abominable Slowmen"))
	fmt.Println(returnData[1])

	//Serve
	fmt.Println("Starting to serve the application on port " + appSet.Port + ".")
	log.Fatal(http.ListenAndServe(":"+appSet.Port, Handlers()))

	//
}

//Handlers registers handlers for routes
func Handlers() *http.ServeMux {
	mux := http.NewServeMux()
	staticFiles := http.FileServer(http.Dir(StaticDir))
	mux.Handle("/", staticFiles)
	return mux
}
