package urlshort

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Feedback from solution:
		// path := r.URL.Path
		// url, ok := pathsToUrls[path]; ok {/* Redirect here */}
		redirectURL, exists := pathsToUrls[r.URL.Path]

		if exists {
			// Feedback from solution: Should use 302 Status Found (http.StatusFound)
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
			// Feedback from solution: Should return
		}

		fallback.ServeHTTP(w, r)
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml *os.File, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := getYAML(yml)

	if err != nil {
		return nil, err
	}

	pathsToUrls := pathMapFromYAML(paths)

	return MapHandler(pathsToUrls, fallback), nil
}

func getYAML(yml *os.File) ([]yamlPath, error) {
	var paths []yamlPath
	err := yaml.NewDecoder(yml).Decode(&paths)

	if err != nil {
		return nil, err
	}

	return paths, nil
}

func pathMapFromYAML(paths []yamlPath) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, p := range paths {
		pathsToUrls[p.Path] = p.URL
	}
	return pathsToUrls
}

type yamlPath struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// JSONHandler will parse the provided JSON string and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//  [
//		{
//			"path": "/some-path",
//       	"url": "https://www.some-url.com/demo"
//		}
//	]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(jsonFile *os.File, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := getJSON(jsonFile)

	if err != nil {
		return nil, err
	}

	pathsToUrls := pathMapFromJSON(paths)

	return MapHandler(pathsToUrls, fallback), nil
}

func getJSON(jsonFile *os.File) ([]jsonPath, error) {
	var paths []jsonPath
	err := json.NewDecoder(jsonFile).Decode(&paths)

	if err != nil {
		return nil, err
	}

	return paths, nil
}

type jsonPath struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func pathMapFromJSON(paths []jsonPath) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, p := range paths {
		pathsToUrls[p.Path] = p.URL
	}
	return pathsToUrls
}

// DBHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths to their corresponding URL. DBHandler will use the
// database provided as an argument to map paths to URLs.
// If the path is not provided in the database, then the fallback
// http.Handler will be called instead.
func DBHandler(db *sql.DB, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var url string

		query := "SELECT url FROM urlshort WHERE path=$1;"
		path := r.URL.Path
		result := db.QueryRow(query, path)

		switch err := result.Scan(&url); err {
		case nil:
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		default:
			fmt.Println("Error accessing database:", err)
			fallback.ServeHTTP(w, r)
		}
	})
}
