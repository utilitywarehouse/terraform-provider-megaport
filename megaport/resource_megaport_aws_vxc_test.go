package megaport

import (
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

	cfg, err := testAccGetConfig("megaport_aws_vxc_basic", map[string]interface{}{
		"uid":            rName,
		"location":       "Equinix LD5",
		"aws_account_id": rId,
		"customer_asn":   rAsn,
	})
	if err != nil {
		t.Fatal(err)
	}
	testAccLogConfig(cfg)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.port", &vxcBefore),
					testAccCheckResourceExists("megaport_aws_vxc.test", &vxcBefore),
				),
			},
		},
	})
}
