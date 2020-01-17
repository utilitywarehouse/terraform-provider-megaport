package megaport

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func resourceAttributePrivatePublic() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "private",
		StateFunc: func(v interface{}) string {
			return strings.ToLower(v.(string))
		},
		ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
			vv := strings.ToLower(v.(string))
			if vv != "public" && vv != "private" {
				errs = append(errs, fmt.Errorf("%q must be either 'public' or 'private', got %s", k, vv))
			}
			return
		},
	}
}

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

func validateCIDRAddress(v interface{}, k string) (warns []string, errs []error) {
	vv, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	_, _, err := net.ParseCIDR(vv)
	if err != nil {
		errs = append(errs, fmt.Errorf("expected %q to be a valid IPv4 CIDR, got %v: %v", k, vv, err))
	}
	return
}

func flattenVxcEnd(v api.ProductAssociatedVxcEnd) []interface{} {
	return []interface{}{map[string]interface{}{
		"product_uid": v.ProductUid,
		"vlan":        int(v.Vlan),
	}}
}

func isResourceDeleted(provisioningStatus string) bool {
	switch provisioningStatus {
	case api.ProductStatusCancelled:
		fallthrough
	case api.ProductStatusCancelledParent:
		fallthrough
	case api.ProductStatusDecommissioned:
		return true
	default:
		return false
	}
}
