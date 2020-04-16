package urlshort

import (
	"net/http"

    yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(writer http.ResponseWriter, reader *http.Request) {

		path := reader.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(writer, reader, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(writer, reader)
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {

    pathURLs, error := parseYaml(yamlBytes)
    if error != nil {
        return nil, error
    }

    pathsToURLs := buildMap(pathURLs)

	return MapHandler(pathsToURLs, fallback), nil
}

func parseYaml(data []byte) ([]pathURL, error) {

    var pathURLs []pathURL
	error := yaml.Unmarshal(yamlBytes, &pathURLs)
	if error != nil {
		return nil, error
	}

    return pathURLs, nil
}

func buildMap(pathURLs []pathURL) map[string]string {

    pathsToURLs := make(map[string]string)
	for _, pathURL := range pathURLs {
		pathsToUrls[pathURL.Path] = pathURL.URL
	}

    return pathsToURLs
}

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
