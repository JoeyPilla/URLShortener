package URLShortener

import (
	"net/http"

	yamlV2 "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusPermanentRedirect)
		}
		fallback.ServeHTTP(w, r)
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

func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	yamlMap, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathsToURLMap := buildMap(yamlMap)
	return MapHandler(pathsToURLMap, fallback), err
}

func parseYAML(yaml []byte) (yamlMap []map[string]string, err error) {
	err = yamlV2.Unmarshal(yaml, &yamlMap)
	return yamlMap, err
}

func buildMap(yamlMap []map[string]string) map[string]string {
	returnMap := make(map[string]string)
	for _, value := range yamlMap {
		key := value["path"]
		returnMap[key] = value["url"]
	}
	return returnMap
}
