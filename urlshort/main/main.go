package main

import (
	"flag"
	"fmt"
	"gophercises/urlshort/db"
	"gophercises/urlshort/handler"
	"net/http"
	"os"
	"path/filepath"
)

const dbName = "db/test.db"

func getHandler(ext string, data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	switch ext {
	case ".json":
		return handler.JSONHandler(data, fallback)
	case ".yml", ".yaml":
		return handler.YAMLHandler(data, fallback)
	case ".db":
		return handler.DBHandler(dbName, []byte(db.DefaultBucket), fallback)
	}
	return nil, fmt.Errorf("no handler defined for extension: %s", ext)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func main() {
	db.AddData(dbName)
	file := flag.String("config", dbName, "file containing path to URL mapping")
	flag.Parse()

	data, err := os.ReadFile(*file)
	if err != nil {
		fmt.Println("Couldn't open file:", *file)
		fmt.Println(err)
		os.Exit(1)
	}

	extension := filepath.Ext(*file)
	mux := defaultMux()

	handler, err := getHandler(extension, []byte(data), mux)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}
