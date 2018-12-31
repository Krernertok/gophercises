package main

import (
	"flag"
	"fmt"
	"github.com/krernertok/gophercises/urlshort"
	"net/http"
	"os"
)

func main() {
	defaultPathValue := ""
	yamlPath := flag.String("yaml", defaultPathValue, "specifies a path for loading a yaml file")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	handler := urlshort.MapHandler(pathsToUrls, mux)

	// TODO: refactor this into its own function
	if *yamlPath != defaultPathValue {
		yamlFile, fileErr := os.Open(*yamlPath)

		if fileErr != nil {
			panic(fileErr)
		}

		yamlHandler, parseErr := urlshort.YAMLHandler(yamlFile, handler)

		if parseErr != nil {
			panic(parseErr)
		}

		handler = yamlHandler
	}

	// TODO: convert this to be similar to the YAML handler
	jsonData := `
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
	jsonHandler, err := urlshort.JSONHandler(jsonData, handler)

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
