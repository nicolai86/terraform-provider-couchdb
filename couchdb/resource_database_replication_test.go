package couchdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	couchdb "github.com/nicolai86/couchdb-go"
)

func TestAccCouchDBDatabaseReplication_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCouchDBDatabaseReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchDBDatabaseReplication,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDatabaseReplicationExists("couchdb_database_replication.test"),
				),
			},
			{
				Config: testAccCouchDBDatabaseReplication_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDatabaseReplicationExists("couchdb_database_replication.test"),
				),
			},
		},
	})
}

func testAccCouchDBDatabaseReplicationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("replication ID is not set")
		}

		client := testAccProvider.Meta().(*Client).couch
		db := client.Database(rs.Primary.Attributes["database"])
		err := db.Get(context.Background(), rs.Primary.ID, "")

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCouchDBDatabaseReplicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchdb" {
			continue
		}

		client := testAccProvider.Meta().(*Client).couch
		db := client.Database(rs.Primary.Attributes["database"])
		err := db.Get(context.Background(), rs.Primary.ID, "")
		if err != nil {
			if err == couchdb.ErrNotFound {
				return nil
			}
			return err
		}
	}

	return nil
}

var testAccCouchDBDatabaseReplication = `
resource "couchdb_database" "test" {
	name = "test"
}

resource "couchdb_database_replication" "test" {
	name = "test"
	source = "${couchdb_database.test.name}"
	target = "bar"
	create_target = true
	continuous = true
} 
`

var testAccCouchDBDatabaseReplication_update = `
resource "couchdb_database" "test" {
	name = "test"
}

resource "couchdb_database_replication" "test" {
	name = "test"
	
	source = "${couchdb_database.test.name}"
	target = "bar"
	create_target = true
	continuous = true

	context {
		user = "admin"
	}

	filter = "documents/by_author"
	query_params {
		author = "alex"
	}
} 
`
