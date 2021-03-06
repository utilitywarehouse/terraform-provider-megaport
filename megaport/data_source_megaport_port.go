package megaport

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

var (
	megaportPorts []*api.Product
)

func dataSourceMegaportPort() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMegaportPortRead,

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

func dataSourceUpdatePorts(c *api.Client) error {
	megaportMutexKV.Lock("ports")
	defer megaportMutexKV.Unlock("ports")
	if megaportPorts != nil {
		return nil
	}
	log.Printf("[INFO] Updating port list")
	pp, err := c.ListPorts()
	if err != nil {
		return err
	}
	megaportPorts = pp
	return nil
}

func dataSourceMegaportPortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	if err := dataSourceUpdatePorts(cfg.Client); err != nil {
		return diag.FromErr(err)
	}
	var filtered []*api.Product
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		nr := regexp.MustCompile(nameRegex.(string))
		for _, port := range megaportPorts {
			if nr.MatchString(port.ProductName) {
				filtered = append(filtered, port)
			}
		}
	}
	if len(filtered) < 1 {
		return diag.FromErr(fmt.Errorf("No ports were found."))
	}
	if len(filtered) > 1 {
		return diag.FromErr(fmt.Errorf("Multiple ports were found. Please use a more specific query."))
	}
	d.SetId(filtered[0].ProductUid)
	return nil
}
