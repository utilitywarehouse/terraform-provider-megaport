package megaport

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

var (
	megaportPartnerPorts []*api.Megaport
)

func dataSourceMegaportPartnerPort() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMegaportPartnerPortRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			// computed attributes
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUpdatePartnerPorts(c *api.Client) error {
	megaportMutexKV.Lock("partner_ports")
	defer megaportMutexKV.Unlock("partner_ports")
	if megaportPartnerPorts != nil {
		return nil
	}
	log.Printf("Updating partner port list")
	pp, err := c.GetMegaports() // TODO: rename in api
	if err != nil {
		return err
	}
	megaportPartnerPorts = pp
	return nil
}

func dataSourceMegaportPartnerPortRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	if err := dataSourceUpdatePartnerPorts(cfg.Client); err != nil {
		return err
	}
	var filtered []*api.Megaport
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		nr := regexp.MustCompile(nameRegex.(string))
		for _, port := range megaportPartnerPorts {
			if nr.MatchString(port.Title) {
				filtered = append(filtered, port)
			}
		}
	}
	if len(filtered) < 1 {
		return fmt.Errorf("No partner ports were found.")
	}
	if len(filtered) > 1 {
		return fmt.Errorf("Multiple partner ports were found. Please use a more specific query.")
	}
	d.SetId(newUUID(filtered[0].ProductUid)) // TODO: simply use the uuid?
	d.Set("uid", filtered[0].ProductUid)
	return nil
}
