package megaport

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func init() {
	resource.AddTestSweepers("megaport_port", &resource.Sweeper{
		Name: "megaport_port",
		Dependencies: []string{
			"megaport_aws_vxc",
			"megaport_gcp_vxc",
			"megaport_private_vxc",
		},
		F: func(region string) error {
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
				if strings.HasPrefix(p.ProductName, "terraform_acctest_") && !client.IsResourceDeleted(p.ProvisioningStatus) {
					if err := client.DeletePort(p.ProductUid); err != nil {
						log.Printf("[ERROR] Could not destroy port %q (%s) during sweep: %s", p.ProductName, p.ProductUid, err)
					}
				}
			}
			return nil
		},
	})
}

func TestAccMegaportPort_basic(t *testing.T) {
	var port, portUpdated, portNew api.Product
	rName := "t" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	configValues := map[string]interface{}{
		"uid":      rName,
		"location": "Telehouse North",
	}
	cfg, err := newTestAccConfig("megaport_port_basic", configValues, 0)
	if err != nil {
		t.Fatal(err)
	}
	cfgUpdate, err := newTestAccConfig("megaport_port_basic_update", configValues, 1)
	if err != nil {
		t.Fatal(err)
	}
	cfgForceNew, err := newTestAccConfig("megaport_port_basic_forcenew", configValues, 2)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { cfg.log() },
				Config:    cfg.Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &port),
					resource.TestCheckResourceAttr("megaport_port.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_port.foo", "speed", "1000"),
					resource.TestCheckResourceAttr("megaport_port.foo", "term", "1"),
					resource.TestCheckResourceAttrPair("megaport_port.foo", "location_id", "data.megaport_location.foo", "id"),
					resource.TestCheckResourceAttr("megaport_port.foo", "invoice_reference", ""),
					resource.TestCheckResourceAttr("megaport_port.foo", "marketplace_visibility", "private"),
				),
			},
			{
				PreConfig: func() { cfgUpdate.log() },
				Config:    cfgUpdate.Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &portUpdated),
					resource.TestCheckResourceAttr("megaport_port.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_port.foo", "speed", "1000"),
					resource.TestCheckResourceAttr("megaport_port.foo", "term", "1"),
					resource.TestCheckResourceAttrPair("megaport_port.foo", "location_id", "data.megaport_location.foo", "id"),
					resource.TestCheckResourceAttr("megaport_port.foo", "invoice_reference", rName),
					resource.TestCheckResourceAttr("megaport_port.foo", "marketplace_visibility", "public"),
				),
			},
			{
				PreConfig: func() { cfgForceNew.log() },
				Config:    cfgForceNew.Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &portNew),
					resource.TestCheckResourceAttr("megaport_port.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_port.foo", "speed", "10000"),
					resource.TestCheckResourceAttr("megaport_port.foo", "term", "12"),
					resource.TestCheckResourceAttrPair("megaport_port.foo", "location_id", "data.megaport_location.foo", "id"),
					resource.TestCheckResourceAttr("megaport_port.foo", "invoice_reference", rName),
					resource.TestCheckResourceAttr("megaport_port.foo", "marketplace_visibility", "public"),
				),
			},
		},
	})

	if port.ProductUid != portUpdated.ProductUid {
		t.Errorf("TestAccMegaportPort_basic: expected the port to be updated but the resource ids differ")
	}
	if port.ProductUid == portNew.ProductUid {
		t.Errorf("TestAccMegaportPort_basic: expected the port to be recreated but the resource ids are identical")
	}
}
