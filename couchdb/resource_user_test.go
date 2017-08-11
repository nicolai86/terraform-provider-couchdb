package couchdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	couchdb "github.com/nicolai86/couchdb-go"
)

func TestAccCouchDBUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCouchDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchDBUser,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBUserExists("couchdb_user.test"),
				),
			},
		},
	})
}

func testAccCouchDBUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("user ID is not set")
		}

		client := testAccProvider.Meta().(*Client).couch
		_, err := client.Users.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		return testAccCouchDBUserWorks(client.Host, rs.Primary.Attributes["name"], rs.Primary.Attributes["password"], "developer")
	}
}

func testAccCouchDBUserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchdb" {
			continue
		}

		client := testAccProvider.Meta().(*Client).couch
		_, err := client.Users.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if err == couchdb.ErrNotFound {
				return nil
			}
			return err
		}
	}

	return nil
}

var testAccCouchDBUser = `
resource "couchdb_user" "test" {
	name = "test"
	password = "test"

	roles = ["developer"]
}`
