package megaport

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func init() {
	resource.AddTestSweepers("megaport_private_vxc", &resource.Sweeper{
		Name: "megaport_private_vxc",
		F:    testAccVxcSweeper(api.VxcTypePrivate),
	})
}

func TestAccMegaportPrivateVxc_basic(t *testing.T) {
	var (
		vxc, vxcUpdated, vxcNew api.ProductAssociatedVxc
		portA, portB            api.Product
	)
	rName := "t" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	configValues := map[string]interface{}{
		"uid":       rName,
		"locationA": "Equinix LD5",
		"locationB": "Global Switch London East",
	}
	cfg, err := newTestAccConfig("megaport_private_vxc_basic", configValues, 0)
	if err != nil {
		t.Fatal(err)
	}
	configValuesUpdate := mergeMaps(configValues, map[string]interface{}{
		"vlanA": 567,
		"vlanB": 567,
	})
	cfgUpdate, err := newTestAccConfig("megaport_private_vxc_full", configValuesUpdate, 1)
	if err != nil {
		t.Fatal(err)
	}
	configValuesForceNew := mergeMaps(configValues, map[string]interface{}{
		"locationB": "Telehouse North$",
		"vlanA":     456,
		"vlanB":     456,
	})
	cfgForceNew, err := newTestAccConfig("megaport_private_vxc_full", configValuesForceNew, 2)
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
				Destroy:   false,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &portA),
					testAccCheckResourceExists("megaport_port.bar", &portB),
					testAccCheckResourceExists("megaport_private_vxc.foobar", &vxc),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "rate_limit", "100"),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "invoice_reference", ""),
					resource.TestCheckResourceAttrPair("megaport_private_vxc.foobar", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttrSet("megaport_private_vxc.foobar", "a_end.0.vlan"),
					resource.TestCheckResourceAttrPair("megaport_private_vxc.foobar", "b_end.0.product_uid", "megaport_port.bar", "id"),
					resource.TestCheckResourceAttrSet("megaport_private_vxc.foobar", "b_end.0.vlan"),
				),
			},
			{
				PreConfig: func() { cfgUpdate.log() },
				Config:    cfgUpdate.Config,
				Destroy:   false,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &portA),
					testAccCheckResourceExists("megaport_port.bar", &portB),
					testAccCheckResourceExists("megaport_private_vxc.foobar", &vxcUpdated),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "rate_limit", "200"),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "invoice_reference", rName),
					resource.TestCheckResourceAttrPair("megaport_private_vxc.foobar", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "a_end.0.vlan", strconv.Itoa(configValuesUpdate["vlanA"].(int))),
					resource.TestCheckResourceAttrPair("megaport_private_vxc.foobar", "b_end.0.product_uid", "megaport_port.bar", "id"),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "b_end.0.vlan", strconv.Itoa(configValuesUpdate["vlanB"].(int))),
				),
			},
			{
				PreConfig: func() { cfgForceNew.log() },
				Config:    cfgForceNew.Config,
				Destroy:   false,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &portA),
					testAccCheckResourceExists("megaport_port.bar", &portB),
					testAccCheckResourceExists("megaport_private_vxc.foobar", &vxcNew),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "rate_limit", "200"),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "invoice_reference", rName),
					resource.TestCheckResourceAttrPair("megaport_private_vxc.foobar", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "a_end.0.vlan", strconv.Itoa(configValuesForceNew["vlanA"].(int))),
					resource.TestCheckResourceAttrPair("megaport_private_vxc.foobar", "b_end.0.product_uid", "megaport_port.bar", "id"),
					resource.TestCheckResourceAttr("megaport_private_vxc.foobar", "b_end.0.vlan", strconv.Itoa(configValuesForceNew["vlanB"].(int))),
				),
			},
		},
	})

	if vxc.ProductUid != vxcUpdated.ProductUid {
		t.Errorf("TestAccMegaportPrivateVxc_basic: expected the VXC to be updated but the resource ids differ")
	}
	if vxc.ProductUid == vxcNew.ProductUid {
		t.Errorf("TestAccMegaportPrivateVxc_basic: expected the VXC to be recreated but the resource ids are identical")
	}
}
