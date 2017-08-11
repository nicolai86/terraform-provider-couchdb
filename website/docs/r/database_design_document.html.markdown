---
layout: "couchdb"
page_title: "couchdb: database_design_document"
sidebar_current: "docs-couchdb-resource-database_design_document"
description: |-
  Manages couchdb database design documents.
---

# couchdb\_database\_design\_document

Manages couchdb design documents and views for specific databases

## Example Usage

```hcl
resource "couchdb_database" "test" {
	name = "test"
}

resource "couchdb_database_design_document" "test" {
	database = "${couchdb_database.test.name}"
	name     = "test"

	view {
		name = "test"
		map  = "function(doc) { emit(doc); }"
	}
}
```

## Argument Reference

The following arguments are supported:

* `database` - (Required) The name of the database
* `name` - (Required) The name of the design document
* `language` - (Optional) The language of the map & reduce functions (default: javascript)
* `view` - (Optional) List of views inside the design document

*view* has the following attributes

* `name` - (Required) the name of the view
* `map` - (Required) the map function
* `reduce` - (Optional) the reduce function

## Attributes Reference

The following attributes are exported:

* `database` - database of the resource
* `name` - name of the resource
* `language` - language of the resource
* `view` - view of the resource
* `revision` - revision of the resource
