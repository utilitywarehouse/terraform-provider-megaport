package megaport

import (
	"github.com/hashicorp/terraform/helper/schema"
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
			"product_uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vlan": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"untag": { // TODO: what's the key for the API request?
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"b_end": {
				Type:       schema.TypeSet,
				Required:   true,
				MaxItems:   1,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_uid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vlan": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"invoice_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceMegaportPrivateVxcRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceMegaportPrivateVxcCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	var a, b map[string]interface{}
	//if t := d.Get("a_end"); t != nil {
	//a = t.(*schema.Set).List()[0].(*schema.ResourceData)
	//}
	if t := d.Get("b_end").(*schema.Set).List(); t != nil {
		//b = (t.(*schema.Set)).List()[0].(map[string]interface)
		b = t[0].(map[string]interface{})
	}
	var vlanA uint64
	if a != nil {
		vlanA = uint64(a["vlan"].(int))
	}
	uid, err := cfg.Client.Vxc.Create(
		d.Get("product_uid").(string),
		b["product_uid"].(string),
		d.Get("name").(string),
		vlanA,
		uint64(b["vlan"].(int)),
		uint64(d.Get("rate_limit").(int)),
	)
	if err != nil {
		return err
	}
	d.SetId(uid)
	return nil
}

func resourceMegaportPrivateVxcUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceMegaportPrivateVxcDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceMegaportPrivateVxcImportState(*schema.ResourceData, interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}
