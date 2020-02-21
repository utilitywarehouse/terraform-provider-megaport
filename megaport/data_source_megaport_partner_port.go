package megaport

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-uuid"
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
			"aws": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"aws", "gcp", "marketplace"},
				Elem:         dataSourceMegaportPartnerPortMarketplace(),
			},
			"gcp": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"aws", "gcp", "marketplace"},
				Elem:         dataSourceMegaportPartnerPortGcp(),
			},
			"marketplace": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"aws", "marketplace"},
				Elem:          dataSourceMegaportPartnerPortMarketplace(),
			},
			"bandwidths": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceMegaportPartnerPortMarketplace() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"vxc_permitted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceMegaportPartnerPortGcp() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pairing_key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[[:xdigit:]]{8}-([[:xdigit:]]{4}-){3}[[:xdigit:]]{12}\/[\w]+-[\w]+\d\/\d$`), "Invalid GCP pairing key format"),
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
	nameRegex := d.Get("name_regex").(string)
	if v, ok := d.GetOk("aws"); ok {
		if err := dataSourceUpdatePartnerPorts(cfg.Client); err != nil {
			return err
		}
		p, err := filterPartnerPorts(megaportPartnerPorts, "AWS", nameRegex, expandFilters(v))
		if err != nil {
			return err
		}
		if err := d.Set("bandwidths", []int{}); err != nil {
			return err
		}
		d.SetId(p.ProductUid)
		return nil
	}
	if v, ok := d.GetOk("marketplace"); ok {
		if err := dataSourceUpdatePartnerPorts(cfg.Client); err != nil {
			return err
		}
		p, err := filterPartnerPorts(megaportPartnerPorts, "DEFAULT", nameRegex, expandFilters(v))
		if err != nil {
			return err
		}
		if err := d.Set("bandwidths", []int{}); err != nil {
			return err
		}
		d.SetId(p.ProductUid)
		return nil
	}
	if v, ok := d.GetOk("gcp"); ok {
		pk := expandFilters(v)["pairing_key"].(string)
		randomUUID, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}
		pk = randomUUID + "/" + strings.SplitN(pk, "/", 2)[1]
		ports, bandwidths, err := cfg.Client.GetMegaportsForGcpPairingKey(pk)
		if err != nil {
			return err
		}
		p, err := filterCloudPartnerPorts(ports, nameRegex)
		if err != nil {
			return err
		}
		bw := make([]int, len(bandwidths))
		for i, v := range bandwidths {
			bw[i] = int(v)
		}
		if err := d.Set("bandwidths", bw); err != nil {
			return err
		}
		d.SetId(p.ProductUid)
		return nil
	}
	return nil
}

func expandFilters(v interface{}) map[string]interface{} {
	return v.([]interface{})[0].(map[string]interface{})
}

func filterPartnerPorts(ports []*api.Megaport, connectType, nameRegex string, d map[string]interface{}) (*api.Megaport, error) {
	unfiltered := ports
	filtered := []*api.Megaport{}
	for _, port := range unfiltered {
		if port.ConnectType == connectType {
			filtered = append(filtered, port)
		}
	}
	unfiltered = filtered
	filtered = []*api.Megaport{}
	nr := regexp.MustCompile(nameRegex)
	for _, port := range unfiltered {
		if nr.MatchString(port.Title) {
			filtered = append(filtered, port)
		}
	}
	if lid, ok := d["location_id"]; ok {
		unfiltered = filtered
		filtered = []*api.Megaport{}
		for _, port := range unfiltered {
			if port.LocationId == uint64(lid.(int)) {
				filtered = append(filtered, port)
			}
		}
	}
	if vp, ok := d["vxc_permitted"]; ok {
		unfiltered = filtered
		filtered = []*api.Megaport{}
		for _, port := range unfiltered {
			if port.VxcPermitted == vp.(bool) {
				filtered = append(filtered, port)
			}
		}
	}
	if len(filtered) < 1 {
		return nil, fmt.Errorf("No ports were found. You might want to use a less specific query.")
	}
	if len(filtered) > 1 {
		return nil, fmt.Errorf("Multiple ports were found. Please use a more specific query.")
	}
	return filtered[0], nil
}

func filterCloudPartnerPorts(ports []*api.MegaportCloud, nameRegex string) (*api.MegaportCloud, error) {
	filtered := []*api.MegaportCloud{}
	nr := regexp.MustCompile(nameRegex)
	for _, port := range ports {
		if nr.MatchString(port.Name) {
			filtered = append(filtered, port)
		}
	}
	if len(filtered) < 1 {
		return nil, fmt.Errorf("No ports were found. You might want to use a less specific query.")
	}
	if len(filtered) > 1 {
		return nil, fmt.Errorf("Multiple ports were found. Please use a more specific query.")
	}
	return filtered[0], nil
}
