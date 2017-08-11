package couchdb

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	couchdb "github.com/nicolai86/couchdb-go"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"couchdb": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("COUCHDB_ENDPOINT"); v == "" {
		t.Fatal("COUCHDB_ENDPOINT must be set for acceptance tests")
	}
}

func testAccCouchDBUserWorks(endpoint, username, password, expectedRole string) error {
	client, err := couchdb.New(endpoint, &http.Client{}, couchdb.WithBasicAuthentication(username, password))
	if err != nil {
		return err
	}

	sess, err := client.Sessions.Get(context.Background())
	if err != nil {
		return err
	}
	if sess.Context.Name != username {
		return fmt.Errorf("Expected user %s, but got %s", username, sess.Context.Name)
	}
	if sess.Context.Roles[0] != expectedRole {
		return fmt.Errorf("Expected user role %s, but got %s", expectedRole, sess.Context.Roles)
	}
	return nil
}
