package urlshort

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		redirectURL := pathsToUrls[path]

		if redirectURL != "" {
			http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
//
//	type PathURL struct {
//		path string `yaml:"path"`
//		url  string `yaml:"url"`
//	}
type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func YAMLHandler(filepath string, fallback http.Handler) (http.HandlerFunc, error) {
	yml, err := readYAML(filepath)
	if err != nil {
		return nil, err
	}

	pathUrls, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	pathToUrls := convertYamlToMap(pathUrls)

	return MapHandler(pathToUrls, fallback), nil
}

func readYAML(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return buf, err
}

func parseYAML(yml []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(yml, &pathURLs)
	return pathURLs, err
}

func convertYamlToMap(pathUrls []pathURL) map[string]string {
	pathToUrls := map[string]string{}

	for _, pathURL := range pathUrls {
		pathToUrls[pathURL.Path] = pathURL.URL
	}

	return pathToUrls
}

type jsonRedirects struct {
	Redirects []jsonPathURL `json:"redirects"`
}

type jsonPathURL struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func JSONHandler(filepath string, fallback http.Handler) (http.HandlerFunc, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	buf, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var pathURLs jsonRedirects
	err = json.Unmarshal(buf, &pathURLs)

	pathToUrls := map[string]string{}
	for _, pathURL := range pathURLs.Redirects {
		pathToUrls[pathURL.Path] = pathURL.URL
	}

	return MapHandler(pathToUrls, fallback), err
}
