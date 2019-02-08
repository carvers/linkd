package storers

import (
	"context"
	"strings"

	"cloud.google.com/go/datastore"
)

const datastoreLinkKind = "Link"

type Datastore struct {
	client *datastore.Client
}

func NewDatastore(client *datastore.Client) Datastore {
	return Datastore{client: client}
}

type DatastoreLink struct {
	Domain string
	Path   string
	Target string
}

func (d Datastore) Key(domain, path string) *datastore.Key {
	key := datastore.NameKey(datastoreLinkKind, strings.Trim(domain, "/")+"/"+strings.Trim(path, "/"), nil)
	return key
}

func (d Datastore) GetLink(ctx context.Context, domain, path string) (string, error) {
	key := d.Key(domain, path)
	var link DatastoreLink
	err := d.client.Get(ctx, key, &link)
	if err == datastore.ErrNoSuchEntity {
		return "", ErrLinkNotFound
	} else if err != nil {
		return "", err
	}
	return link.Target, nil
}

func (d Datastore) SetLink(ctx context.Context, domain, path, target string) error {
	key := d.Key(domain, path)
	link := DatastoreLink{Target: target, Domain: domain, Path: path}
	_, err := d.client.Put(ctx, key, &link)
	if err != nil {
		return err
	}
	return nil
}

func (d Datastore) DeleteLink(ctx context.Context, domain, path string) error {
	return d.client.Delete(ctx, d.Key(domain, path))
}
