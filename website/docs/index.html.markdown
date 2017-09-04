---
layout: "couchdb"
page_title: "Provider: CouchDB"
sidebar_current: "docs-couchdb-index"
description: |-
  The CouchDB provider is used to interact with CouchDB.
---

# CouchDB Provider

The CouchDB provider is used to manage couchdb resources.

Use the navigation to the left to read about the available resources.

## Example Usage

Here is an example that will setup the following:
+ A database.

(create this as sl.tf and run terraform commands from this directory):

```hcl
provider "couchdb" {
  endpoint = "http://localhost:5984"
}

resource "couchdb_admin" "jenny" {
  name = "jenny"
  password = "secret" 
}

resource "couchdb_database" "db1" {
  name = "example"
}

resource "couchdb_database_replication" "db2db" {
  name = "example"
  source = "${couchdb_database.db1.name}"
  target = "example-clone"
  create_target = true
  continuous = true
}

resource "couchdb_database_design_document" "test" {
	database = "${couchdb_database.db1.name}"
	name = "types"

	view {
		name = "people"
		map = "function(doc) { if (doc.type == 'person') { emit(doc); } }"
	}
}
```

You'll need to provide your CouchDB endpoint, **username** and **password** for an administrator or privileged user.

If you do not want to put credentials in your configuration file,
you can leave them out and instead set these environment variables:

- **COUCHDB_ENDPOINT**: Your CouchDB endpoint
- **COUCHDB_USERNAME**: Your CouchDB admin username
- **COUCHDB_PASSWORD**: Your CouchDB admin password
