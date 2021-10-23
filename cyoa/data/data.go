package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Choice `json:"options"`
}

type Choice struct {
	Arc  string `json:"arc"`
	Text string `json:"text"`
}

func ChaptersFromJson(filePath string) map[string]Chapter {
	var chapters map[string]Chapter

	jsonFile, err := os.Open(filePath)
	if err != nil {
		handleError("Could not open file:", filePath)
	}

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		handleError("Could not read file:", filePath)
	}

	err = json.Unmarshal(jsonBytes, &chapters)
	if err != nil {
		handleError("Could not decode JSON in:", filePath)
	}

	return chapters
}

func handleError(description ...interface{}) {
	fmt.Println(description...)
	os.Exit(1)
}
