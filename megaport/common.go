package megaport

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

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
