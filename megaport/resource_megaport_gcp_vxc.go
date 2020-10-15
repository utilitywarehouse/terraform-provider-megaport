package megaport

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func resourceMegaportGcpVxc() *schema.Resource {
	return &schema.Resource{
		Create: resourceMegaportGcpVxcCreate,
		Read:   resourceMegaportGcpVxcRead,
		Update: resourceMegaportGcpVxcUpdate,
		Delete: resourceMegaportGcpVxcDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Elem:     resourceMegaportVxcGcpEndElem(),
			},
			"invoice_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceMegaportVxcGcpEndElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "" && d.Get("b_end.0.connected_product_uid").(string) != ""
				},
			},
			"connected_product_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pairing_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func flattenVxcEndGcp(configProductUid string, v *api.ProductAssociatedVxc) []interface{} {
	var cc *api.ProductAssociatedVxcResourcesCspConnectionGcp
	if cc_ := v.Resources.GetCspConnection(api.VxcConnectTypeGoogle); cc_ != nil {
		cc = cc_.(*api.ProductAssociatedVxcResourcesCspConnectionGcp)
	}
	return []interface{}{map[string]interface{}{
		"product_uid":           configProductUid,
		"connected_product_uid": v.BEnd.ProductUid,
		"pairing_key":           cc.PairingKey,
	}}
}

func resourceMegaportGcpVxcRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	p, err := cfg.Client.GetVxc(d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not get VXC information: %v", err)
		d.SetId("")
		return nil
	}
	if p.ProvisioningStatus == api.ProductStatusDecommissioned {
		d.SetId("")
		return nil
	}
	if err := d.Set("name", p.ProductName); err != nil {
		return err
	}
	if err := d.Set("rate_limit", int(p.RateLimit)); err != nil {
		return err
	}
	if err := d.Set("a_end", flattenVxcEnd(p.AEnd)); err != nil {
		return err
	}
	puid := ""
	if v := d.Get("b_end").([]interface{}); len(v) > 0 {
		puid = v[0].(map[string]interface{})["product_uid"].(string)
	}
	if err := d.Set("b_end", flattenVxcEndGcp(puid, p)); err != nil {
		return err
	}
	if err := d.Set("invoice_reference", p.CostCentre); err != nil {
		return err
	}
	return nil
}

func resourceMegaportGcpVxcCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	a := d.Get("a_end").([]interface{})[0].(map[string]interface{})
	b := d.Get("b_end").([]interface{})[0].(map[string]interface{})
	input := &api.CloudVxcCreateInput{
		ProductUidA:   api.String(a["product_uid"]),
		ProductUidB:   api.String(b["product_uid"]),
		Name:          api.String(d.Get("name")),
		PartnerConfig: &api.PartnerConfigGcp{PairingKey: api.String(b["pairing_key"])},
		RateLimit:     api.Uint64FromInt(d.Get("rate_limit")),
	}
	if v, ok := d.GetOk("invoice_reference"); ok {
		input.InvoiceReference = api.String(v)
	}
	if v := a["vlan"].(int); v != 0 {
		input.VlanA = api.Uint64FromInt(v)
	}
	if *input.VlanA > 0 {
		ok, err := cfg.Client.GetPortVlanIdAvailable(*input.ProductUidA, *input.VlanA)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("VLAN id %d is unavailable on product %s", *input.VlanA, *input.ProductUidA)
		}
	}
	uid, err := cfg.Client.CreateCloudVxc(input)
	if err != nil {
		return err
	}
	d.SetId(*uid)
	if err := waitUntilVxcIsConfigured(cfg.Client, *uid, 5*time.Minute); err != nil {
		return err
	}
	return resourceMegaportGcpVxcRead(d, m)
}

func resourceMegaportGcpVxcUpdate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	a := d.Get("a_end").([]interface{})[0].(map[string]interface{})
	b := d.Get("b_end").([]interface{})[0].(map[string]interface{})
	input := &api.CloudVxcUpdateInput{
		Name:          api.String(d.Get("name")),
		PartnerConfig: &api.PartnerConfigGcp{PairingKey: api.String(b["pairing_key"])},
		ProductUid:    api.String(d.Id()),
		RateLimit:     api.Uint64FromInt(d.Get("rate_limit")),
	}
	if v, ok := d.GetOk("invoice_reference"); ok {
		input.InvoiceReference = api.String(v)
	}
	if v := a["vlan"].(int); v != 0 {
		input.VlanA = api.Uint64FromInt(v)
	}
	if *input.VlanA > 0 {
		ok, err := cfg.Client.GetPortVlanIdAvailable(a["product_uid"].(string), *input.VlanA)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("VLAN id %d is unavailable on product %s", *input.VlanA, a["product_uid"].(string))
		}
	}
	if err := cfg.Client.UpdateCloudVxc(input); err != nil {
		return err
	}
	if err := waitUntilVxcIsConfigured(cfg.Client, d.Id(), 5*time.Minute); err != nil {
		return err
	}
	if err := waitUntilGcpVxcIsUpdated(cfg.Client, input, 5*time.Minute); err != nil {
		return err
	}
	return resourceMegaportGcpVxcRead(d, m)
}

func resourceMegaportGcpVxcDelete(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	err := cfg.Client.DeleteVxc(d.Id())
	if err != nil && err != api.ErrNotFound {
		return err
	}
	if err == api.ErrNotFound {
		log.Printf("[DEBUG] VXC (%s) not found, deleting from state anyway", d.Id())
		return nil
	}
	if err := waitUntilVxcIsDeleted(cfg.Client, d.Id(), 5*time.Minute); err != nil {
		return err
	}
	return nil
}

func waitUntilGcpVxcIsUpdated(client *api.Client, input *api.CloudVxcUpdateInput, timeout time.Duration) error {
	scc := &resource.StateChangeConf{
		Target: []string{api.ProductStatusConfigured, api.ProductStatusLive},
		Refresh: func() (interface{}, string, error) {
			v, err := client.GetVxc(*input.ProductUid)
			if err != nil {
				log.Printf("[ERROR] Could not retrieve VXC while waiting for update to finish: %v", err)
				return nil, "", err
			}
			if v == nil {
				return nil, "", nil
			}
			if !compareNillableStrings(input.InvoiceReference, v.CostCentre) {
				return nil, "", nil
			}
			if !compareNillableStrings(input.Name, v.ProductName) {
				return nil, "", nil
			}
			if !compareNillableUints(input.RateLimit, v.RateLimit) {
				return nil, "", nil
			}
			if !compareNillableUints(input.VlanA, v.AEnd.Vlan) {
				return nil, "", nil
			}
			pc := input.PartnerConfig.(*api.PartnerConfigGcp)
			var cc *api.ProductAssociatedVxcResourcesCspConnectionGcp
			if cc_ := v.Resources.GetCspConnection(api.VxcConnectTypeGoogle); cc_ != nil {
				cc = cc_.(*api.ProductAssociatedVxcResourcesCspConnectionGcp)
			}
			if !compareNillableStrings(pc.PairingKey, cc.PairingKey) {
				return nil, "", nil
			}
			return v, v.ProvisioningStatus, nil
		},
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      5 * time.Second,
	}
	log.Printf("[INFO] Waiting for VXC (%s) to be updated", *input.ProductUid)
	_, err := scc.WaitForState()
	return err
}
