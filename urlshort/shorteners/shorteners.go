// shorteners contains some handlers that parse files like yaml or json and return
// http handlers.
package shorteners

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the path in the request
		path := r.URL.Path

		// Get the url in the pathsToUrls map.
		url, ok := pathsToUrls[path]

		// If url not exist, call the fallback handler.
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}

		// Redirect to the url.
		http.Redirect(w, r, url, http.StatusSeeOther)
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
	// Unmarshal the yaml.
	var inputs []struct {
		Path string `yaml:"path"`
		Url  string `yaml:"url"`
	}
	err := yaml.Unmarshal(yml, &inputs)
	if err != nil {
		return nil, err
	}

	// Build the map.
	pathsToUrls := make(map[string]string)
	for _, input := range inputs {
		pathsToUrls[input.Path] = input.Url
	}

	// Return the handler.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the path in the request
		path := r.URL.Path

		// Get the url in the pathsToUrls map.
		url, ok := pathsToUrls[path]

		// If url not exist, call the fallback handler.
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}

		// Redirect to the url.
		http.Redirect(w, r, url, http.StatusSeeOther)

	}), nil
}
