package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"

	"carvers.dev/linkd/storers"
)

const (
	formTmpl = `<html>
  <head>
    <title>{{ .Name }}</title>
  </head>
  <body>
    <h1>Add a new link</h1>
    <form method="post">
      <select name="domain">{{ range .Domains }}
        <option value="{{ .Domain }}">{{ .Domain }}</option>{{ end }}
      </select>
      <label>Path</label>
      <input type="text" name="path" />
      <label>Target</label>
      <input type="text" name="target" />
      <input type="submit">
    </form>
    <h1>Existing links</h1>
    <table>
      <tr><th>Domain</th><th>Path</th><th>Target</th><th>Created At</th><th>Created By</th></tr>{{ range .Links }}
      <tr><td>{{ .Domain }}</td><td>{{ .Path }}</td><td><a href="{{ .Target }}">{{ .Target }}</a></td><td>{{ .CreatedAt }}</td><td>{{ .CreatedBy }}</td></tr>{{ end }}
    </table>
  </body>
</html>`
)

type TemplateData struct {
	Name    string
	Domains []storers.DatastoreDomain
	Links   []storers.DatastoreLink
}

type server struct {
	datastore storers.Datastore

	adminDomain string
	name        string
}

func (s *server) matchURL(ctx context.Context, domain, path string) (string, error) {
	return s.datastore.GetLink(ctx, domain, path)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domain := r.Host
	domain = strings.Split(domain, ":")[0]
	if domain == s.adminDomain {
		s.ServeAdminHTTP(w, r)
		return
	}
	to, err := s.matchURL(r.Context(), domain, r.URL.Path)
	if err == storers.ErrLinkNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	} else if err != nil {
		log.Println("Error retrieving link:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		return
	}
	http.Redirect(w, r, to, http.StatusFound)
}

func (s *server) ServeAdminHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	// TODO: get current user
	var user interface{}
	if user == nil {
		// TODO: generate login url
		var url string
		fmt.Fprintf(w, `<a href="%s">Sign in</a>`, url)
		return
	}
	if r.Method == "POST" {
		domain := strings.TrimSpace(r.PostFormValue("domain"))
		path := "/" + strings.Trim(strings.TrimSpace(r.PostFormValue("path")), "/")
		target := strings.TrimSpace(r.PostFormValue("target"))
		// TODO: set email to user's email
		var email string
		err := s.datastore.SetLink(r.Context(), domain, path, target, email)
		if err != nil {
			log.Println("Error setting link:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error"))
			return
		}
		http.Redirect(w, r, "https://"+s.adminDomain, http.StatusFound)
		return
	}
	tmpl, err := template.New("form").Parse(formTmpl)
	if err != nil {
		log.Println("Error parsing template:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}
	links, err := s.datastore.ListLinks(r.Context())
	if err != nil {
		log.Println("Error retrieving links:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}
	domains, err := s.datastore.ListDomains(r.Context())
	if err != nil {
		log.Println("Error retrieving domains:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}
	sort.Slice(domains, func(i, j int) bool {
		return domains[i].Domain < domains[j].Domain
	})
	data := TemplateData{
		Domains: domains,
		Links:   links,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error writing template:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}
}
