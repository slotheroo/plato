package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/bogem/id3v2"
	"github.com/tcolgate/mp3"
)

func test(path string) {
	fmt.Println("TESTING " + path)
	start := time.Now()

	const (
		Xing = "Xing"
		Info = "Info"
	)
	var (
		skipped = 0
		frame   mp3.Frame
	)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Println("After file open:", time.Since(start))
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		fmt.Println(err)
		return
	}
	tag.Close()
	tagSizeBytes := tag.Size()
	fmt.Println("Seeking bytes:", tagSizeBytes)
	_, err = file.Seek(int64(tagSizeBytes), 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	mp3Decoder := mp3.NewDecoder(file)
	err = mp3Decoder.Decode(&frame, &skipped)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Skipped:", skipped)
	fmt.Println("After decode mp3 frame:", time.Since(start))
	headerByteNum := 4
	channelModeBytes := 32
	/*fmt.Print("FRAME HEADER: ")
	fmt.Println(frame.Header())
	fmt.Print("FRAME: ")
	fmt.Println(frame)
	fmt.Print("CRC: ")
	fmt.Println(frame.CRC())
	fmt.Print("Channel Mode: ")
	fmt.Println(frame.Header().ChannelMode())*/
	if frame.Header().ChannelMode() == mp3.SingleChannel {
		channelModeBytes = 17
	}
	reader := frame.Reader()
	bytes := make([]byte, 64)
	_, err = reader.Read(bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("After read bytes:", time.Since(start))
	/*	fmt.Print("NUM READ: ")
		fmt.Println(numRead)*/
	/*fmt.Print("BYTES: ")
	fmt.Println(bytes)*/
	xingBytes := bytes[headerByteNum+channelModeBytes : headerByteNum+channelModeBytes+4]
	xingString := string(xingBytes)
	fmt.Println("Xing String: " + xingString)
	//	fmt.Println(xingString == Xing || xingString == Info)
	if xingString == Xing || xingString == Info {
		xingFlagByte := bytes[headerByteNum+channelModeBytes+7]
		//fmt.Print("Xing Flag Byte: ")
		//fmt.Println(xingFlagByte)
		hasNumFrames := xingFlagByte&1 == 1
		//fmt.Print("Has Number of Frames: ")
		//fmt.Println(hasNumFrames)
		if hasNumFrames {
			numFramesBytes := bytes[headerByteNum+channelModeBytes+8 : headerByteNum+channelModeBytes+12]
			/*fmt.Print("Number Frames Bytes: ")
			fmt.Println(numFramesBytes)*/
			numFrames := int(binary.BigEndian.Uint32(numFramesBytes))
			//fmt.Print("Number of Frames: ")
			//fmt.Println(numFrames)
			length := int(math.Round(float64(numFrames*1152) / float64(frame.Header().SampleRate())))
			m := length / 60
			s := length % 60
			fmt.Println(fmt.Sprintf("VBR calculated length: %dm %ds", m, s))
		}
	} else {
		stat, err := file.Stat()
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println("Tag Size: ", tag.Size())
		length := int(math.Round(float64((int(stat.Size())-tagSizeBytes)*8) / float64(frame.Header().BitRate())))
		m := length / 60
		s := length % 60
		fmt.Println(fmt.Sprintf("CBR calculated length: %dm %ds", m, s))
	}
	fmt.Println("End:", time.Since(start), "\n")
}
