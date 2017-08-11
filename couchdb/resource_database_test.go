package couchdb

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCouchDBDatabase_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCouchDBDatabaseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCouchDBDatabase,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDatabaseExists("couchdb_database.test"),
				),
			},
			resource.TestStep{
				Config: testAccCouchDBDatabase_security,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDatabaseSecurity("couchdb_database.test"),
				),
			},
			resource.TestStep{
				Config: testAccCouchDBDatabase,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBDatabaseSecurity("couchdb_database.test"),
				),
			},
		},
	})
}

func testAccCouchDBDatabaseExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("database ID is set")
		}

		client := testAccProvider.Meta().(*Client).couch
		ok, err := client.Databases.Exists(rs.Primary.ID)

		if err != nil {
			return err
		}

		if !ok {
			return fmt.Errorf("Record not found")
		}

		return nil
	}
}

func testAccCouchDBDatabaseSecurity(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("database ID is not set")
		}

		client := testAccProvider.Meta().(*Client).couch
		db := client.Database(rs.Primary.ID)
		sec, err := db.GetSecurity(context.Background())

		if err != nil {
			return err
		}

		if rs.Primary.Attributes["security.0.members.#"] != strconv.Itoa(len(sec.Members.Names)) {
			return fmt.Errorf("Expected %d members, got %s", len(sec.Members.Names), rs.Primary.Attributes["security.0.members.#"])
		}
		if rs.Primary.Attributes["security.0.member_roles.#"] != strconv.Itoa(len(sec.Members.Roles)) {
			return fmt.Errorf("Expected %d member roles, got %s", len(sec.Members.Roles), rs.Primary.Attributes["security.0.member_roles.#"])
		}
		if rs.Primary.Attributes["security.0.admins.#"] != strconv.Itoa(len(sec.Admins.Names)) {
			return fmt.Errorf("Expected %d admins, got %s", len(sec.Admins.Names), rs.Primary.Attributes["security.0.admins.#"])
		}
		if rs.Primary.Attributes["security.0.admin_roles.#"] != strconv.Itoa(len(sec.Admins.Roles)) {
			return fmt.Errorf("Expected %d admin roles, got %s", len(sec.Admins.Roles), rs.Primary.Attributes["security.0.admin_roles.#"])
		}
		return nil
	}
}

func testAccCouchDBDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).couch

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchdb" {
			continue
		}

		ok, err := client.Databases.Exists(rs.Primary.ID)

		if err == nil || ok {
			return fmt.Errorf("DB still exists")
		}
	}

	return nil
}

var testAccCouchDBDatabase_security = `
resource "couchdb_database" "test" {
	name = "test"

	security {
		admins  = ["admin"]
		members = ["max"]
		
		admin_roles  = ["maintenance"]
		member_roles = ["users"]
	}
}
`

var testAccCouchDBDatabase = `
resource "couchdb_database" "test" {
	name = "test"
	
	security {}

	clustering {
		shards   = 6
		replicas = 2
	}
}
`
