package main

import (
	// "database/sql"
	"flag"
	"fmt"
	"github.com/krernertok/gophercises/urlshort"
	// "github.com/krernertok/gophercises/urlshort/config"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	/* dbConfig := config.GetDBConfig()

	dbInfoStr := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	dbInfo := fmt.Sprintf(
		dbInfoStr,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
	)

	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	mux := defaultMux()
	dbHandler := urlshort.DBHandler(db, mux)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler) */

	// ORIGINAL HANDLER
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	handler := urlshort.MapHandler(pathsToUrls, mux)

	// USING FLAGS
	defaultPathValue := ""
	yamlPath := flag.String("yaml", defaultPathValue, "specifies a path for loading a yaml file")
	jsonPath := flag.String("json", defaultPathValue, "specifies a path for loading a JSON file")
	flag.Parse()

	// USING YAML
	if *yamlPath != defaultPathValue {
		yamlHandler, err := getFileHandler(*yamlPath, urlshort.YAMLHandler, handler)

		if err != nil {
			panic(err)
		}

		handler = yamlHandler
	}

	if *jsonPath != defaultPathValue {
		jsonHandler, err := getFileHandler(*jsonPath, urlshort.JSONHandler, handler)

		if err != nil {
			panic(err)
		}

		handler = jsonHandler
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)

}

func getFileHandler(path string, handler func(*os.File, http.Handler) (http.HandlerFunc, error),
	fallback http.Handler) (http.HandlerFunc, error) {
	file, fileErr := os.Open(path)

	if fileErr != nil {
		return nil, fileErr
	}

	fileHandler, parseErr := handler(file, fallback)

	if parseErr != nil {
		return nil, parseErr
	}

	return fileHandler, nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
