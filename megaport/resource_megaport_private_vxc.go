package megaport

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func resourceMegaportPrivateVxc() *schema.Resource {
	return &schema.Resource{
		Create: resourceMegaportPrivateVxcCreate,
		Read:   resourceMegaportPrivateVxcRead,
		Update: resourceMegaportPrivateVxcUpdate,
		Delete: resourceMegaportPrivateVxcDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMegaportPrivateVxcImportState,
		},

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
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     resourceMegaportVxcEndElem(),
			},
			"b_end": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     resourceMegaportVxcEndElem(),
			},
			"invoice_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceMegaportPrivateVxcRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	p, err := cfg.Client.GetPrivateVxc(d.Id())
	if err != nil {
		log.Printf("resourceMegaportPrivateVxcRead: %v", err)
		d.SetId("")
		return nil
	}
	if p.ProvisioningStatus == api.ProductStatusDecommissioned {
		d.SetId("")
		return nil
	}
	d.Set("name", p.ProductName)
	d.Set("rate_limit", p.RateLimit)
	d.Set("a_end", flattenVxcEnd(p.AEnd))
	d.Set("b_end", flattenVxcEnd(p.BEnd))
	d.Set("invoice_reference", p.CostCentre)
	return nil
}

func resourceMegaportPrivateVxcCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	a := d.Get("a_end").([]interface{})[0].(map[string]interface{})
	b := d.Get("b_end").([]interface{})[0].(map[string]interface{})
	uid, err := cfg.Client.CreatePrivateVxc(&api.PrivateVxcCreateInput{
		ProductUidA:      api.String(a["product_uid"]),
		ProductUidB:      api.String(b["product_uid"]),
		Name:             api.String(d.Get("name")),
		InvoiceReference: api.String(d.Get("invoice_reference")),
		VlanA:            api.Uint64FromInt(a["vlan"]),
		VlanB:            api.Uint64FromInt(b["vlan"]),
		RateLimit:        api.Uint64FromInt(d.Get("rate_limit")),
	})
	if err != nil {
		return err
	}
	d.SetId(*uid)
	return resourceMegaportPrivateVxcRead(d, m)
}

func resourceMegaportPrivateVxcUpdate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	a := d.Get("a_end").([]interface{})[0].(map[string]interface{})
	b := d.Get("b_end").([]interface{})[0].(map[string]interface{})
	var vlanB uint64
	if d.HasChange("b_end.0.vlan") {
		vlanB = uint64(b["vlan"].(int))
	}
	//if
	log.Printf(">>1 %#v", a)
	log.Printf(">>2 %#v", a["vlan"])
	if err := cfg.Client.UpdatePrivateVxc(&api.PrivateVxcUpdateInput{
		InvoiceReference: api.String(d.Get("invoice_reference")),
		Name:             api.String(d.Get("name")),
		ProductUid:       api.String(d.Id()),
		RateLimit:        api.Uint64FromInt(d.Get("rate_limit")),
		VlanA:            api.Uint64FromInt(a["vlan"]),
		VlanB:            api.Uint64FromInt(vlanB),
	}); err != nil {
		return err
	}
	return resourceMegaportPrivateVxcRead(d, m)
}

func resourceMegaportPrivateVxcDelete(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	err := cfg.Client.DeletePrivateVxc(d.Id())
	if err != nil && err != api.ErrNotFound {
		return err
	}
	if err == api.ErrNotFound {
		log.Printf("resourceMegaportPortDelete: resource not found, deleting anyway")
	}
	return nil
}

func resourceMegaportPrivateVxcImportState(*schema.ResourceData, interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}

func resourceMegaportVxcEndElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vlan": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func flattenVxcEnd(v api.ProductAssociatedVxcEnd) []interface{} {
	return []interface{}{map[string]interface{}{
		"product_uid": v.ProductUid,
		"vlan":        int(v.Vlan),
	}}
}
