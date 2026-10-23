package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/carvers/linkd/storers"
	"google.golang.org/api/option"
)

var BuildID, GitSHA string

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	log.Println("Using build ID", BuildID, "and git SHA", GitSHA)
	ctx := context.Background()

	project := os.Getenv("DATASTORE_PROJECT")
	creds := os.Getenv("DATASTORE_CREDS")

	var client *datastore.Client
	var err error

	if project == "" {
		log.Println("DATASTORE_PROJECT must be set")
		os.Exit(1)
	}

	if creds != "" {
		client, err = datastore.NewClient(ctx, project, option.WithAuthCredentialsFile(option.ServiceAccount, creds))
	} else {
		client, err = datastore.NewClient(ctx, project)
	}
	if err != nil {
		log.Println("Error setting up datastore client:", err.Error())
		os.Exit(1)
	}

	s := server{
		datastore: storers.NewDatastore(client),
	}

	http.HandleFunc("/version.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Build", BuildID, "from", GitSHA)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.Handle("/", &s)
	err = http.ListenAndServe(":9876", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
