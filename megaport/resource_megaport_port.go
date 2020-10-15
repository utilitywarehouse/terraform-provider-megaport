package megaport

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func resourceMegaportPort() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMegaportPortCreate,
		ReadContext:   resourceMegaportPortRead,
		UpdateContext: resourceMegaportPortUpdate,
		DeleteContext: resourceMegaportPortDelete,

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
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1000, 10000, 100000}),
			},
			"term": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 12, 24, 36}),
			},
			"invoice_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"marketplace_visibility": resourceAttributePrivatePublic(),
		},
	}
}

func resourceMegaportPortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	p, err := cfg.Client.GetPort(d.Id())
	if err != nil {
		log.Printf("resourceMegaportPortRead: %v", err)
		d.SetId("")
		return nil
	}
	if err := d.Set("location_id", int(p.LocationId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", p.ProductName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("speed", int(p.PortSpeed)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("term", int(p.ContractTermMonths)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("invoice_reference", p.CostCentre); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("marketplace_visibility", "private"); err != nil {
		return diag.FromErr(err)
	}
	if p.MarketplaceVisibility {
		if err := d.Set("marketplace_visibility", "public"); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceMegaportPortCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}
	d.SetId(*uid)
	if err := waitUntilPortIsConfigured(cfg.Client, *uid, 5*time.Minute); err != nil {
		return diag.FromErr(err)
	}
	return resourceMegaportPortRead(ctx, d, m)
}

func resourceMegaportPortUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	if err := cfg.Client.UpdatePort(&api.PortUpdateInput{
		InvoiceReference:      api.String(d.Get("invoice_reference")),
		Name:                  api.String(d.Get("name")),
		ProductUid:            api.String(d.Id()),
		MarketplaceVisibility: api.Bool(d.Get("marketplace_visibility") == "public"),
	}); err != nil {
		return diag.FromErr(err)
	}
	if err := waitUntilPortIsConfigured(cfg.Client, d.Id(), 5*time.Minute); err != nil {
		return diag.FromErr(err)
	}
	return resourceMegaportPortRead(ctx, d, m)
}

func resourceMegaportPortDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*Config)
	err := cfg.Client.DeletePort(d.Id())
	if err != nil && err != api.ErrNotFound {
		return diag.FromErr(err)
	}
	if err == api.ErrNotFound {
		log.Printf("resourceMegaportPortDelete: resource not found, deleting anyway")
	}
	return nil
}

func waitUntilPortIsConfigured(client *api.Client, productUid string, timeout time.Duration) error {
	scc := &resource.StateChangeConf{
		Target: []string{api.ProductStatusConfigured, api.ProductStatusLive},
		Refresh: func() (interface{}, string, error) {
			v, err := client.GetPort(productUid)
			if err != nil {
				log.Printf("[ERROR] Could not retrieve Port while waiting for setup to finish: %v", err)
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
	log.Printf("[INFO] Waiting for Port (%s) to be configured", productUid)
	_, err := scc.WaitForState()
	return err
}
