package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bogem/id3v2"
	"github.com/tcolgate/mp3"
)

func ScanAll() ([]Track, error) {
	var tracks []Track
	err := filepath.Walk(MusicsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Scanner encountered an error accessing the path %q: %v\n", path, err)
			return err
		}
		if filepath.Ext(path) == ".mp3" {
			// Open file and parse tag in it.
			tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
			if err != nil {
				log.Fatal("Error while opening mp3 file: ", path, err)
			}
			tag.Close()
			subPath := strings.TrimPrefix(path, MusicsDir)
			subPathBytes := []byte(subPath)
			track := Track{
				Artist:   tag.Artist(),
				Album:    tag.Album(),
				FileName: info.Name(),
				ID:       base64.StdEncoding.EncodeToString(subPathBytes),
				Location: subPath,
				Title:    tag.Title(),
				Year:     tag.Year(),
			}
			duration, err := CalculateDuration(path)
			if err != nil {
				return err
			}
			track.Duration = duration
			tracks = append(tracks, track)
		}
		fmt.Printf("visited file or dir: %q\n", path)
		return nil
	})
	if err != nil {
		fmt.Printf("Scanner encountered an error walking the path %q: %v\n", MusicsDir, err)
		return nil, err
	}
	return tracks, nil
}

func CalculateDuration(path string) (int, error) {
	var (
		skipped       = 0
		frame         mp3.Frame
		totalDuration time.Duration
	)
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	mp3Decoder := mp3.NewDecoder(file)
	for {
		if err := mp3Decoder.Decode(&frame, &skipped); err != nil {
			if err == io.EOF {
				break
			} else {
				return 0, err
			}
		}
		totalDuration += frame.Duration()
	}
	return int(totalDuration.Seconds()), nil
}
