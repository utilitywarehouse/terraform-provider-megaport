package megaport

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func resourceMegaportPrivateVxc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMegaportPrivateVxcCreate,
		ReadContext:   resourceMegaportPrivateVxcRead,
		UpdateContext: resourceMegaportPrivateVxcUpdate,
		DeleteContext: resourceMegaportPrivateVxcDelete,

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
				Elem:     resourceMegaportVxcEndElem(),
			},
			"invoice_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceMegaportPrivateVxcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	p, err := cfg.Client.GetVxc(d.Id())
	if err != nil {
		log.Printf("resourceMegaportPrivateVxcRead: %v", err)
		d.SetId("")
		return nil
	}
	if p.ProvisioningStatus == api.ProductStatusDecommissioned {
		d.SetId("")
		return nil
	}
	if err := d.Set("name", p.ProductName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rate_limit", int(p.RateLimit)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("a_end", flattenVxcEnd(p.AEnd)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("b_end", flattenVxcEnd(p.BEnd)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("invoice_reference", p.CostCentre); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceMegaportPrivateVxcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	a := d.Get("a_end").([]interface{})[0].(map[string]interface{})
	b := d.Get("b_end").([]interface{})[0].(map[string]interface{})
	input := &api.PrivateVxcCreateInput{
		ProductUidA: api.String(a["product_uid"]),
		ProductUidB: api.String(b["product_uid"]),
		Name:        api.String(d.Get("name")),
		RateLimit:   api.Uint64FromInt(d.Get("rate_limit")),
	}
	if v, ok := d.GetOk("invoice_reference"); ok {
		input.InvoiceReference = api.String(v)
	}
	if v := a["vlan"].(int); v != 0 {
		input.VlanA = api.Uint64FromInt(v)
	}
	if v := b["vlan"].(int); v != 0 {
		input.VlanB = api.Uint64FromInt(v)
	}
	uid, err := cfg.Client.CreatePrivateVxc(input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*uid)
	if err := waitUntilVxcIsConfigured(ctx, cfg.Client, *uid, 5*time.Minute); err != nil {
		return diag.FromErr(err)
	}
	return resourceMegaportPrivateVxcRead(ctx, d, m)
}

func resourceMegaportPrivateVxcUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	a := d.Get("a_end").([]interface{})[0].(map[string]interface{})
	b := d.Get("b_end").([]interface{})[0].(map[string]interface{})
	input := &api.PrivateVxcUpdateInput{
		Name:       api.String(d.Get("name")),
		ProductUid: api.String(d.Id()),
		RateLimit:  api.Uint64FromInt(d.Get("rate_limit")),
	}
	if v, ok := d.GetOk("invoice_reference"); ok {
		input.InvoiceReference = api.String(v)
	}
	if v := a["vlan"].(int); v != 0 {
		input.VlanA = api.Uint64FromInt(v)
	}
	if v := b["vlan"].(int); v != 0 {
		input.VlanB = api.Uint64FromInt(v)
	}
	if err := cfg.Client.UpdatePrivateVxc(input); err != nil {
		return diag.FromErr(err)
	}
	if err := waitUntilVxcIsConfigured(ctx, cfg.Client, d.Id(), 5*time.Minute); err != nil {
		return diag.FromErr(err)
	}
	if err := waitUntilPrivateVxcIsUpdated(ctx, cfg.Client, input, 5*time.Minute); err != nil {
		return diag.FromErr(err)
	}
	return resourceMegaportPrivateVxcRead(ctx, d, m)
}

func resourceMegaportPrivateVxcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	err := cfg.Client.DeleteVxc(d.Id())
	if err != nil && err != api.ErrNotFound {
		return diag.FromErr(err)
	}
	if err == api.ErrNotFound {
		log.Printf("[DEBUG] VXC (%s) not found, deleting from state anyway", d.Id())
		return nil
	}
	if err := waitUntilVxcIsDeleted(ctx, cfg.Client, d.Id(), 5*time.Minute); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitUntilPrivateVxcIsUpdated(ctx context.Context, client *api.Client, input *api.PrivateVxcUpdateInput, timeout time.Duration) error {
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
			if !compareNillableUints(input.VlanB, v.BEnd.Vlan) {
				return nil, "", nil
			}
			return v, v.ProvisioningStatus, nil
		},
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      5 * time.Second,
	}
	log.Printf("[INFO] Waiting for VXC (%s) to be updated", *input.ProductUid)
	_, err := scc.WaitForStateContext(ctx)
	return err
}
