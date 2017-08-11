package couchdb

import "github.com/hashicorp/terraform/helper/schema"

func stringsFromSet(d interface{}) []string {
	s := d.(*schema.Set)
	ret := []string{}
	for _, v := range s.List() {
		ret = append(ret, v.(string))
	}
	return ret
}
