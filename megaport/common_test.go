package megaport

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func TestValidateAwsBGPAuthKey(t *testing.T) {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ012346789"
	short := acctest.RandStringFromCharSet(5, charset)
	long := acctest.RandStringFromCharSet(25, charset)
	whitespace := acctest.RandStringFromCharSet(5, charset) + " " + acctest.RandStringFromCharSet(5, charset)
	ok := acctest.RandStringFromCharSet(8, charset)
	e, w := validateAwsBGPAuthKey(ok, "test_value")
	if len(e) > 0 || len(w) > 0 {
		t.Errorf("validateAwsBGPAuthKey: %q failed validation, expected it to pass", ok)
	}
	e, w = validateAwsBGPAuthKey(short, "test_value")
	if len(e) == 0 && len(w) == 0 {
		t.Errorf("validateAwsBGPAuthKey: %q passed validation, expected it to fail", short)
	}
	e, w = validateAwsBGPAuthKey(long, "test_value")
	if len(e) == 0 && len(w) == 0 {
		t.Errorf("validateAwsBGPAuthKey: %q passed validation, expected it to fail", long)
	}
	e, w = validateAwsBGPAuthKey(whitespace, "test_value")
	if len(e) == 0 && len(w) == 0 {
		t.Errorf("validateAwsBGPAuthKey: %q passed validation, expected it to fail", whitespace)
	}
}

func testAccCheckResourceExists(n string, o interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cfg := testAccProvider.Meta().(*Config)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckResourceExists: cannot find %q", n)
		}
		switch t := o.(type) {
		case *api.Product:
			v, err := cfg.Client.GetPort(rs.Primary.ID)
			if err != nil {
				return err
			}
			*(o.(*api.Product)) = *v
		case *api.ProductAssociatedVxc:
			v, err := cfg.Client.GetVxc(rs.Primary.ID)
			if err != nil {
				return err
			}
			*(o.(*api.ProductAssociatedVxc)) = *v
		default:
			return fmt.Errorf("testAccCheckResourceExists: not implemented, cannot check %q of type %s", n, t)
		}
		return nil
	}
}

func testAccCheckResourceDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*Config)
	for n, rs := range s.RootModule().Resources {
		if strings.Split(n, ".")[0] == "data" {
			continue
		}
		switch rs.Type {
		case "megaport_port":
			v, err := cfg.Client.GetPort(rs.Primary.ID)
			if err != nil {
				return err
			}
			if v != nil && !isResourceDeleted(v.ProvisioningStatus) {
				return fmt.Errorf("testAccCheckResourceDestroy: %q (%s) has not been destroyed", n, rs.Primary.ID)
			}
		case "megaport_aws_vxc":
			fallthrough
		case "megaport_private_vxc":
			v, err := cfg.Client.GetVxc(rs.Primary.ID)
			if err != nil {
				return err
			}
			if v != nil && !isResourceDeleted(v.ProvisioningStatus) {
				return fmt.Errorf("testAccCheckResourceDestroy: %q (%s) has not been destroyed", n, rs.Primary.ID)
			}
		default:
			return fmt.Errorf("testAccCheckResourceDestroy: not implemented, cannot check %q (%s)", n, rs.Primary.ID)
		}
	}
	return nil
}

func testAccVxcSweeper(vxcType string) func(string) error {
	return func(region string) error {
		c, err := sharedClientForRegion(region)
		if err != nil {
			return fmt.Errorf("Error getting client: %s", err)
		}
		client := c.(*api.Client)
		ports, err := client.ListPorts()
		if err != nil {
			return err
		}
		for _, p := range ports {
			for _, v := range p.AssociatedVxcs {
				if strings.HasPrefix(v.ProductName, "terraform_acctest_") && !client.IsResourceDeleted(v.ProvisioningStatus) {
					vxc, err := client.GetVxc(v.ProductUid)
					if err != nil {
						return err
					}
					if vxc.Type() != vxcType {
						continue
					}
					if err := client.DeleteVxc(vxc.ProductUid); err != nil {
						log.Printf("[ERROR] Could not destroy VXC %q (%s) during sweep: %s", vxc.ProductName, vxc.ProductUid, err)
					}
				}
			}
		}
		return nil
	}
}
