package megaport

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

var (
	megaportLocations []*api.Location
)

func dataSourceMegaportLocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMegaportLocationRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			// computed attributes
			"location_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceUpdateLocations(c *api.Client) error {
	megaportMutexKV.Lock("locations")
	defer megaportMutexKV.Unlock("locations")
	if megaportLocations != nil {
		return nil
	}
	log.Printf("Updating location list")
	loc, err := c.GetLocations()
	if err != nil {
		return err
	}
	megaportLocations = loc
	return nil
}

func dataSourceMegaportLocationRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	if err := dataSourceUpdateLocations(cfg.Client); err != nil {
		return err
	}
	var filtered []*api.Location
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		nr := regexp.MustCompile(nameRegex.(string))
		for _, loc := range megaportLocations {
			if nr.MatchString(loc.Name) {
				filtered = append(filtered, loc)
			}
		}
	}
	if len(filtered) < 1 {
		return fmt.Errorf("No locations were found.")
	}
	if len(filtered) > 1 {
		return fmt.Errorf("Multiple locations were found. Please use a more specific query.")
	}
	d.SetId(newUUID(strconv.FormatUint(filtered[0].Id, 10)))
	d.Set("location_id", filtered[0].Id)
	return nil
}
