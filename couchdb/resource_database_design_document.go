package couchdb

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	couchdb "github.com/nicolai86/couchdb-go"
)

func resourceDesignDocument() *schema.Resource {
	return &schema.Resource{
		Create: resourceDesignDocumentCreate,
		Read:   resourceDesignDocumentRead,
		Update: resourceDesignDocumentUpdate,
		Delete: resourceDesignDocumentDelete,

		Schema: map[string]*schema.Schema{
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Database to associate design with",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the design document",
			},
			"revision": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Revision",
			},
			"language": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "javascript",
				Description: "Language of map/ reduce functions",
			},
			"view": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A view inside the design document",
				Set: func(v interface{}) int {
					view := v.(map[string]interface{})
					name := view["name"].(string)
					id := 0
					for _, b := range md5.Sum([]byte(name)) {
						id += int(b)
					}
					return id
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the view",
						},
						"map": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Map function",
						},
						"reduce": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Reduce functionn",
						},
					},
				},
			},
		},
	}
}

func resourceDesignDocumentCreate(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	db := couch.Database(d.Get("database").(string))

	doc := couchdb.DesignDocument{
		Language: d.Get("language").(string),
		Views:    map[string]couchdb.View{},
	}
	if vs, ok := d.GetOk("view"); ok {
		views := vs.(*schema.Set)
		for _, v := range views.List() {
			view := v.(map[string]interface{})
			doc.Views[view["name"].(string)] = couchdb.View{
				MapFn:    view["map"].(string),
				ReduceFn: view["reduce"].(string),
			}
		}
	}
	id := fmt.Sprintf("_design/%s", d.Get("name").(string))
	rev, err := db.Put(context.Background(), id, doc)
	if err != nil {
		return err
	}
	d.SetId(id)
	d.Set("revision", rev)

	return resourceDesignDocumentRead(d, m)
}

func resourceDesignDocumentRead(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	db := couch.Database(d.Get("database").(string))
	doc := couchdb.DesignDocument{}
	err := db.Get(context.Background(), d.Id(), &doc)
	if err != nil {
		return err
	}

	d.Set("language", doc.Language)
	views := []map[string]string{}
	for name, view := range doc.Views {
		v := map[string]string{
			"name":   name,
			"map":    view.MapFn,
			"reduce": view.ReduceFn,
		}
		views = append(views, v)
	}

	d.Set("view", views)
	d.Set("revision", doc.Rev)

	return nil
}

func resourceDesignDocumentUpdate(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	db := couch.Database(d.Get("database").(string))
	doc := couchdb.DesignDocument{
		Document: couchdb.Document{
			ID:  d.Id(),
			Rev: d.Get("revision").(string),
		},
		Language: d.Get("language").(string),
		Views:    map[string]couchdb.View{},
	}
	if vs, ok := d.GetOk("view"); ok {
		views := vs.(*schema.Set)
		for _, v := range views.List() {
			view := v.(map[string]interface{})
			doc.Views[view["name"].(string)] = couchdb.View{
				MapFn:    view["map"].(string),
				ReduceFn: view["reduce"].(string),
			}
		}
	}
	rev, err := db.Put(context.Background(), d.Id(), doc)
	if err != nil {
		return err
	}
	d.Set("revision", rev)

	return nil
}

func resourceDesignDocumentDelete(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	db := couch.Database(d.Get("database").(string))
	_, err := db.Delete(context.Background(), d.Id(), d.Get("revision").(string))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
