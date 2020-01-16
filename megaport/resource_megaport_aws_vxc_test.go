package megaport

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func TestAccMegaportAwsVxc_basic(t *testing.T) {
	var vxcBefore api.ProductAssociatedVxc
	rName := "t" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rId := acctest.RandStringFromCharSet(12, "012346789")
	rAsn := uint64(acctest.RandIntRange(1, 65535))

	configValues := map[string]interface{}{
		"uid":            rName,
		"location":       "Equinix LD5",
		"aws_account_id": rId,
		"customer_asn":   rAsn,
	}
	cfg, err := testAccGetConfig("megaport_aws_vxc_basic", configValues, 0)
	if err != nil {
		t.Fatal(err)
	}
	cfgUpdate, err := testAccGetConfig("megaport_aws_vxc_basic_update", configValues, 1)
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
					testAccCheckResourceExists("megaport_port.foo", &vxcBefore),
					testAccCheckResourceExists("megaport_aws_vxc.foo", &vxcBefore),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "rate_limit", "100"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "invoice_reference", ""),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "a_end.0.vlan", "567"),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "b_end.0.product_uid", "data.megaport_partner_port.aws", "id"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_account_id", rId),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_asn", strconv.Itoa(int(rAsn))),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.type", "private"),
				),
			},
			{
				Config: cfgUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &vxcBefore),
					testAccCheckResourceExists("megaport_aws_vxc.foo", &vxcBefore),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "rate_limit", "1000"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "invoice_reference", rName),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "a_end.0.vlan", "568"),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "b_end.0.product_uid", "data.megaport_partner_port.aws", "id"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_account_id", rId),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_asn", strconv.Itoa(int(rAsn))),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.type", "private"),
				),
			},
		},
	})
}
