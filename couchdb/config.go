package couchdb

import (
	"net/http"

	couchdb "github.com/nicolai86/couchdb-go"
)

// Config contains couchdb configuration values
type Config struct {
	Endpoint string
	Username string
	Password string
}

// Client contains couchdb api clients
type Client struct {
	couch *couchdb.Client
}

// Client configures and returns a fully initialized couchdb client
func (c *Config) Client() (*Client, error) {
	authentication := func(c *couchdb.Client) error { return nil }
	if c.Username != "" && c.Password != "" {
		authentication = couchdb.WithBasicAuthentication(c.Username, c.Password)
	}
	client, err := couchdb.New(c.Endpoint, &http.Client{}, authentication)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}
