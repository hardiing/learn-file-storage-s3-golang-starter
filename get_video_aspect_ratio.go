package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"math"
	"os/exec"
)

type VideoInfo struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	var info VideoInfo
	err = json.Unmarshal(out.Bytes(), &info)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	if len(info.Streams) == 0 {
		return "", errors.New("No streams found")
	}

	aspectRatio := GetAspectRatio(info.Streams[0].Width, info.Streams[0].Height)

	switch aspectRatio {
	case "16:9":
		return aspectRatio, err
	case "9:16":
		return aspectRatio, err
	case "other":
		return "other", err
	default:
		return "", err
	}
}

func almostEqual(actual, expected, tolerance float64) bool {
	return math.Abs(actual-expected) <= tolerance
}

func GetAspectRatio(width, height int) string {
	if width <= 0 || height <= 0 {
		return "0:0"
	}
	ratio := float64(width) / float64(height)
	tolerance := 0.02
	if almostEqual(ratio, 16.0/9.0, tolerance) {
		return "16:9"
	} else if almostEqual(ratio, 9.0/16.0, tolerance) {
		return "9:16"
	} else {
		return "other"
	}
}
