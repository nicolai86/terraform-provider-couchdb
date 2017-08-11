---
layout: "couchdb"
page_title: "couchdb: user"
sidebar_current: "docs-couchdb-resource-user"
description: |-
  Manages couchdb users.
---

# couchdb\_user

Manages couchdb users with limited privileges.

## Example Usage

```hcl
resource "couchdb_user" "jenny" {
  name     = "jenny"
  password = "secret"
  roles    = ["developers"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user. Used to log-in.
* `password` - (Required) The password of the user.
* `roles` - (Optional) The roles of the user. Used for per-database access control.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource
* `name` - Name of the resource
* `password` - Password of the resource
* `roles` - Roles of the resource
* `revision` - Revision of the resource
