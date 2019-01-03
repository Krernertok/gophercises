package story

import (
	"encoding/json"
	"os"
)

// Option contains Text and Arc fields
type Option struct {
	Text    string `json:"text"`
	ArcName string `json:"arc"`
}

// Arc depicts a single story object
type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

// ParseStory parses the JSON file at the given path and returns
// a map[string]Arc.
func ParseStory(filePath string) (map[string]Arc, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	story := make(map[string]Arc)
	err = json.NewDecoder(file).Decode(&story)
	if err != nil {
		return nil, err
	}

	return story, nil
}
