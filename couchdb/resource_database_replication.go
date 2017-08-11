package couchdb

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	couchdb "github.com/nicolai86/couchdb-go"
)

func resourceDatabaseReplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseReplicationCreate,
		Read:   resourceDatabaseReplicationRead,
		Update: resourceDatabaseReplicationUpdate,
		Delete: resourceDatabaseReplicationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the replication document",
			},
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source of the replication",
			},
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target of the replication",
			},
			"create_target": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Create target if it does not exist?",
			},
			"continuous": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Keep the replication permanently running?",
			},
			"context": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true, // Default to []string{}
				Description: "Execution context of the replication. Can be a specific user or a set of roles",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Execution context username",
						},
						"roles": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Execution context roles",
						},
					},
				},
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter documents when replicating",
			},
			"query_params": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Additional query parameters for the filter function",
			},
			"replication_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal replication ID",
			},
			"replication_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "currennt replication state",
			},
			"replication_state_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "current replication state transition reason",
			},
		},
	}
}

func stringMap(d interface{}) map[string]string {
	ps, ok := d.(map[string]interface{})
	if !ok {
		return nil
	}
	params := map[string]string{}
	for k, v := range ps {
		params[k] = v.(string)
	}
	return params
}

func extractUserContext(d interface{}) *couchdb.UserContext {
	ctx, ok := d.(map[string]interface{})
	if !ok {
		return nil
	}

	return &couchdb.UserContext{
		Name:  ctx["user"].(string),
		Roles: stringsFromSet(ctx["roles"]),
	}
}

func resourceDatabaseReplicationCreate(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	var authorization *couchdb.UserContext
	if v, ok := d.GetOk("context"); ok {
		vs := v.([]interface{})
		if len(vs) == 1 {
			authorization = extractUserContext(vs[0])
		}
	}
	payload := couchdb.ReplicationPayload{
		ID:           d.Get("name").(string),
		Source:       d.Get("source").(string),
		Target:       d.Get("target").(string),
		CreateTarget: d.Get("create_target").(bool),
		Continuous:   d.Get("continuous").(bool),
		Filter:       d.Get("filter").(string),
		QueryParams:  stringMap(d.Get("query_params")),
		Context:      authorization,
	}
	rep, err := couch.Replications.Create(context.Background(), payload)
	if err != nil {
		return err
	}
	d.SetId(rep.ID)

	return resourceDatabaseReplicationRead(d, m)
}

func resourceDatabaseReplicationRead(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	rep, err := couch.Replications.Get(context.Background(), d.Id())
	if err != nil {
		return err
	}

	d.Set("source", rep.Source)
	d.Set("target", rep.Target)
	d.Set("continuous", rep.Continuous)
	d.Set("create_target", rep.CreateTarget)
	d.Set("filter", rep.Filter)
	d.Set("query_params", rep.QueryParams)
	d.Set("replication_id", rep.ReplicationID)
	d.Set("replication_state", rep.ReplicationState)
	d.Set("replication_state_reason", rep.ReplicationStateReason)

	return nil
}

func resourceDatabaseReplicationUpdate(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	var authorization *couchdb.UserContext
	if v, ok := d.GetOk("context"); ok {
		vs := v.([]interface{})
		if len(vs) == 1 {
			authorization = extractUserContext(vs[0])
		}
	}
	payload := couchdb.ReplicationPayload{
		ID:           d.Id(),
		Source:       d.Get("source").(string),
		Target:       d.Get("target").(string),
		CreateTarget: d.Get("create_target").(bool),
		Continuous:   d.Get("continuous").(bool),
		Filter:       d.Get("filter").(string),
		QueryParams:  stringMap(d.Get("query_params")),
		Context:      authorization,
	}
	_, err := couch.Replications.Update(context.Background(), payload)
	if err != nil {
		return err
	}

	return resourceDatabaseReplicationRead(d, m)
}

func resourceDatabaseReplicationDelete(d *schema.ResourceData, m interface{}) error {
	couch := m.(*Client).couch

	err := couch.Replications.Delete(context.Background(), d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
