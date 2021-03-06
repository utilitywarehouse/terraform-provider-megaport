package megaport

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func resourceMegaportAwsVxc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMegaportAwsVxcCreate,
		ReadContext:   resourceMegaportAwsVxcRead,
		UpdateContext: resourceMegaportAwsVxcUpdate,
		DeleteContext: resourceMegaportAwsVxcDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Elem:     resourceMegaportVxcAwsEndElem(),
			},
			"invoice_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceMegaportVxcAwsEndElem() *schema.Resource {
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
			"aws_connection_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aws_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"aws_prefixes": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},
			"bgp_auth_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validateAwsBGPAuthKey,
			},
			"customer_asn": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"customer_ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"type": resourceAttributePrivatePublic(),
		},
	}
}

func flattenVxcEndAws(configProductUid string, v *api.ProductAssociatedVxc) []interface{} {
	var cc *api.ProductAssociatedVxcResourcesCspConnectionAws
	if cc_ := v.Resources.GetCspConnection(api.VxcConnectTypeAws); cc_ != nil {
		cc = cc_.(*api.ProductAssociatedVxcResourcesCspConnectionAws)
	}
	var prefixes []string
	if cc.Prefixes != "" {
		prefixes = strings.Split(cc.Prefixes, ",")
	}
	return []interface{}{map[string]interface{}{
		"product_uid":           configProductUid,
		"connected_product_uid": v.BEnd.ProductUid,
		"aws_connection_name":   cc.Name,
		"aws_account_id":        cc.OwnerAccount,
		"aws_ip_address":        cc.AmazonIpAddress,
		"aws_prefixes":          prefixes,
		"bgp_auth_key":          cc.AuthKey,
		"customer_asn":          int(cc.Asn),
		"customer_ip_address":   cc.CustomerIpAddress,
		"type":                  strings.ToLower(cc.Type),
	}}
}

func expandVxcEndAws(e map[string]interface{}) *api.PartnerConfigAws {
	pc := &api.PartnerConfigAws{
		AwsAccountId: api.String(e["aws_account_id"]),
		CustomerASN:  api.Uint64FromInt(e["customer_asn"]),
		Type:         api.String(e["type"]),
	}
	if v := e["aws_connection_name"]; v != "" {
		pc.AwsConnectionName = api.String(v)
	}
	if v := e["aws_ip_address"]; v != "" {
		pc.AmazonIPAddress = api.String(v)
	}
	if v := e["aws_prefixes"].(*schema.Set).List(); len(v) > 0 {
		pc.AmazonPrefixes = make([]string, len(v))
		for i, vv := range v {
			pc.AmazonPrefixes[i] = vv.(string)
		}
	}
	if v := e["bgp_auth_key"]; v != "" {
		pc.BGPAuthKey = api.String(v)
	}
	if v := e["customer_ip_address"]; v != "" {
		pc.CustomerIPAddress = api.String(v)
	}
	return pc
}

func resourceMegaportAwsVxcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}
	if err := d.Set("rate_limit", int(p.RateLimit)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("a_end", flattenVxcEnd(p.AEnd)); err != nil {
		return diag.FromErr(err)
	}
	puid := ""
	if v := d.Get("b_end").([]interface{}); len(v) > 0 {
		puid = v[0].(map[string]interface{})["product_uid"].(string)
	}
	if err := d.Set("b_end", flattenVxcEndAws(puid, p)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("invoice_reference", p.CostCentre); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceMegaportAwsVxcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	a := d.Get("a_end").([]interface{})[0].(map[string]interface{})
	b := d.Get("b_end").([]interface{})[0].(map[string]interface{})
	input := &api.CloudVxcCreateInput{
		ProductUidA:   api.String(a["product_uid"]),
		ProductUidB:   api.String(b["product_uid"]),
		Name:          api.String(d.Get("name")),
		PartnerConfig: expandVxcEndAws(b),
		RateLimit:     api.Uint64FromInt(d.Get("rate_limit")),
	}
	if v := b["aws_prefixes"].(*schema.Set).List(); len(v) > 0 && b["type"].(string) != "public" {
		return diag.FromErr(fmt.Errorf("cannot specify 'aws_prefixes' for a private VXC"))
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
			return diag.FromErr(err)
		}
		if !ok {
			return diag.FromErr(fmt.Errorf("VLAN id %d is unavailable on product %s", *input.VlanA, *input.ProductUidA))
		}
	}
	uid, err := cfg.Client.CreateCloudVxc(input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*uid)
	if err := waitUntilVxcIsConfigured(ctx, cfg.Client, *uid, 5*time.Minute); err != nil {
		return diag.FromErr(err)
	}
	return resourceMegaportAwsVxcRead(ctx, d, m)
}

func resourceMegaportAwsVxcUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	a := d.Get("a_end").([]interface{})[0].(map[string]interface{})
	b := d.Get("b_end").([]interface{})[0].(map[string]interface{})
	input := &api.CloudVxcUpdateInput{
		Name:          api.String(d.Get("name")),
		PartnerConfig: expandVxcEndAws(b),
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
			return diag.FromErr(err)
		}
		if !ok {
			return diag.FromErr(fmt.Errorf("VLAN id %d is unavailable on product %s", *input.VlanA, a["product_uid"].(string)))
		}
	}
	if err := cfg.Client.UpdateCloudVxc(input); err != nil {
		return diag.FromErr(err)
	}
	if err := waitUntilVxcIsConfigured(ctx, cfg.Client, d.Id(), 5*time.Minute); err != nil {
		return diag.FromErr(err)
	}
	if err := waitUntilAwsVxcIsUpdated(ctx, cfg.Client, input, 5*time.Minute); err != nil {
		return diag.FromErr(err)
	}
	return resourceMegaportAwsVxcRead(ctx, d, m)
}

func resourceMegaportAwsVxcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func waitUntilAwsVxcIsUpdated(ctx context.Context, client *api.Client, input *api.CloudVxcUpdateInput, timeout time.Duration) error {
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
			pc := input.PartnerConfig.(*api.PartnerConfigAws)
			var cc *api.ProductAssociatedVxcResourcesCspConnectionAws
			if cc_ := v.Resources.GetCspConnection(api.VxcConnectTypeAws); cc_ != nil {
				cc = cc_.(*api.ProductAssociatedVxcResourcesCspConnectionAws)
			}
			if !compareNillableStrings(pc.AmazonIPAddress, cc.AmazonIpAddress) {
				return nil, "", nil
			}
			if !compareNillableStrings(pc.AwsAccountId, cc.OwnerAccount) {
				return nil, "", nil
			}
			if !compareNillableStrings(pc.AwsConnectionName, cc.Name) {
				return nil, "", nil
			}
			if !compareNillableStrings(pc.BGPAuthKey, cc.AuthKey) {
				return nil, "", nil
			}
			if !compareNillableUints(pc.CustomerASN, cc.Asn) {
				return nil, "", nil
			}
			if !compareNillableStrings(pc.CustomerIPAddress, cc.CustomerIpAddress) {
				return nil, "", nil
			}
			if !compareNillableStrings(pc.Type, strings.ToLower(cc.Type)) {
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
