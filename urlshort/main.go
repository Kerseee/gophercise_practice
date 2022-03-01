package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"kerseeehuang.com/urlshort/shorteners"
)

type config struct {
	yamlFile string
}

func main() {
	// Configure the settings.
	var cfg config
	flag.StringVar(&cfg.yamlFile, "yaml", "yamls/example.yaml", "yaml file path")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := shorteners.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := os.ReadFile(cfg.yamlFile)
	if err != nil {
		panic(err)
	}
	yamlHandler, err := shorteners.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
