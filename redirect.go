package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type server struct {
	mappings map[string]map[string]string
	lock     sync.RWMutex
}

func (s *server) matchURL(domain, path string) string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	paths, ok := s.mappings[strings.ToLower(domain)]
	if !ok {
		return ""
	}
	return paths[path]
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domain := r.Host
	domain = strings.Split(domain, ":")[0]
	fmt.Println(domain, r.URL.Path)
	to := s.matchURL(domain, r.URL.Path)
	if to == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}
	http.Redirect(w, r, to, http.StatusMovedPermanently)
}

func (s *server) setMappings(m map[string]map[string]string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.mappings = m
}
