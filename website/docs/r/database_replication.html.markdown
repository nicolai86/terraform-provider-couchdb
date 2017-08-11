---
layout: "couchdb"
page_title: "couchdb: database_replication"
sidebar_current: "docs-couchdb-resource-database_replication"
description: |-
  Manages couchdb database replications.
---

# couchdb\_database\_replication

Manages couchdb database replicationn settings.

## Example Usage

```hcl
resource "couchdb_database" "example" {
  name = "test"
}

resource "couchdb_database_replication" "example" {
  name          = "one-off"
  source        = "_users"
  target        = "users"
  create_target = true
  continuous    = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the replication.
* `source` - (Required) The source of the replication
* `target` - (Required) The target of the replicationn
* `create_target` - (Optional) Create the target if it does not exist (default: false)
* `continuous` - (Optional) Keep the replication running all the time (default: false)
* `filter` - (Optional) Document filter to use durign replication (default: none)
* `query_params` - (Optional) Additional params to pass along to filter function (default: {})
* `context` - (Optional) The execution context for replication 

*context* has the following attributes:

* `user` - (Optional) Execute replication with specific user privileges
* `roles` - (Optional) List of roles the replication is executed with

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource
* `name` - Name of the resource
* `source` - Source of the resource
* `target` - Target of the resource
* `create_target` - Create Target unless it exists
* `continuous` - Continuous replication enabled
* `filter` - Filter of the resource
* `query_params` - Query Params of the filter
* `context` - execution context of the resource
* `replication_id` - replication_id of the resource
* `replication_state` - replication_state of the resource
* `replication_state_reason` - replication_state_reason of the resource