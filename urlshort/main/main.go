package main

import (
	"flag"
	"fmt"
	"gophercises/urlshort"
	"net/http"
	"os"
	"path/filepath"
)

func getHandler(ext string, data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	switch ext {
	case ".json":
		return urlshort.JSONHandler(data, fallback)
	case ".yml", ".yaml":
		return urlshort.YAMLHandler(data, fallback)
	}
	return nil, fmt.Errorf("No handler defined for extension: %s", ext)
}

func main() {
	file := flag.String("config", "main/config.json", "file containing path to URL mapping")
	flag.Parse()

	data, err := os.ReadFile(*file)
	if err != nil {
		fmt.Println("Couldn't open file:", *file)
		fmt.Println(err)
		os.Exit(1)
	}

	extension := filepath.Ext(*file)
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	handler, err := getHandler(extension, []byte(data), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
	// http.ListenAndServe(":8080", mapHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
