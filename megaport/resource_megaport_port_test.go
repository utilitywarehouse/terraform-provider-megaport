package megaport

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func TestAccMegaportPort_basic(t *testing.T) {
	var portBefore api.Product
	rName := "t" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	cfg, err := testAccGetConfig("megaport_port_basic", map[string]interface{}{
		"uid":      rName,
		"location": "Telehouse North",
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
					testAccCheckResourceExists("megaport_port.test", &portBefore),
				),
			},
		},
	})
}
