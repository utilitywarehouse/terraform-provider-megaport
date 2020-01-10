package megaport

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func TestAccMegaportPort_basic(t *testing.T) {
	var portBefore api.Product
	rName := "t" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPortBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.test", &portBefore),
				),
			},
		},
	})
}

func testAccPortBasicConfig(name string) string {
	return fmt.Sprintf(`
data "megaport_location" "port" {
  name_regex = "Telehouse North"
}

resource "megaport_port" "test" {
  name        = "terraform_acctest_%s"
  location_id = data.megaport_location.port.id
  speed       = 1000
  term        = 1
}
`, name)
}
