package megaport

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationsRead,

		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceLocationsRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	d.SetId(time.Now().UTC().String())
	loc, err := cfg.Client.GetLocations()
	if err != nil {
		return err
	}
	names := make([]string, len(loc), len(loc))
	ids := make([]uint64, len(loc), len(loc))
	for i, v := range loc {
		names[i] = v.Name
		ids[i] = v.Id
	}
	if err := d.Set("names", names); err != nil {
		return fmt.Errorf("Error setting Location names: %s", err)
	}
	if err := d.Set("ids", ids); err != nil {
		return fmt.Errorf("Error setting Location IDs: %s", err)
	}
	return nil
}
