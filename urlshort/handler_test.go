package urlshort

import "testing"

func TestGetYAML(t *testing.T) {
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	paths, err := getYAML([]byte(yaml))
	correctPaths := []yamlPath{
		{"/urlshort", "https://github.com/gophercises/urlshort"},
		{"/urlshort-final", "https://github.com/gophercises/urlshort/tree/solution"},
	}

	if err != nil {
		t.Error("Should be able to parse valid YAML:", yaml)
	}

	for k, v := range correctPaths {
		if path := paths[k]; path != v {
			t.Error("Expected", v, "Got", path)
		}
	}

	invalidYAML := `
	- path: /urlshort
	  url: https://github.com/gophercises/urlshort
	- path: /urlshort-final
	  url: https://github.com/gophercises/urlshort/tree/solution
`
	paths, err = getYAML([]byte(invalidYAML))

	if err == nil {
		t.Error("Should throw error for invalid YAML")
	}
}

func TestBuildPathMap(t *testing.T) {
	yamlPaths := []yamlPath{
		{"/urlshort", "https://github.com/gophercises/urlshort"},
		{"/urlshort-final", "https://github.com/gophercises/urlshort/tree/solution"},
	}
	correct := map[string]string{
		"/urlshort":       "https://github.com/gophercises/urlshort",
		"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
	}

	pathsToUrls := buildPathMap(yamlPaths)

	for k, v := range correct {
		if path := pathsToUrls[k]; path != v {
			t.Error("Expected", v, "for", k, "Got", path)
		}
	}
}
