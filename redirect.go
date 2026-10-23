package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/carvers/linkd/storers"
)

type server struct {
	datastore storers.Datastore
}

func (s *server) matchURL(ctx context.Context, domain, path string) (string, error) {
	return s.datastore.GetLink(ctx, domain, path)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domain := r.Host
	domain, _, _ = strings.Cut(domain, ":")
	to, err := s.matchURL(r.Context(), domain, r.URL.Path)
	if errors.Is(err, storers.ErrLinkNotFound) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	} else if err != nil {
		log.Println("Error retrieving link:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		return
	}
	http.Redirect(w, r, to, http.StatusMovedPermanently)
}
