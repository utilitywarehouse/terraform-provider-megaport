package megaport

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceInternetExchanges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInternetExchangesRead,

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceInternetExchangesRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	d.SetId(time.Now().UTC().String())
	loc, err := cfg.Client.GetInternetExchanges(uint64(d.Get("location_id").(int)))
	if err != nil {
		return err
	}
	names := make([]string, len(loc), len(loc))
	for i, v := range loc {
		names[i] = v.Name
	}
	if err := d.Set("names", names); err != nil {
		return fmt.Errorf("Error setting Internet Exchanges names: %s", err)
	}
	return nil
}
