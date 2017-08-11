---
layout: "couchdb"
page_title: "couchdb: database"
sidebar_current: "docs-couchdb-resource-database"
description: |-
  Manages couchdb databases.
---

# couchdb\_database

Manages couchdb databases.

## Example Usage

```hcl
resource "couchdb_database" "example" {
  name = "test"

  security {
    admins      = ["jenny"]
    admin_roles = ["owners"]

    members      = ["jan"]
    member_roles = ["developers"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database
* `security` - (Optional) The access control for the database
* `clustering` - (Optional) The couchdb 2.x clustering configuration

*security* has the following attributes:

* `admins` - (Optional) list of usernames who can administer this database
* `admin_roles` - (Optional) list of roles who can administer this database
* `members` - (Optional) list of usernames who can access this database
* `member_roles` - (Optional) list of roles who can access this database

*clustering* has the following attributes

* `replicas` - (Optional) the number of replicas for the database. Default: 3
* `shards` - (Optional) the number of shards for the database. Default: 8

## Attributes Reference

The following attributes are exported:

* `id` - ID of the new resource
* `name` - Name of the new resource
* `security` - Security configuration for the resource
* `document_count` - Number of documents inside the resource
* `document_deletion_count` - Number of tombstones inside the resource
* `disk_size` - Disk size of the resource
* `disk_format_version` - Disk format version of the resource
* `data_size` - Data size of the resource
* `update_sequence_number` - Update sequence number of the resource
* `purge_sequence_number` - Purge sequence number of the resource
* `committed_update_sequence_number` - Comitted updated sequence number of the resource