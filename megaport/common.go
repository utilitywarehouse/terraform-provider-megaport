package megaport

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

func validateAWSBGPAuthKey(v interface{}, k string) (warns []string, errs []error) {
	vv, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	if len(vv) < 6 || len(vv) > 24 {
		errs = append(errs, fmt.Errorf("%q must be between 6 and 24 characters long", k))
		return
	}
	if strings.Contains(vv, " ") {
		errs = append(errs, fmt.Errorf("%q cannot contain any whitespace", k))
		return
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

func compareNillableStrings(a *string, b string) bool {
	return a == nil || *a == b
}

func compareNillableUints(a *uint64, b uint64) bool {
	return a == nil || *a == b
}

func waitUntilVxcIsConfigured(client *api.Client, productUid string, timeout time.Duration) error {
	scc := &resource.StateChangeConf{
		Pending: []string{api.ProductStatusDeployable},
		Target:  []string{api.ProductStatusConfigured, api.ProductStatusLive},
		Refresh: func() (interface{}, string, error) {
			v, err := client.GetVxc(productUid)
			if err != nil {
				log.Printf("[ERROR] Could not retrieve VXC while waiting for setup to finish: %v", err)
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
	log.Printf("[INFO] Waiting for VXC (%s) to be configured", productUid)
	_, err := scc.WaitForState()
	return err
}

func waitUntilVxcIsDeleted(client *api.Client, productUid string, timeout time.Duration) error {
	initial, err := client.GetVxc(productUid)
	if err != nil {
		log.Printf("[ERROR] Could not retrieve VXC while waiting for deletion to finish: %v", err)
		return err
	}
	scc := &resource.StateChangeConf{
		Target: []string{api.ProductStatusDecommissioned},
		Refresh: func() (interface{}, string, error) {
			v, err := client.GetVxc(productUid)
			if err != nil {
				log.Printf("[ERROR] Could not retrieve VXC while waiting for deletion to finish: %v", err)
				return nil, "", err
			}
			if v == nil {
				return nil, "", nil
			}
			if initial.AEnd.Vlan > 0 {
				ok, err := client.GetPortVlanIdAvailable(initial.AEnd.ProductUid, initial.AEnd.Vlan)
				if err != nil {
					return v, "", err
				}
				if !ok {
					return v, "", nil
				}
			}
			if initial.BEnd.Vlan > 0 && initial.Type() == api.VXCTypePrivate {
				ok, err := client.GetPortVlanIdAvailable(initial.BEnd.ProductUid, initial.BEnd.Vlan)
				if err != nil {
					return v, "", err
				}
				if !ok {
					return v, "", nil
				}
			}
			return v, v.ProvisioningStatus, nil
		},
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      5 * time.Second,
	}
	log.Printf("[INFO] Waiting for VXC (%s) to be deleted", productUid)
	_, err = scc.WaitForState()
	return err
}
