package couchdb

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	couchdb "github.com/nicolai86/couchdb-go"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseCreate,
		Read:   resourceDatabaseRead,
		Update: resourceDatabaseUpdate,
		Delete: resourceDatabaseDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the database",
			},
			"security": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Security configuration of the database",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admins": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Database administrators",
						},
						"admin_roles": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Database administration roles",
						},
						"members": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Database members",
						},
						"member_roles": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Database member roles",
						},
					},
				},
			},
			"clustering": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "database clustering configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replicas": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     3,
							Description: "Number of replicas",
						},
						"shards": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     8,
							Description: "Number of shards",
						},
					},
				},
			},
			"document_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of documents in database",
			},
			"document_deletion_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tombstones in database",
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of storage disk",
			},
			"disk_format_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Disk format version",
			},
			"data_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of database data",
			},
		},
	}
}

func extractDatabaseSecurity(d interface{}) couchdb.DatabaseSecurity {
	sec := couchdb.DatabaseSecurity{
		Admins: couchdb.AuthorizationRules{
			Names: []string{},
			Roles: []string{},
		},
		Members: couchdb.AuthorizationRules{
			Names: []string{},
			Roles: []string{},
		},
	}

	security, ok := d.(map[string]interface{})
	if !ok {
		return sec
	}

	sec.Admins.Names = stringsFromSet(security["admins"])
	sec.Admins.Roles = stringsFromSet(security["admin_roles"])
	sec.Members.Names = stringsFromSet(security["members"])
	sec.Members.Roles = stringsFromSet(security["member_roles"])
	return sec
}

func extractClusterOptions(v interface{}) (ret couchdb.DatabaseClusterOptions) {
	vs := v.([]interface{})
	if len(vs) != 1 {
		return ret
	}
	vi := vs[0].(map[string]interface{})
	ret.Replicas = vi["replicas"].(int)
	ret.Shards = vi["shards"].(int)
	return ret
}

func resourceDatabaseCreate(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	err := couch.Databases.Create(d.Get("name").(string), extractClusterOptions(d.Get("clustering")))
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	if v, ok := d.GetOk("security"); ok {
		vs := v.([]interface{})
		if len(vs) == 1 {
			db := couch.Database(d.Get("name").(string))
			err := db.SetSecurity(context.Background(), extractDatabaseSecurity(vs[0]))
			if err != nil {
				return err
			}
		}
	}

	return resourceDatabaseRead(d, m)
}

func resourceDatabaseUpdate(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	if d.HasChange("security") {
		db := couch.Database(d.Get("name").(string))
		if v, ok := d.GetOk("security"); ok {
			vs := v.([]interface{})
			if len(vs) == 1 {
				err := db.SetSecurity(context.Background(), extractDatabaseSecurity(vs[0]))
				if err != nil {
					return err
				}
			}
		} else {
			err := db.SetSecurity(context.Background(), extractDatabaseSecurity(nil))
			if err != nil {
				return err
			}
		}
	}

	return resourceDatabaseRead(d, m)
}

func resourceDatabaseRead(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	meta, err := couch.Databases.Meta(d.Id())
	if err != nil {
		return err
	}

	d.Set("document_count", strconv.Itoa(meta.DocumentCount))
	d.Set("document_deletion_count", strconv.Itoa(meta.DocumentDeletionCount))
	d.Set("disk_size", strconv.Itoa(meta.DiskSize))
	d.Set("data_size", strconv.Itoa(meta.DataSize))
	d.Set("disk_format_versionn", strconv.Itoa(meta.DiskFormatVersion))

	db := couch.Database(d.Get("name").(string))
	sec, err := db.GetSecurity(context.Background())
	if err != nil {
		return err
	}
	security := []map[string][]string{
		{
			"admins":       sec.Admins.Names,
			"admin_roles":  sec.Admins.Roles,
			"members":      sec.Members.Names,
			"member_roles": sec.Members.Roles,
		},
	}
	d.Set("security", security)

	return nil
}

func resourceDatabaseDelete(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	err := couch.Databases.Delete(d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
