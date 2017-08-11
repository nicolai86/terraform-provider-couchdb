package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/nicolai86/terraform-provider-couchdb/couchdb"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: couchdb.Provider})
}
