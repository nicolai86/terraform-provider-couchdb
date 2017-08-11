package couchdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	couchdb "github.com/nicolai86/couchdb-go"
)

func TestAccCouchDBAdminUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCouchDBAdminUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchDBAdminUser_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCouchDBAdminUser("couchdb_admin_user.test"),
				),
			},
		},
	})
}

func testAccCouchDBAdminUser(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("admin user ID is not set")
		}

		client := testAccProvider.Meta().(*Client).couch
		memberships, err := client.Membership()
		opts := couchdb.ClusterOptions{}
		if err == nil {
			opts.Node = memberships.AllNodes[0]
		}
		admins, err := client.Admins.List(context.Background(), opts)

		if err != nil {
			return fmt.Errorf("Failed to list admins: %v", err.Error())
		}

		for _, admin := range admins {
			if admin == rs.Primary.ID {
				return testAccCouchDBUserWorks(client.Host, rs.Primary.Attributes["name"], rs.Primary.Attributes["password"], "_admin")
			}
		}

		return fmt.Errorf("Admin %s does not exist", rs.Primary.ID)
	}
}

func testAccCouchDBAdminUserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchdb" {
			continue
		}

		client := testAccProvider.Meta().(*Client).couch
		memberships, err := client.Membership()
		opts := couchdb.ClusterOptions{}
		if err == nil {
			opts.Node = memberships.AllNodes[0]
		}
		admins, err := client.Admins.List(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("Failed to list admins: %v", err.Error())
		}

		for _, admin := range admins {
			if admin == rs.Primary.ID {
				return fmt.Errorf("Admin %s still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

var testAccCouchDBAdminUser_basic = `
resource "couchdb_admin_user" "test" {
	name = "test"
	password = "test"
	
	node = "nonode@nohost"
}`
