package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"carvers.dev/linkd/storers"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	ctx := context.Background()

	project := os.Getenv("DATASTORE_PROJECT")
	creds := os.Getenv("DATASTORE_CREDS")
	name := os.Getenv("ADMIN_PANEL_NAME")
	adminDomain := os.Getenv("ADMIN_PANEL_DOMAIN")

	var client *datastore.Client
	var err error

	if project == "" {
		log.Println("DATASTORE_PROJECT must be set")
		os.Exit(1)
	}

	if name == "" {
		log.Println("ADMIN_PANEL_NAME must be set")
		os.Exit(1)
	}
	if adminDomain == "" {
		log.Println("ADMIN_DOMAIN must be set")
		os.Exit(1)
	}

	if creds != "" {
		client, err = datastore.NewClient(ctx, project, option.WithServiceAccountFile(creds))
	} else {
		client, err = datastore.NewClient(ctx, project)
	}
	if err != nil {
		log.Println("Error setting up datastore client:", err.Error())
		os.Exit(1)
	}

	s := server{
		datastore:   storers.NewDatastore(client),
		adminDomain: adminDomain,
		name:        name,
	}

	http.Handle("/", &s)
	err = http.ListenAndServe(":9876", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
