package urlshort

import (
	"os"
	"testing"
)

func TestGetYAML(t *testing.T) {
	filepath := "paths.yaml"
	yamlFile, fileErr := os.Open(filepath)

	if fileErr != nil {
		t.Fatal("Could not open file:")
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

	invalidPath := "invalid.yaml"
	invalidYaml, invalidFileErr := os.Open(invalidPath)
	if invalidFileErr != nil {
		t.Fatal("Could not open file with invalid YAML: ", invalidPath)
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

func TestParseJSON(t *testing.T) {
	validJSON := `
[
	{
		"path": "/json",
		"url": "https://godoc.org/encoding/json"
	},
	{
		"path": "/flag",
		"url": "https://godoc.org/flag"
	}
]`
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

	paths, err := parseJSON([]byte(validJSON))

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
