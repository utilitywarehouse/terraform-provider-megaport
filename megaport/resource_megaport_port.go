package megaport

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem:       vxcResource,
				Set:        schema.HashResource(vxcResource),
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

var (
	// TODO: should these be functions?
	vxcResource = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rate_limit": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"a_end": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem:       vxcAEndResource,
			},
			"b_end": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem:       vxcBEndResource,
			},
			"invoice_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}

	vxcAEndResource = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vlan": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// TODO: untag?
			// TODO: product_uid might be needed for independant?
		},
	}

	vxcBEndResource = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vlan": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
)

func resourceMegaportPortRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	p, err := cfg.Client.Port.Get(d.Id())
	if err != nil {
		log.Printf("resourceMegaportPortRead: %v", err)
		d.SetId("")
		return nil
	}
	d.Set("location_id", p.LocationId)
	d.Set("name", p.ProductName)
	d.Set("speed", p.PortSpeed)
	d.Set("term", p.ContractTermMonths)
	vxcs := make([]interface{}, len(p.AssociatedVxcs))
	for i, v := range p.AssociatedVxcs {
		vxcs[i] = map[string]interface{}{
			"name":       v.ProductName,
			"rate_limit": int(v.RateLimit),
			"a_end": schema.NewSet(schema.HashResource(vxcAEndResource), []interface{}{map[string]interface{}{
				"vlan": int(v.AEnd.Vlan),
			}}),
			"b_end": schema.NewSet(schema.HashResource(vxcBEndResource), []interface{}{map[string]interface{}{
				"product_uid": v.BEnd.ProductUid,
				"vlan":        int(v.BEnd.Vlan),
			}}),
		}
	}
	d.Set("associated_vxcs", schema.NewSet(schema.HashResource(vxcResource), vxcs))
	d.Set("marketplace_visibility", p.MarketplaceVisibility)
	//d.Set("invoice_reference", p.) // TODO: is this even exported?
	return nil
}

func resourceMegaportPortCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	log.Printf("!!! CREATE")
	uid, err := cfg.Client.Port.Create(d.Get("name").(string),
		uint64(d.Get("location_id").(int)), uint64(d.Get("speed").(int)),
		uint64(d.Get("term").(int)), true)
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
	err := cfg.Client.Port.Delete(d.Id())
	if err != api.ErrNotFound {
		return err
	}
	log.Printf("resourceMegaportPortDelete: resource not found, deleting anyway")
	return nil
}

func resourceMegaportPortImportState(*schema.ResourceData, interface{}) ([]*schema.ResourceData, error) {
	log.Printf("!!! IMPORT")
	return nil, nil
}
