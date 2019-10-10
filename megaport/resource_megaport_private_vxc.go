package megaport

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

var (
	vxcEndResourceElem = &schema.Resource{
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
				Elem:     vxcEndResourceElem,
			},
			"b_end": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     vxcEndResourceElem,
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
	p, err := cfg.Client.Vxc.Get(d.Id())
	if err != nil {
		log.Printf("resourceMegaportPrivateVxcRead: %v", err)
		d.SetId("")
		return nil
	}
	d.Set("name", p.ProductName)
	d.Set("rate_limit", p.RateLimit)
	d.Set("a_end", []interface{}{map[string]interface{}{
		"product_uid": p.AEnd.ProductUid,
		"vlan":        int(p.AEnd.Vlan),
	}})
	d.Set("b_end", []interface{}{map[string]interface{}{
		"product_uid": p.BEnd.ProductUid,
		"vlan":        int(p.BEnd.Vlan),
	}})
	//d.Set("invoice_reference", p.) // TODO: is this even exported?
	return nil
}

func resourceMegaportPrivateVxcCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	a := d.Get("a_end").([]interface{})[0].(map[string]interface{})
	b := d.Get("b_end").([]interface{})[0].(map[string]interface{})
	uid, err := cfg.Client.Vxc.Create(
		a["product_uid"].(string),
		b["product_uid"].(string),
		d.Get("name").(string),
		uint64(a["vlan"].(int)),
		uint64(b["vlan"].(int)),
		uint64(d.Get("rate_limit").(int)),
	)
	if err != nil {
		return err
	}
	d.SetId(uid)
	return resourceMegaportPrivateVxcRead(d, m)
}

func resourceMegaportPrivateVxcUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceMegaportPrivateVxcDelete(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	err := cfg.Client.Vxc.Delete(d.Id())
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
