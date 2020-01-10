package main

type Track struct {
	Album    string `json:"album"`
	Artist   string `json:"artist"`
	Duration int    `json:"duration"`
	FileName string `json:"fileName"`
	ID       string `json:"id"`
	Location string `json:"location"`
	Title    string `json:"title"`
	Year     string `json:"year"`
}
