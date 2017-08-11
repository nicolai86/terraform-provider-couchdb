package couchdb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	couchdb "github.com/nicolai86/couchdb-go"
)

func resourceAdminUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAdminUserCreate,
		Read:   resourceAdminUserRead,
		Update: resourceAdminUserUpdate,
		Delete: resourceAdminUserDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the administrator",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password of the administrator",
			},
			"node": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Node to create admin user one (couchdb 2.x only)",
				Default:     "",
			},
		},
	}
}

func resourceAdminUserCreate(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	err := couch.Admins.Create(context.Background(), d.Get("name").(string), d.Get("password").(string), couchdb.ClusterOptions{
		Node: d.Get("node").(string),
	})
	if err != nil {
		return err
	}
	d.SetId(d.Get("name").(string))
	return nil
}

func resourceAdminUserRead(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	admins, err := couch.Admins.List(context.Background(), couchdb.ClusterOptions{
		Node: d.Get("node").(string),
	})
	if err != nil {
		return err
	}
	exists := false
	for _, name := range admins {
		exists = exists || name == d.Id()
	}
	if !exists {
		return fmt.Errorf("Admin %s was not found", d.Id())
	}

	return nil
}

func resourceAdminUserUpdate(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	err := couch.Admins.Update(context.Background(), d.Get("name").(string), d.Get("password").(string), couchdb.ClusterOptions{
		Node: d.Get("node").(string),
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceAdminUserDelete(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	err := couch.Admins.Delete(context.Background(), d.Id(), couchdb.ClusterOptions{
		Node: d.Get("node").(string),
	})
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
