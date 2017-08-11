---
layout: "couchdb"
page_title: "couchdb: admin_user"
sidebar_current: "docs-couchdb-resource-admin_user"
description: |-
  Manages couchdb admins.
---

# couchdb\_admin\_user

Manages couchdb administrators.

## Example Usage

```hcl
resource "couchdb_admin_user" "jenny" {
  name     = "jenny"
  password = "secret"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the administrator. Used to log-in.
* `password` - (Required) The password of the administrator.
* `node` - (Optional) The node on which to create this user. Ignored for CouchDB 1.6, required for CouchDB 2.x

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource
* `name` - Name of the resource
* `password` - Password of the resource
