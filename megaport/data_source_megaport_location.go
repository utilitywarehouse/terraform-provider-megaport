package megaport

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

var (
	megaportLocations []*api.Location
)

func dataSourceMegaportLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMegaportLocationRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
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

func dataSourceMegaportLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	if err := dataSourceUpdateLocations(cfg.Client); err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("No locations were found."))
	}
	if len(filtered) > 1 {
		return diag.FromErr(fmt.Errorf("Multiple locations were found. Please use a more specific query."))
	}
	d.SetId(strconv.FormatUint(filtered[0].Id, 10))
	return nil
}
