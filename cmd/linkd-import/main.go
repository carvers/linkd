package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/carvers/linkd/storers"
	"google.golang.org/api/option"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: linkd-import FILE")
		os.Exit(1)
	}
	file := os.Args[1]
	mappings, err := loadMapping(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	domain := parseDomain(file)

	ctx := context.Background()

	project := os.Getenv("DATASTORE_PROJECT")
	creds := os.Getenv("DATASTORE_CREDS")

	var client *datastore.Client

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
	store := storers.NewDatastore(client)
	for from, to := range mappings {
		err := store.SetLink(ctx, domain, from, to)
		if err != nil {
			log.Printf("Error setting link %q to %q\n", domain+from, to)
			os.Exit(1)
		}
		log.Printf("Set link %q to %q\n", domain+from, to)
	}
}
