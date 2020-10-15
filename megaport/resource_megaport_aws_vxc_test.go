package megaport

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func init() {
	resource.AddTestSweepers("megaport_aws_vxc", &resource.Sweeper{
		Name: "megaport_aws_vxc",
		F:    testAccVxcSweeper(api.VxcTypeAws),
	})
}

func TestAccMegaportAwsVxc_basic(t *testing.T) {
	var (
		vxc, vxcUpdated, vxcNew api.ProductAssociatedVxc
		port                    api.Product
	)
	rName := "t" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rand.Seed(time.Now().UnixNano())
	n := &net.IPNet{
		IP:   net.IPv4(169, 254, byte(rand.Intn(255)), byte(rand.Intn(63)*4+1)),
		Mask: net.CIDRMask(30, 32),
	}
	ipA := n.String()
	n.IP[len(n.IP)-1]++
	ipB := n.String()
	configValues := map[string]interface{}{
		"uid":                 rName,
		"location":            "Equinix LD5",
		"aws_account_id":      acctest.RandStringFromCharSet(12, "012346789"),
		"customer_asn":        acctest.RandIntRange(1, 65536),
		"aws_ip_address":      ipA,
		"customer_ip_address": ipB,
		"type":                "private",
	}
	cfg, err := newTestAccConfig("megaport_aws_vxc_basic", configValues, 0)
	if err != nil {
		t.Fatal(err)
	}
	configValuesUpdate := mergeMaps(configValues, map[string]interface{}{
		"aws_account_id": acctest.RandStringFromCharSet(12, "012346789"),
		"customer_asn":   acctest.RandIntRange(1, 65536),
	})
	cfgUpdate, err := newTestAccConfig("megaport_aws_vxc_full", configValuesUpdate, 1)
	if err != nil {
		t.Fatal(err)
	}
	configValuesForceNew := mergeMaps(configValuesUpdate, map[string]interface{}{
		"location": "Interxion DUB2",
	})
	cfgForceNew, err := newTestAccConfig("megaport_aws_vxc_full", configValuesForceNew, 2)
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
					testAccCheckResourceExists("megaport_aws_vxc.foo", &vxc),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "rate_limit", "100"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "invoice_reference", ""),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "a_end.0.vlan", "567"),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "b_end.0.product_uid", "data.megaport_partner_port.aws", "id"),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.connected_product_uid"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_account_id", configValues["aws_account_id"].(string)),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.aws_connection_name"),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.aws_ip_address"),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.bgp_auth_key"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_asn", strconv.Itoa(configValues["customer_asn"].(int))),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.customer_ip_address"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.type", "private"),
				),
			},
			{
				ResourceName:            "megaport_aws_vxc.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"b_end.0.product_uid"},
			},
			{
				PreConfig: func() { cfgUpdate.log() },
				Config:    cfgUpdate.Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &port),
					testAccCheckResourceExists("megaport_aws_vxc.foo", &vxcUpdated),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "rate_limit", "1000"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "invoice_reference", rName),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "a_end.0.vlan", "568"),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "b_end.0.product_uid", "data.megaport_partner_port.aws", "id"),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.connected_product_uid"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_account_id", configValuesUpdate["aws_account_id"].(string)),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_connection_name", rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_ip_address", ipA),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.bgp_auth_key", rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_asn", strconv.Itoa(configValuesUpdate["customer_asn"].(int))),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_ip_address", ipB),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.type", "private"),
				),
			},
			{
				ResourceName:            "megaport_aws_vxc.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"b_end.0.product_uid"},
			},
			{
				PreConfig: func() { cfgForceNew.log() },
				Config:    cfgForceNew.Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("megaport_port.foo", &port),
					testAccCheckResourceExists("megaport_aws_vxc.foo", &vxcNew),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "rate_limit", "1000"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "invoice_reference", rName),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "a_end.0.vlan", "568"),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "b_end.0.product_uid", "data.megaport_partner_port.aws", "id"),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.connected_product_uid"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_account_id", configValuesForceNew["aws_account_id"].(string)),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_connection_name", rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_ip_address", ipA),
					resource.TestCheckNoResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_prefixes"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.bgp_auth_key", rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_asn", strconv.Itoa(configValuesForceNew["customer_asn"].(int))),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_ip_address", ipB),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.type", "private"),
				),
			},
			{
				ResourceName:            "megaport_aws_vxc.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"b_end.0.product_uid"},
			},
		},
	})

	if vxc.ProductUid != vxcUpdated.ProductUid {
		t.Errorf("TestAccMegaportAwsVxc_basic: expected the VXC to be updated but the resource ids differ")
	}
	if vxc.ProductUid == vxcNew.ProductUid {
		t.Errorf("TestAccMegaportAwsVxc_basic: expected the VXC to be recreated but the resource ids are identical")
	}
}

func TestAccMegaportAwsVxc_basicPublic(t *testing.T) {
	var (
		vxc  api.ProductAssociatedVxc
		port api.Product
	)
	rName := "t" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rand.Seed(time.Now().UnixNano())
	n := &net.IPNet{
		IP:   net.IPv4(88, byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(63)*4+1)),
		Mask: net.CIDRMask(30, 32),
	}
	ipA := n.String()
	n.IP[len(n.IP)-1]++
	ipB := n.String()
	n.Mask = net.CIDRMask(24, 32)
	n.IP[len(n.IP)-1] = 0
	prefix := n.String()
	configValues := map[string]interface{}{
		"uid":                 rName,
		"location":            "Equinix LD5",
		"aws_account_id":      acctest.RandStringFromCharSet(12, "012346789"),
		"customer_asn":        acctest.RandIntRange(1, 65536),
		"aws_ip_address":      ipA,
		"customer_ip_address": ipB,
		"prefixes":            []string{prefix},
		"type":                "public",
	}
	cfg, err := newTestAccConfig("megaport_aws_vxc_full", configValues, 0)
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
					testAccCheckResourceExists("megaport_aws_vxc.foo", &vxc),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "rate_limit", "1000"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "invoice_reference", rName),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "a_end.0.vlan", "568"),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "b_end.0.product_uid", "data.megaport_partner_port.aws", "id"),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.connected_product_uid"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_account_id", configValues["aws_account_id"].(string)),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_connection_name", rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_ip_address", ipA),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["megaport_aws_vxc.foo"]
						if !ok {
							return fmt.Errorf("megaport_aws_vxc.foo: Not found")
						}
						is := rs.Primary
						if is == nil {
							return fmt.Errorf("megaport_aws_vxc.foo: No primary instance")
						}
						p := []string{}
						for k, v := range is.Attributes {
							if k != "b_end.0.aws_prefixes.#" && strings.HasPrefix(k, "b_end.0.aws_prefixes.") {
								p = append(p, v)
							}
						}
						if diff := cmp.Diff(configValues["prefixes"], p); diff != "" {
							return fmt.Errorf("megaport_aws_vxc.foo: Attribute 'b_end.0.aws_prefix' unexpected value:\n%s", diff)
						}
						return nil
					},
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.bgp_auth_key", rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_asn", strconv.Itoa(configValues["customer_asn"].(int))),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_ip_address", ipB),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.type", "public"),
				),
			},
			{
				ResourceName:            "megaport_aws_vxc.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"b_end.0.product_uid"},
			},
		},
	})
}
