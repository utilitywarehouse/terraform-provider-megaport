package megaport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func init() {
	resource.AddTestSweepers("megaport_gcp_vxc", &resource.Sweeper{
		Name: "megaport_gcp_vxc",
		F:    testAccVxcSweeper(api.VxcTypeGcp),
	})
}

func TestAccMegaportGcpVxc_basic(t *testing.T) {
	var (
		vxc, vxcUpdated, vxcNew api.ProductAssociatedVxc
		port                    api.Product
	)
	rName := "t" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rpk, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	rpk = fmt.Sprintf("%s/europe-west1/%d", rpk, acctest.RandIntRange(1, 3))
	configValues := map[string]interface{}{
		"uid":        rName,
		"nameRegex":  "London",
		"pairingKey": rpk,
		"rateLimit":  "data.megaport_partner_port.gcp.bandwidths[0]",
		"vlan":       456,
	}
	cfg, err := newTestAccConfig("megaport_gcp_vxc_basic", configValues, 0)
	if err != nil {
		t.Fatal(err)
	}
	configValuesUpdate := mergeMaps(configValues, map[string]interface{}{
		"rateLimit": "data.megaport_partner_port.gcp.bandwidths[1]",
		"vlan":      567,
	})
	cfgUpdate, err := newTestAccConfig("megaport_gcp_vxc_full", configValuesUpdate, 1)
	if err != nil {
		t.Fatal(err)
	}
	configValuesForceNew := mergeMaps(configValuesUpdate, map[string]interface{}{
		"nameRegex": "Amsterdam",
	})
	cfgForceNew, err := newTestAccConfig("megaport_gcp_vxc_full", configValuesForceNew, 1)
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
					testAccCheckResourceExists("megaport_gcp_vxc.foo", &vxc),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttrPair("megaport_gcp_vxc.foo", "rate_limit", "data.megaport_partner_port.gcp", "bandwidths.0"),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "invoice_reference", ""),
					resource.TestCheckResourceAttrPair("megaport_gcp_vxc.foo", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "a_end.0.vlan", "456"),
					resource.TestCheckResourceAttrPair("megaport_gcp_vxc.foo", "b_end.0.product_uid", "data.megaport_partner_port.gcp", "id"),
					resource.TestCheckResourceAttrSet("megaport_gcp_vxc.foo", "b_end.0.connected_product_uid"),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "b_end.0.pairing_key", rpk),
				),
			},
			{
				ResourceName:            "megaport_gcp_vxc.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"b_end.0.product_uid"},
			},
			{
				PreConfig: func() { cfgUpdate.log() },
				Config:    cfgUpdate.Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &port),
					testAccCheckResourceExists("megaport_gcp_vxc.foo", &vxcUpdated),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttrPair("megaport_gcp_vxc.foo", "rate_limit", "data.megaport_partner_port.gcp", "bandwidths.1"),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "invoice_reference", rName),
					resource.TestCheckResourceAttrPair("megaport_gcp_vxc.foo", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "a_end.0.vlan", "567"),
					resource.TestCheckResourceAttrPair("megaport_gcp_vxc.foo", "b_end.0.product_uid", "data.megaport_partner_port.gcp", "id"),
					resource.TestCheckResourceAttrSet("megaport_gcp_vxc.foo", "b_end.0.connected_product_uid"),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "b_end.0.pairing_key", rpk),
				),
			},
			{
				ResourceName:            "megaport_gcp_vxc.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"b_end.0.product_uid"},
			},
			{
				PreConfig: func() { cfgForceNew.log() },
				Config:    cfgForceNew.Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &port),
					testAccCheckResourceExists("megaport_gcp_vxc.foo", &vxcNew),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttrPair("megaport_gcp_vxc.foo", "rate_limit", "data.megaport_partner_port.gcp", "bandwidths.1"),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "invoice_reference", rName),
					resource.TestCheckResourceAttrPair("megaport_gcp_vxc.foo", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "a_end.0.vlan", "567"),
					resource.TestCheckResourceAttrPair("megaport_gcp_vxc.foo", "b_end.0.product_uid", "data.megaport_partner_port.gcp", "id"),
					resource.TestCheckResourceAttrSet("megaport_gcp_vxc.foo", "b_end.0.connected_product_uid"),
					resource.TestCheckResourceAttr("megaport_gcp_vxc.foo", "b_end.0.pairing_key", rpk),
				),
			},
			{
				ResourceName:            "megaport_gcp_vxc.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"b_end.0.product_uid"},
			},
		},
	})

	if vxc.ProductUid != vxcUpdated.ProductUid {
		t.Errorf("TestAccMegaportGcpVxc_basic: expected the VXC to be updated but the resource ids differ")
	}
	if vxc.ProductUid == vxcNew.ProductUid {
		t.Errorf("TestAccMegaportGcpVxc_basic: expected the VXC to be recreated but the resource ids are identical")
	}
}
