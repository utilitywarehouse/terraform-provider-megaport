package megaport

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

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
			v, err := cfg.Client.GetCloudVxc(rs.Primary.ID)
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
			v, err := cfg.Client.GetCloudVxc(rs.Primary.ID)
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
