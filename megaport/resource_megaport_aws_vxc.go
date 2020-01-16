package megaport

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func resourceMegaportAwsVxc() *schema.Resource {
	return &schema.Resource{
		Create: resourceMegaportAwsVxcCreate,
		Read:   resourceMegaportAwsVxcRead,
		Update: resourceMegaportAwsVxcUpdate,
		Delete: resourceMegaportAwsVxcDelete,

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
				ForceNew: true,
			},
			"aws_ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCIDRAddress,
			},
			"bgp_auth_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"customer_asn": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"customer_ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCIDRAddress,
			},
			"type": resourceAttributePrivatePublic(),
		},
	}
}

func flattenVxcEndAws(configProductUid string, v api.ProductAssociatedVxcEnd, r api.ProductAssociatedVxcResources) []interface{} {
	return []interface{}{map[string]interface{}{
		"product_uid":           configProductUid,
		"connected_product_uid": v.ProductUid,
		"aws_connection_name":   r.AwsVirtualInterface.Name,
		"aws_account_id":        r.AwsVirtualInterface.OwnerAccount,
		"aws_ip_address":        r.AwsVirtualInterface.AmazonIpAddress,
		"bgp_auth_key":          r.AwsVirtualInterface.AuthKey,
		"customer_asn":          int(r.AwsVirtualInterface.Asn),
		"customer_ip_address":   r.AwsVirtualInterface.CustomerIpAddress,
		"type":                  strings.ToLower(r.AwsVirtualInterface.Type),
	}}
}

func resourceMegaportAwsVxcRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	p, err := cfg.Client.GetCloudVxc(d.Id())
	if err != nil {
		log.Printf("resourceMegaportAwsVxcRead: %v", err)
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
	if err := d.Set("rate_limit", p.RateLimit); err != nil {
		return err
	}
	if err := d.Set("a_end", flattenVxcEnd(p.AEnd)); err != nil {
		return err
	}
	puid := d.Get("b_end").([]interface{})[0].(map[string]interface{})["product_uid"].(string)
	if err := d.Set("b_end", flattenVxcEndAws(puid, p.BEnd, p.Resources)); err != nil {
		return err
	}
	if err := d.Set("invoice_reference", p.CostCentre); err != nil {
		return err
	}
	return nil
}

func resourceMegaportAwsVxcCreate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	a := d.Get("a_end").([]interface{})[0].(map[string]interface{})
	b := d.Get("b_end").([]interface{})[0].(map[string]interface{})
	input := &api.CloudVxcCreateInput{
		ProductUidA: api.String(a["product_uid"]),
		ProductUidB: api.String(b["product_uid"]),
		Name:        api.String(d.Get("name")),
		RateLimit:   api.Uint64FromInt(d.Get("rate_limit")),
	}
	if v, ok := d.GetOk("invoice_reference"); ok {
		input.InvoiceReference = api.String(v)
	}
	if v := a["vlan"]; v != 0 {
		input.VlanA = api.Uint64FromInt(a["vlan"])
	}
	inputPartnerConfig := &api.PartnerConfigAWS{
		AWSAccountID: api.String(b["aws_account_id"]),
		CustomerASN:  api.Uint64FromInt(b["customer_asn"]),
		Type:         api.String(b["type"]),
	}
	if v := b["aws_connection_name"]; v != "" {
		inputPartnerConfig.AWSConnectionName = api.String(v)
	}
	if v := b["amazon_ip_address"]; v != "" {
		inputPartnerConfig.AmazonIPAddress = api.String(v)
	}
	if v := b["bgp_auth_key"]; v != "" {
		inputPartnerConfig.BGPAuthKey = api.String(v)
	}
	if v := b["customer_ip_address"]; v != "" {
		inputPartnerConfig.CustomerIPAddress = api.String(v)
	}
	input.PartnerConfig = inputPartnerConfig
	uid, err := cfg.Client.CreateCloudVxc(input)
	if err != nil {
		return err
	}
	d.SetId(*uid)
	if err := waitUntilVxcIsConfigured(cfg.Client, *uid, 5*time.Minute); err != nil {
		return err
	}
	return resourceMegaportAwsVxcRead(d, m)
}

func resourceMegaportAwsVxcUpdate(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	a := d.Get("a_end").([]interface{})[0].(map[string]interface{})
	//b := d.Get("b_end").([]interface{})[0].(map[string]interface{}) // TODO
	input := &api.CloudVxcUpdateInput{
		InvoiceReference: api.String(d.Get("invoice_reference")),
		Name:             api.String(d.Get("name")),
		ProductUid:       api.String(d.Id()),
		RateLimit:        api.Uint64FromInt(d.Get("rate_limit")),
		VlanA:            api.Uint64FromInt(a["vlan"]),
	}
	if err := cfg.Client.UpdateCloudVxc(input); err != nil {
		return err
	}
	if err := waitUntilVxcIsConfigured(cfg.Client, d.Id(), 5*time.Minute); err != nil {
		return err
	}
	if err := waitUntilVxcIsUpdated(cfg.Client, input, 5*time.Minute); err != nil {
		return err
	}
	return resourceMegaportAwsVxcRead(d, m)
}

func resourceMegaportAwsVxcDelete(d *schema.ResourceData, m interface{}) error {
	cfg := m.(*Config)
	err := cfg.Client.DeleteCloudVxc(d.Id())
	if err != nil && err != api.ErrNotFound {
		return err
	}
	if err == api.ErrNotFound {
		log.Printf("resourceMegaportPortDelete: resource not found, deleting anyway")
	}
	return nil
}

func waitUntilVxcIsConfigured(client *api.Client, productUid string, timeout time.Duration) error {
	scc := &resource.StateChangeConf{
		Pending:    []string{api.ProductStatusDeployable},
		Target:     []string{api.ProductStatusConfigured, api.ProductStatusLive},
		Refresh:    resourceMegaportAwsVxcStateRefreshFunc(client, productUid),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      5 * time.Second,
	}
	log.Printf("[INFO] Waiting for VXC (%s) to be available", productUid)
	_, err := scc.WaitForState()
	return err
}

func resourceMegaportAwsVxcStateRefreshFunc(client *api.Client, uid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.GetCloudVxc(uid)
		if err != nil {
			log.Printf("Error retrieving VXC while waiting for setup to finish: %v", err)
			return nil, "", err
		}
		if v == nil {
			return nil, "", nil
		}
		return v, v.ProvisioningStatus, nil
	}
}

func waitUntilVxcIsUpdated(client *api.Client, input *api.CloudVxcUpdateInput, timeout time.Duration) error {
	scc := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"UPDATED"},
		Refresh: func() (interface{}, string, error) {
			v, err := client.GetCloudVxc(*input.ProductUid)
			if err != nil {
				log.Printf("Error retrieving VXC while waiting for setup to finish: %v", err)
				return nil, "", err
			}
			if v == nil {
				return nil, "", nil
			}
			if v.CostCentre != *input.InvoiceReference {
				return nil, "PENDING", nil
			}
			if v.ProductName != *input.Name {
				return nil, "PENDING", nil
			}
			if v.RateLimit != *input.RateLimit {
				return nil, "PENDING", nil
			}
			if v.AEnd.Vlan != *input.VlanA {
				return nil, "PENDING", nil
			}
			return v, "UPDATED", nil
		},
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      5 * time.Second,
	}
	log.Printf("[INFO] Waiting for VXC (%s) to be available", *input.ProductUid)
	_, err := scc.WaitForState()
	return err
}
