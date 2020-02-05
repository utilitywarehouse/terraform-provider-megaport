package megaport

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"connect_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "GOOGLE"}, false),
			},
			"location_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"vxc_permitted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
	log.Printf("[INFO] Updating partner port list")
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
	unfiltered := megaportPartnerPorts
	filtered := []*api.Megaport{}
	vp := d.Get("vxc_permitted")
	for _, port := range unfiltered {
		if port.VxcPermitted == vp.(bool) {
			filtered = append(filtered, port)
		}
	}
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		unfiltered = filtered
		filtered = []*api.Megaport{}
		nr := regexp.MustCompile(nameRegex.(string))
		for _, port := range unfiltered {
			if nr.MatchString(port.Title) {
				filtered = append(filtered, port)
			}
		}
	}
	if ct, ok := d.GetOk("connect_type"); ok {
		unfiltered = filtered
		filtered = []*api.Megaport{}
		for _, port := range unfiltered {
			if port.ConnectType == ct.(string) {
				filtered = append(filtered, port)
			}
		}
	}
	if lid, ok := d.GetOk("location_id"); ok {
		unfiltered = filtered
		filtered = []*api.Megaport{}
		for _, port := range unfiltered {
			if port.LocationId == uint64(lid.(int)) {
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
	d.SetId(filtered[0].ProductUid)
	return nil
}
