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

func resourceMegaportMcr() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMegaportMcrCreate,
		ReadContext:   resourceMegaportMcrRead,
		UpdateContext: resourceMegaportMcrUpdate,
		DeleteContext: resourceMegaportMcrDelete,

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
			"rate_limit": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"asn": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"invoice_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceMegaportMcrRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	p, err := cfg.Client.GetMcr(d.Id())
	if err != nil {
		log.Printf("resourceMegaportMcrRead: %v", err)
		d.SetId("")
		return nil
	}
	if err := d.Set("location_id", int(p.LocationId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", p.ProductName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rate_limit", int(p.PortSpeed)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("asn", int(p.Resources.VirtualRouter.McrASN)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("invoice_reference", p.CostCentre); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceMegaportMcrCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	input := &api.Mcr2CreateInput{
		LocationId:       api.Uint64FromInt(d.Get("location_id")),
		Name:             api.String(d.Get("name")),
		RateLimit:        api.Uint64FromInt(d.Get("rate_limit")),
		Asn:              api.Uint64FromInt(d.Get("asn")),
		InvoiceReference: api.String(d.Get("invoice_reference")),
	}
	uid, err := cfg.Client.CreateMcr(input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*uid)
	if err := waitUntilMcrIsConfigured(ctx, cfg.Client, *uid, 5*time.Minute); err != nil {
		return diag.FromErr(err)
	}
	return resourceMegaportMcrRead(ctx, d, m)
}

func resourceMegaportMcrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	input := &api.Mcr2UpdateInput{
		InvoiceReference: api.String(d.Get("invoice_reference")),
		Name:             api.String(d.Get("name")),
		ProductUid:       api.String(d.Id()),
	}
	if err := cfg.Client.UpdateMcr(input); err != nil {
		return diag.FromErr(err)
	}
	if err := waitUntilMcrIsConfigured(ctx, cfg.Client, d.Id(), 5*time.Minute); err != nil {
		return diag.FromErr(err)
	}
	return resourceMegaportMcrRead(ctx, d, m)
}

func resourceMegaportMcrDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	err := cfg.Client.DeleteMcr(d.Id())
	if err != nil && err != api.ErrNotFound {
		return diag.FromErr(err)
	}
	if err == api.ErrNotFound {
		log.Printf("resourceMegaportMcrDelete: resource not found, deleting anyway")
	}
	return nil
}

func waitUntilMcrIsConfigured(ctx context.Context, client *api.Client, productUid string, timeout time.Duration) error {
	scc := &resource.StateChangeConf{
		Target: []string{api.ProductStatusConfigured, api.ProductStatusLive},
		Refresh: func() (interface{}, string, error) {
			v, err := client.GetMcr(productUid)
			if err != nil {
				log.Printf("[ERROR] Could not retrieve MCR while waiting for setup to finish: %v", err)
				return nil, "", err
			}
			if v == nil {
				return nil, "", nil
			}
			return v, v.ProvisioningStatus, nil
		},
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      5 * time.Second,
	}
	log.Printf("[INFO] Waiting for MCR (%s) to be configured", productUid)
	_, err := scc.WaitForStateContext(ctx)
	return err
}
