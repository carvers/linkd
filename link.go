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

	var client *datastore.Client
	var err error

	if project == "" {
		log.Println("DATASTORE_PROJECT must be set")
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
		datastore: storers.NewDatastore(client),
	}

	http.Handle("/", &s)
	err = http.ListenAndServe(":9876", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
