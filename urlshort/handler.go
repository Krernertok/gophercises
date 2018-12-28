package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectURL, exists := pathsToUrls[r.URL.Path]

		if exists {
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := getYAML(yml)

	if err != nil {
		return nil, err
	}

	pathsToUrls := buildPathMap(paths)

	return MapHandler(pathsToUrls, fallback), nil
}

func getYAML(yml []byte) ([]yamlPath, error) {
	var paths []yamlPath
	err := yaml.Unmarshal(yml, &paths)

	if err != nil {
		return nil, err
	}

	return paths, nil
}

func buildPathMap(paths []yamlPath) map[string]string {
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
//  [{
//		"path": "/some-path",
//       "url": "https://www.some-url.com/demo"
//	}]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(json string, fallback http.Handler) (http.HandlerFunc, error) {
	return nil, nil
}
