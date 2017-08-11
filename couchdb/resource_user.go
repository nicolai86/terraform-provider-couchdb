package couchdb

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	couchdb "github.com/nicolai86/couchdb-go"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"revision": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Revision",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Username",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password",
			},
			"roles": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User roles",
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	user, err := couch.Users.Create(context.Background(), couchdb.CreateUserPayload{
		Name:     d.Get("name").(string),
		Password: d.Get("password").(string),
		Roles:    stringsFromSet(d.Get("roles")),
	})
	if err != nil {
		return err
	}
	d.SetId(user.ID)
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	user, err := couch.Users.Get(context.Background(), d.Id())
	if err != nil {
		return err
	}
	d.Set("revision", user.Rev)
	d.Set("roles", user.Roles)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	_, err := couch.Users.Update(context.Background(), couchdb.UpdateUserPayload{
		ID:       d.Id(),
		Name:     d.Get("name").(string),
		Password: d.Get("password").(string),
		Roles:    stringsFromSet(d.Get("roles")),
	})
	if err != nil {
		return err
	}

	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	err := couch.Users.Delete(context.Background(), d.Id())
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
