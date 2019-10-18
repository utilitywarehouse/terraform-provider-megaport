package megaport

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func resourceMegaportPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceMegaportPortCreate,
		Read:   resourceMegaportPortRead,
		Update: resourceMegaportPortUpdate,
		Delete: resourceMegaportPortDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMegaportPortImportState,
		},

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"speed": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"term": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"invoice_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"associated_vxcs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     resourceMegaportPrivateVxc(),
				Set:      schema.HashResource(resourceMegaportPrivateVxc()),
			},
			"marketplace_visibility": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "private",
				StateFunc: func(v interface{}) string {
					return strings.ToLower(v.(string))
				},
				ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
					vv := strings.ToLower(v.(string))
					if vv != "public" && vv != "private" {
						errs = append(errs, fmt.Errorf("%s must be either 'public' or 'private', got %s", k, vv))
					}
					return
				},
			},
			// TODO: LAG ports
		},
	}
}

func resourceMegaportPortRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	p, err := cfg.Client.GetPort(d.Id())
	if err != nil {
		log.Printf("resourceMegaportPortRead: %v", err)
		d.SetId("")
		return nil
	}
	d.Set("location_id", p.LocationId)
	d.Set("name", p.ProductName)
	d.Set("speed", p.PortSpeed)
	d.Set("term", p.ContractTermMonths)
	d.Set("associated_vxcs", schema.NewSet(schema.HashResource(resourceMegaportPrivateVxc()), flattenVxcList(p.AssociatedVxcs)))
	d.Set("marketplace_visibility", p.MarketplaceVisibility)
	//d.Set("invoice_reference", p.) // TODO: is this even exported?
	return nil
}

func resourceMegaportPortCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	log.Printf("!!! CREATE")
	uid, err := cfg.Client.CreatePort(d.Get("name").(string),
		uint64(d.Get("location_id").(int)), uint64(d.Get("speed").(int)),
		uint64(d.Get("term").(int)))
	if err != nil {
		return err
	}
	d.SetId(uid)
	return nil
}

func resourceMegaportPortUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("!!! UPDATE")
	return nil
}

func resourceMegaportPortDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("!!! DELETE")
	cfg := m.(*Config)
	err := cfg.Client.DeletePort(d.Id())
	if err != nil && err != api.ErrNotFound {
		return err
	}
	if err == api.ErrNotFound {
		log.Printf("resourceMegaportPortDelete: resource not found, deleting anyway")
	}
	return nil
}

func resourceMegaportPortImportState(*schema.ResourceData, interface{}) ([]*schema.ResourceData, error) {
	log.Printf("!!! IMPORT")
	return nil, nil
}

func flattenVxc(v api.ProductAssociatedVxc) interface{} {
	return map[string]interface{}{
		"name":       v.ProductName,
		"rate_limit": int(v.RateLimit),
		"a_end":      flattenVxcEnd(v.AEnd),
		"b_end":      flattenVxcEnd(v.BEnd),
	}
}

func flattenVxcList(vs []api.ProductAssociatedVxc) []interface{} {
	ret := make([]interface{}, len(vs))
	for i, v := range vs {
		ret[i] = flattenVxc(v)
	}
	return ret
}
