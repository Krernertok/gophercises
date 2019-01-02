package urlshort

import (
	"os"
	"testing"
)

func TestGetYAML(t *testing.T) {
	filepath := "data/paths.yaml"
	yamlFile, fileErr := os.Open(filepath)

	if fileErr != nil {
		t.Fatal("Could not open file:", filepath)
	}

	paths, parseErr := getYAML(yamlFile)
	correctPaths := []yamlPath{
		{"/urlshort", "https://github.com/gophercises/urlshort"},
		{"/urlshort-final", "https://github.com/gophercises/urlshort/tree/solution"},
	}

	if parseErr != nil {
		t.Error("Should be able to parse valid YAML:", filepath)
	}

	for k, v := range correctPaths {
		if path := paths[k]; path != v {
			t.Error("Expected", v, "Got", path)
		}
	}

	invalidYAMLPath := "data/invalid.yaml"
	invalidYaml, invalidFileErr := os.Open(invalidYAMLPath)
	if invalidFileErr != nil {
		t.Fatal("Could not open file with invalid YAML: ", invalidYAMLPath)
	}

	paths, parseErr = getYAML(invalidYaml)
	if parseErr == nil {
		t.Error("Should throw error for invalid YAML")
	}
}

func TestPathMapFromYAML(t *testing.T) {
	yamlPaths := []yamlPath{
		{"/urlshort", "https://github.com/gophercises/urlshort"},
		{"/urlshort-final", "https://github.com/gophercises/urlshort/tree/solution"},
	}
	correct := map[string]string{
		"/urlshort":       "https://github.com/gophercises/urlshort",
		"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
	}

	pathsToUrls := pathMapFromYAML(yamlPaths)

	for k, v := range correct {
		if path := pathsToUrls[k]; path != v {
			t.Error("Expected", v, "for", k, "Got", path)
		}
	}
}

func TestGetJSON(t *testing.T) {
	validFilePath := "data/paths.json"
	invalidFilePath := "data/invalid.json"

	expectedPaths := []jsonPath{
		{
			Path: "/json",
			URL:  "https://godoc.org/encoding/json",
		},
		{
			Path: "/flag",
			URL:  "https://godoc.org/flag",
		},
	}

	validFile, err := os.Open(validFilePath)

	if err != nil {
		t.Fatal("Couldn't open valid JSON file.")
	}

	paths, err := getJSON(validFile)

	if err != nil {
		t.Error("Couldn't parse valid JSON: ", err)
	}

	if len(paths) != len(expectedPaths) {
		t.Fatal(
			"Expected", expectedPaths,
			"Got", paths,
		)
	}

	for i, ep := range expectedPaths {
		rp := paths[i]
		if ep.Path != rp.Path {
			t.Error(
				"Expected path", ep.Path,
				"Got path", rp.Path,
			)
		}
		if ep.URL != rp.URL {
			t.Error(
				"Expected URL", ep.URL,
				"Got URL", rp.URL,
			)
		}
	}

	invalidFile, err := os.Open(invalidFilePath)

	if err != nil {
		t.Fatal("Couldn't open invalid JSON file.")
	}

	paths, err = getJSON(invalidFile)

	if err == nil {
		t.Error("Shouldn't be able to parse invalid JSON: ", invalidFilePath)
	}
}

func TestPathMapFromJSON(t *testing.T) {
	jsonPaths := []jsonPath{
		{"/urlshort", "https://github.com/gophercises/urlshort"},
		{"/urlshort-final", "https://github.com/gophercises/urlshort/tree/solution"},
	}
	correct := map[string]string{
		"/urlshort":       "https://github.com/gophercises/urlshort",
		"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
	}

	pathsToUrls := pathMapFromJSON(jsonPaths)

	for k, v := range correct {
		if path := pathsToUrls[k]; path != v {
			t.Error("Expected", v, "for", k, "Got", path)
		}
	}
}
