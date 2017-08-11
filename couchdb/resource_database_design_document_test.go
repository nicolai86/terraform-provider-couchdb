package couchdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	couchdb "github.com/nicolai86/couchdb-go"
)

func TestAccCouchDBDesignDocument_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCouchDBDesignDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchDBDesignDocument,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDesignDocumentExists("couchdb_database_design_document.test"),
				),
			},
			{
				Config: testAccCouchDBDesignDocument_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDesignDocumentExists("couchdb_database_design_document.test"),
				),
			},
		},
	})
}

func testAccCouchDBDesignDocumentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("database design document ID is not set")
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

func testAccCouchDBDesignDocumentDestroy(s *terraform.State) error {
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

var testAccCouchDBDesignDocument = `
resource "couchdb_database" "test" {
	name = "test"
}

resource "couchdb_database_design_document" "test" {
	database = "${couchdb_database.test.name}"
	name = "test"

	view {
		name = "test"
		map = "function(doc) { emit(doc._id, doc); }"
	}
}
`
var testAccCouchDBDesignDocument_update = `
resource "couchdb_database" "test" {
	name = "test"
}

resource "couchdb_database_design_document" "test" {
	database = "${couchdb_database.test.name}"
	name = "test"

	view {
		name = "cat"
		map = "function(doc) { emit(doc._id, doc); }"
	}

	view {
		name = "test"
		map = "function(doc) { emit(doc._id, doc); }"
	}
}
`
