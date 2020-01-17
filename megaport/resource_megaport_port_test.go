package megaport

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func TestAccMegaportPort_basic(t *testing.T) {
	var port, portUpdated, portNew api.Product
	rName := "t" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	configValues := map[string]interface{}{
		"uid":      rName,
		"location": "Telehouse North",
	}
	cfg, err := testAccGetConfig("megaport_port_basic", configValues, 0)
	if err != nil {
		t.Fatal(err)
	}
	cfgUpdate, err := testAccGetConfig("megaport_port_basic_update", configValues, 1)
	if err != nil {
		t.Fatal(err)
	}
	cfgForceNew, err := testAccGetConfig("megaport_port_basic_forcenew", configValues, 2)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &port),
					resource.TestCheckResourceAttr("megaport_port.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_port.foo", "speed", "1000"),
					resource.TestCheckResourceAttr("megaport_port.foo", "term", "1"),
					resource.TestCheckResourceAttrPair("megaport_port.foo", "location_id", "data.megaport_location.foo", "id"),
					resource.TestCheckResourceAttr("megaport_port.foo", "invoice_reference", ""),
					resource.TestCheckNoResourceAttr("megaport_port.foo", "associated_vxcs"),
					resource.TestCheckResourceAttr("megaport_port.foo", "marketplace_visibility", "private"),
				),
			},
			{
				Config: cfgUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &portUpdated),
					resource.TestCheckResourceAttr("megaport_port.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_port.foo", "speed", "1000"),
					resource.TestCheckResourceAttr("megaport_port.foo", "term", "1"),
					resource.TestCheckResourceAttrPair("megaport_port.foo", "location_id", "data.megaport_location.foo", "id"),
					resource.TestCheckResourceAttr("megaport_port.foo", "invoice_reference", rName),
					resource.TestCheckNoResourceAttr("megaport_port.foo", "associated_vxcs"),
					resource.TestCheckResourceAttr("megaport_port.foo", "marketplace_visibility", "public"),
				),
			},
			{
				Config: cfgForceNew,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &portNew),
					resource.TestCheckResourceAttr("megaport_port.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_port.foo", "speed", "10000"),
					resource.TestCheckResourceAttr("megaport_port.foo", "term", "12"),
					resource.TestCheckResourceAttrPair("megaport_port.foo", "location_id", "data.megaport_location.foo", "id"),
					resource.TestCheckResourceAttr("megaport_port.foo", "invoice_reference", rName),
					resource.TestCheckNoResourceAttr("megaport_port.foo", "associated_vxcs"),
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
