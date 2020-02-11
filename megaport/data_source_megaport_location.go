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
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"mcr_available": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2}),
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
	log.Printf("[INFO] Updating location list")
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
	if mcrAvailable, ok := d.GetOk("mcr_available"); ok {
		unfiltered := filtered
		filtered = []*api.Location{}
		for _, loc := range unfiltered {
			if loc.Products.McrVersion == uint64(mcrAvailable.(int)) {
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
	d.SetId(strconv.FormatUint(filtered[0].Id, 10))
	return nil
}
