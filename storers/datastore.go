package storers

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

const (
	datastoreLinkKind   = "Link"
	datastoreDomainKind = "Domain"
)

type Datastore struct {
	client *datastore.Client
}

func NewDatastore(client *datastore.Client) Datastore {
	return Datastore{client: client}
}

type DatastoreLink struct {
	Domain    string
	Path      string
	Target    string
	CreatedBy string
	CreatedAt time.Time
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

func (d Datastore) ListLinks(ctx context.Context) ([]DatastoreLink, error) {
	var links []DatastoreLink
	iter := d.client.Run(ctx, datastore.NewQuery("Link"))
	for {
		var link DatastoreLink
		_, err := iter.Next(&link)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (d Datastore) SetLink(ctx context.Context, domain, path, target, email string) error {
	key := d.Key(domain, path)
	link := DatastoreLink{Target: target, Domain: domain, Path: path, CreatedAt: time.Now(), CreatedBy: email}
	_, err := d.client.Put(ctx, key, &link)
	if err != nil {
		return err
	}
	return nil
}

func (d Datastore) DeleteLink(ctx context.Context, domain, path string) error {
	return d.client.Delete(ctx, d.Key(domain, path))
}

type DatastoreDomain struct {
	Domain string
}

func (d Datastore) AddDomain(ctx context.Context, domain string) error {
	dom := DatastoreDomain{
		Domain: domain,
	}
	_, err := d.client.Put(ctx, datastore.NameKey(datastoreDomainKind, domain, nil), &dom)
	if err != nil {
		return err
	}
	return nil
}

func (d Datastore) DeleteDomain(ctx context.Context, domain string) error {
	return d.client.Delete(ctx, datastore.NameKey(datastoreDomainKind, domain, nil))
}

func (d Datastore) ListDomains(ctx context.Context) ([]DatastoreDomain, error) {
	var domains []DatastoreDomain
	iter := d.client.Run(ctx, datastore.NewQuery("Domain"))
	for {
		var domain DatastoreDomain
		_, err := iter.Next(&domain)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}
	return domains, nil
}
