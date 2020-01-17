package megaport

import (
	"fmt"
	"log"

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
			State: schema.ImportStatePassthrough,
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
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
					vv := v.(int)
					if vv != 1000 && vv != 10000 && vv != 100000 {
						errs = append(errs, fmt.Errorf("%q must be 1000, 10000 or 100000 (Mbps), subject to availability, got %d", k, vv))
					}
					return
				},
			},
			"term": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
					vv := v.(int)
					if vv != 1 && vv != 12 && vv != 24 && vv != 36 {
						errs = append(errs, fmt.Errorf("%q must be 1, 12, 24 or 36, got %d", k, vv))
					}
					return
				},
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
			"marketplace_visibility": resourceAttributePrivatePublic(),
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
	if err := d.Set("location_id", p.LocationId); err != nil {
		return err
	}
	if err := d.Set("name", p.ProductName); err != nil {
		return err
	}
	if err := d.Set("speed", p.PortSpeed); err != nil {
		return err
	}
	if err := d.Set("term", p.ContractTermMonths); err != nil {
		return err
	}
	if err := d.Set("invoice_reference", p.CostCentre); err != nil {
		return err
	}
	if err := d.Set("associated_vxcs", schema.NewSet(schema.HashResource(resourceMegaportPrivateVxc()), flattenVxcList(p.AssociatedVxcs))); err != nil {
		return err
	}
	if err := d.Set("marketplace_visibility", "private"); err != nil {
		return err
	}
	if p.MarketplaceVisibility {
		if err := d.Set("marketplace_visibility", "public"); err != nil {
			return err
		}
	}
	return nil
}

func resourceMegaportPortCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	uid, err := cfg.Client.CreatePort(&api.PortCreateInput{
		LocationId:            api.Uint64FromInt(d.Get("location_id")),
		MarketplaceVisibility: api.Bool(d.Get("marketplace_visibility") == "public"),
		Name:                  api.String(d.Get("name")),
		Speed:                 api.Uint64FromInt(d.Get("speed")),
		Term:                  api.Uint64FromInt(d.Get("term")),
		InvoiceReference:      api.String(d.Get("invoice_reference")),
	})
	if err != nil {
		return err
	}
	d.SetId(*uid)
	return resourceMegaportPortRead(d, m)
}

func resourceMegaportPortUpdate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	if err := cfg.Client.UpdatePort(&api.PortUpdateInput{
		InvoiceReference:      api.String(d.Get("invoice_reference")),
		Name:                  api.String(d.Get("name")),
		ProductUid:            api.String(d.Id()),
		MarketplaceVisibility: api.Bool(d.Get("marketplace_visibility") == "public"),
	}); err != nil {
		return err
	}
	return resourceMegaportPortRead(d, m)
}

func resourceMegaportPortDelete(d *schema.ResourceData, m interface{}) error {
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
