// Package traefik_blockpath a plugin to block a path.
package traefik_blockpath

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"slices"
)

// Config holds the plugin configuration.
type Config struct {
	Methods    []string `yaml:"methods,omitempty"`
	Regex      []string `yaml:"regex,omitempty"`
	StatusCode int      `yaml:"statusCode"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Methods:    make([]string, 0),
		Regex:      make([]string, 0),
		StatusCode: http.StatusForbidden,
	}
}

type blockPath struct {
	methods    []string
	name       string
	next       http.Handler
	regexps    []*regexp.Regexp
	statuscode int
}

// New creates and returns a plugin instance.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	regexps := make([]*regexp.Regexp, len(config.Regex))

	for regexIndex, regex := range config.Regex {
		re, err := regexp.Compile(regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", regex, err)
		}

		regexps[regexIndex] = re
	}

	return &blockPath{
		methods:    config.Methods,
		name:       name,
		next:       next,
		regexps:    regexps,
		statuscode: config.StatusCode,
	}, nil
}

func (b *blockPath) ServeHTTP(responseWriter http.ResponseWriter, req *http.Request) {
	currentPath := req.URL.EscapedPath()
	currentMethod := req.Method

	for _, re := range b.regexps {
		if re.MatchString(currentPath) {
			// if methods not defined, block all methods
			if len(b.methods) == 0 {
				responseWriter.WriteHeader(b.statuscode)
				return
			}
			// else block only for that specific methods
			if slices.Contains(b.methods, currentMethod) {
				responseWriter.WriteHeader(b.statuscode)
				return
			}
		}
	}

	b.next.ServeHTTP(responseWriter, req)
}
