package megaport

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func init() {
	resource.AddTestSweepers("megaport_mcr", &resource.Sweeper{
		Name: "megaport_mcr",
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
			mcrs, err := client.ListMcrs()
			if err != nil {
				return err
			}
			for _, m := range mcrs {
				if strings.HasPrefix(m.ProductName, "terraform_acctest_") && !client.IsResourceDeleted(m.ProvisioningStatus) {
					if err := client.DeleteMcr(m.ProductUid); err != nil {
						log.Printf("[ERROR] Could not destroy mcr %q (%s) during sweep: %s", m.ProductName, m.ProductUid, err)
					}
				}
			}
			return nil
		},
	})
}

func TestAccMegaportMcr2_basic(t *testing.T) {
	var mcr, mcrUpdated, mcrNew api.Product
	rName := "t" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	configValues := map[string]interface{}{
		"uid":        rName,
		"location":   "Global Switch London East",
		"rate_limit": 1000,
	}
	cfg, err := newTestAccConfig("megaport_mcr_basic", configValues, 0)
	if err != nil {
		t.Fatal(err)
	}
	cfgUpdate, err := newTestAccConfig("megaport_mcr_full", configValues, 1)
	if err != nil {
		t.Fatal(err)
	}
	rAsn := 1 + acctest.RandIntRange(0, math.MaxInt32)
	configValuesNew := mergeMaps(configValues, map[string]interface{}{
		"asn":        rAsn,
		"rate_limit": 2500,
	})
	cfgForceNew, err := newTestAccConfig("megaport_mcr_full", configValuesNew, 2)
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
					testAccCheckResourceExists("megaport_mcr.foo", &mcr),
					resource.TestCheckResourceAttr("megaport_mcr.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_mcr.foo", "rate_limit", "1000"),
					resource.TestCheckResourceAttrSet("megaport_mcr.foo", "asn"),
					resource.TestCheckResourceAttrPair("megaport_mcr.foo", "location_id", "data.megaport_location.foo", "id"),
					resource.TestCheckResourceAttr("megaport_mcr.foo", "invoice_reference", ""),
				),
			},
			{
				ResourceName:      "megaport_mcr.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() { cfgUpdate.log() },
				Config:    cfgUpdate.Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_mcr.foo", &mcrUpdated),
					resource.TestCheckResourceAttr("megaport_mcr.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_mcr.foo", "rate_limit", "1000"),
					resource.TestCheckResourceAttrSet("megaport_mcr.foo", "asn"),
					resource.TestCheckResourceAttrPair("megaport_mcr.foo", "location_id", "data.megaport_location.foo", "id"),
					resource.TestCheckResourceAttr("megaport_mcr.foo", "invoice_reference", rName),
				),
			},
			{
				ResourceName:      "megaport_mcr.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() { cfgForceNew.log() },
				Config:    cfgForceNew.Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_mcr.foo", &mcrNew),
					resource.TestCheckResourceAttr("megaport_mcr.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_mcr.foo", "rate_limit", "2500"),
					resource.TestCheckResourceAttr("megaport_mcr.foo", "asn", strconv.FormatUint(uint64(rAsn), 10)),
					resource.TestCheckResourceAttrPair("megaport_mcr.foo", "location_id", "data.megaport_location.foo", "id"),
					resource.TestCheckResourceAttr("megaport_mcr.foo", "invoice_reference", rName),
				),
			},
			{
				ResourceName:      "megaport_mcr.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	if mcr.ProductUid != mcrUpdated.ProductUid {
		t.Errorf("TestAccMegaportMcr_basic: expected the MCR to be updated but the resource ids differ")
	}
	if mcr.ProductUid == mcrNew.ProductUid {
		t.Errorf("TestAccMegaportMcr_basic: expected the MCR to be recreated but the resource ids are identical")
	}
}
