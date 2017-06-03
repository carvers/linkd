package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func parseDomain(filename string) string {
	if filename == "" {
		return filename
	}
	filename = strings.TrimSpace(filename)
	filename = filepath.Base(filename)
	return filename
}

func loadMappings(files []string) (map[string]map[string]string, error) {
	mappings := map[string]map[string]string{}
	var err error
	for _, file := range files {
		mappings[parseDomain(file)], err = loadMapping(file)
		if err != nil {
			return nil, errors.Wrapf(err, "error loading mappings from %s", file)
		}
	}
	return mappings, nil
}

func loadMapping(filename string) (map[string]string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "error opening file")
	}
	mappings, err := parseMapping(string(b))
	if err != nil {
		return nil, errors.Wrap(err, "error parsing file")
	}
	return mappings, nil
}

func parseMapping(contents string) (map[string]string, error) {
	lines := strings.Split(contents, "\n")
	links := map[string]string{}
	for pos, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pieces := strings.Split(line, " -> ")
		if len(pieces) != 2 {
			return nil, fmt.Errorf("Error parsing line %d: invalid format %q", pos, line)
		}
		links[pieces[0]] = pieces[1]
	}
	return links, nil
}
