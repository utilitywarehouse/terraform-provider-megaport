package megaport

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func resourceMegaportMcr() *schema.Resource {
	return &schema.Resource{
		Create: resourceMegaportMcrCreate,
		Read:   resourceMegaportMcrRead,
		Update: resourceMegaportMcrUpdate,
		Delete: resourceMegaportMcrDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"mcr_version": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2}),
			},
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
			"term": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 12, 24, 36}),
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

func resourceMegaportMcrRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	p, err := cfg.Client.GetMcr(d.Id())
	if err != nil {
		log.Printf("resourceMegaportMcrRead: %v", err)
		d.SetId("")
		return nil
	}
	if err := d.Set("mcr_version", p.McrVersion()); err != nil {
		return err
	}
	if err := d.Set("term", int(p.ContractTermMonths)); err != nil {
		return err
	}
	if err := d.Set("location_id", int(p.LocationId)); err != nil {
		return err
	}
	if err := d.Set("name", p.ProductName); err != nil {
		return err
	}
	if err := d.Set("rate_limit", int(p.PortSpeed)); err != nil {
		return err
	}
	if err := d.Set("asn", int(p.Resources.VirtualRouter.McrASN)); err != nil {
		return err
	}
	if err := d.Set("invoice_reference", p.CostCentre); err != nil {
		return err
	}
	return nil
}

func resourceMegaportMcrCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	var input api.McrCreateInput
	switch d.Get("mcr_version").(int) {
	case 1:
		input = &api.Mcr1CreateInput{
			LocationId:       api.Uint64FromInt(d.Get("location_id")),
			Name:             api.String(d.Get("name")),
			RateLimit:        api.Uint64FromInt(d.Get("rate_limit")),
			Asn:              api.Uint64FromInt(d.Get("asn")),
			Term:             api.Uint64FromInt(d.Get("term")),
			InvoiceReference: api.String(d.Get("invoice_reference")),
		}
	case 2:
		input = &api.Mcr2CreateInput{
			LocationId:       api.Uint64FromInt(d.Get("location_id")),
			Name:             api.String(d.Get("name")),
			RateLimit:        api.Uint64FromInt(d.Get("rate_limit")),
			Asn:              api.Uint64FromInt(d.Get("asn")),
			InvoiceReference: api.String(d.Get("invoice_reference")),
		}
	}
	uid, err := cfg.Client.CreateMcr(input)
	if err != nil {
		return err
	}
	d.SetId(*uid)
	if err := waitUntilMcrIsConfigured(cfg.Client, *uid, 5*time.Minute); err != nil {
		return err
	}
	return resourceMegaportMcrRead(d, m)
}

func resourceMegaportMcrUpdate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	var input api.McrUpdateInput
	switch d.Get("mcr_version").(int) {
	case 1:
		input = &api.Mcr1UpdateInput{
			InvoiceReference: api.String(d.Get("invoice_reference")),
			Name:             api.String(d.Get("name")),
			ProductUid:       api.String(d.Id()),
		}
	case 2:
		input = &api.Mcr2UpdateInput{
			InvoiceReference: api.String(d.Get("invoice_reference")),
			Name:             api.String(d.Get("name")),
			ProductUid:       api.String(d.Id()),
		}
	}
	if err := cfg.Client.UpdateMcr(input); err != nil {
		return err
	}
	if err := waitUntilMcrIsConfigured(cfg.Client, d.Id(), 5*time.Minute); err != nil {
		return err
	}
	return resourceMegaportMcrRead(d, m)
}

func resourceMegaportMcrDelete(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	err := cfg.Client.DeleteMcr(d.Id())
	if err != nil && err != api.ErrNotFound {
		return err
	}
	if err == api.ErrNotFound {
		log.Printf("resourceMegaportMcrDelete: resource not found, deleting anyway")
	}
	return nil
}

func waitUntilMcrIsConfigured(client *api.Client, productUid string, timeout time.Duration) error {
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
	_, err := scc.WaitForState()
	return err
}
