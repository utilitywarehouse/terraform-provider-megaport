package megaport

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceMegaports() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMegaportsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceMegaportsRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	d.SetId(time.Now().UTC().String())
	loc, err := cfg.Client.GetMegaports()
	if err != nil {
		return err
	}
	ids := make([]string, len(loc), len(loc))
	for i, v := range loc {
		ids[i] = v.ProductUid
	}
	if err := d.Set("ids", ids); err != nil {
		return fmt.Errorf("Error setting Megaport Product UIDs: %s", err)
	}
	return nil
}
