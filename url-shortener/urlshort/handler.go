package urlshort

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// pathUrl stores a path and corresponding url.
type pathUrl struct {
	path 	string
	url 	string
}

var errInvalidPath = errors.New("handle: invalid path in yaml")

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		urlPath, ok := pathsToUrls[r.URL.Path]
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, urlPath, http.StatusSeeOther)
	})
}

// parseYAML parse the yaml with package gopkg.in/yaml 
// and return parsed data and error if there is any.
func parseYAML(yml []byte) ([]map[string]string, error) {
	var urls []map[string]string
	err := yaml.Unmarshal(yml, &urls)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// pathsToUrlsMap return a map that uses path as key and url as value.
// pathsUrls should be a slice of map that in each map there are only two elements:
// path:{a path} and url:{a url}.
func pathsToUrlsMap(pathsUrls []map[string]string) (map[string]string, error) {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathsUrls {
		path, ok := pu["path"]
		if !ok || path == "" {
			return nil, errInvalidPath
		}
		pathsToUrls[path] = pu["url"]
	}
	return pathsToUrls, nil
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
	// Parse the yaml.
	urls, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	// Build the map.
	pathsToUrls, err := pathsToUrlsMap(urls)
	if err != nil {
		return nil, err
	}
	fmt.Println(urls, pathsToUrls)
	return MapHandler(pathsToUrls, fallback), nil
}