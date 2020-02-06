package main

import (
	"encoding/json"
	"net/http"
)

func HandleGetTracks(w http.ResponseWriter, r *http.Request) {
	tracks := []Track{}

	store.Find(&tracks, nil)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tracks)
}
