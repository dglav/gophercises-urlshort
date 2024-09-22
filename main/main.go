package main

import (
	"flag"
	"fmt"
	"net/http"

	urlshort "github.com/dglav/gophercises-urlshort"
)

func main() {
	filepath := flag.String("yaml", "", "Specify the filepath to the YAML file that contains redirects. The filepath starts from the root directory.")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	handler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	if *filepath != "" {
		yamlHandler, err := urlshort.YAMLHandler(*filepath, handler)
		if err != nil {
			fmt.Printf("There was an error in the YAML handler: %v\n", err)
			panic(err)
		}

		handler = yamlHandler
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	fmt.Println("serve default mux")
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
