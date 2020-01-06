package main

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/bogem/id3v2"
	"github.com/tcolgate/mp3"
)

func ScanAll() ([]Track, error) {
	var (
		tracks []Track
	)
	err := filepath.Walk(StaticDir+MusicsDir, func(path string, info os.FileInfo, err error) error {
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
			location := strings.TrimPrefix(path, StaticDir)
			subPath := strings.TrimPrefix(location, MusicsDir)
			subPathBytes := []byte(subPath)
			track := Track{
				Artist:   tag.Artist(),
				Album:    tag.Album(),
				FileName: info.Name(),
				ID:       base64.StdEncoding.EncodeToString(subPathBytes),
				Location: location,
				Title:    tag.Title(),
				Year:     tag.Year(),
			}
			duration, err := CalculateDuration(path, tag.Size())
			if err != nil {
				fmt.Println(err)
			} else {
				track.Duration = duration
			}
			tracks = append(tracks, track)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Scanner encountered an error walking the path %q: %v\n", StaticDir+MusicsDir, err)
		return nil, err
	}
	fmt.Println("Scanned", len(tracks), "tracks")
	return tracks, nil
}

func CalculateDuration(path string, skipBytes int) (int, error) {
	const (
		Xing = "Xing"
		Info = "Info"
	)
	var (
		skipped          = 0
		frame            mp3.Frame //We'll grab the first frame of the mp3 only
		headerByteNum    = 4       //Header byte length for an mp3
		channelModeBytes = 32      //Default if channel mode is not mono
	)
	//Open the file
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	//Skip past the size of the id3 tag at the beginning of the mp3 before searching
	//for the frame header. Significantly reduces scan time when an image is present
	//in the id3 tags
	_, err = file.Seek(int64(skipBytes), 0)
	if err != nil {
		return 0, err
	}
	//Decode frame
	mp3Decoder := mp3.NewDecoder(file)
	err = mp3Decoder.Decode(&frame, &skipped)
	if err != nil {
		return 0, err
	}
	//If mp3 is mono, change value of channelModeBytes
	if frame.Header().ChannelMode() == mp3.SingleChannel {
		channelModeBytes = 17
	}
	//The mp3 package doesn't scan for the Xing/Info header, so we have to do it
	reader := frame.Reader()
	bytes := make([]byte, 64)
	_, err = reader.Read(bytes)
	if err != nil {
		return 0, err
	}
	//If Xing/Info string is present it will be at this location relative to start
	//of header
	xingString := string(bytes[headerByteNum+channelModeBytes : headerByteNum+channelModeBytes+4])
	if xingString == Xing || xingString == Info {
		//flag byte will tell us what info is present in the following bytes
		xingFlagByte := bytes[headerByteNum+channelModeBytes+7]
		//if flag byte has its 0 bit set then "number of frames" info is present
		hasNumFrames := xingFlagByte&1 == 1
		if hasNumFrames {
			//Get number of frames and then calculate mp3 duration
			numFramesBytes := bytes[headerByteNum+channelModeBytes+8 : headerByteNum+channelModeBytes+12]
			numFrames := int(binary.BigEndian.Uint32(numFramesBytes))
			//Duration = Number of Frames * Samples Per Frame / Sampling Rate
			duration := int(math.Round(float64(numFrames*1152) / float64(frame.Header().SampleRate())))
			return duration, nil
		} else {
			return 0, errors.New("Xing/Info header present but number of frames count is missing.")
		}
	}
	//No Xing/Info bytes, assuming file is constant bit rate
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	//Duration = ((File Size as bytes - bytes in id3 tag) * 8) / Bitrate
	duration := int(math.Round(float64((int(stat.Size())-skipBytes)*8) / float64(frame.Header().BitRate())))
	return duration, nil
}
