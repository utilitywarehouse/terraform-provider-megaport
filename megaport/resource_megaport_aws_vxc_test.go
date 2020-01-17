package megaport

import (
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func TestAccMegaportAwsVxc_basic(t *testing.T) {
	var (
		vxc, vxcUpdated api.ProductAssociatedVxc
		port            api.Product
	)
	rName := "t" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rId := acctest.RandStringFromCharSet(12, "012346789")
	rAsn := uint64(acctest.RandIntRange(1, 65535))
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
		"aws_account_id":      rId,
		"customer_asn":        rAsn,
		"aws_ip_address":      ipA,
		"customer_ip_address": ipB,
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
					testAccCheckResourceExists("megaport_port.foo", &port),
					testAccCheckResourceExists("megaport_aws_vxc.foo", &vxc),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "name", "terraform_acctest_"+rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "rate_limit", "100"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "invoice_reference", ""),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "a_end.0.product_uid", "megaport_port.foo", "id"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "a_end.0.vlan", "567"),
					resource.TestCheckResourceAttrPair("megaport_aws_vxc.foo", "b_end.0.product_uid", "data.megaport_partner_port.aws", "id"),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.connected_product_uid"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_account_id", rId),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.aws_connection_name"),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.aws_ip_address"),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.bgp_auth_key"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_asn", strconv.Itoa(int(rAsn))),
					resource.TestCheckResourceAttrSet("megaport_aws_vxc.foo", "b_end.0.customer_ip_address"),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.type", "private"),
				),
			},
			{
				Config: cfgUpdate,
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
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_account_id", rId),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_connection_name", rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.aws_ip_address", ipA),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.bgp_auth_key", rName),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_asn", strconv.Itoa(int(rAsn))),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.customer_ip_address", ipB),
					resource.TestCheckResourceAttr("megaport_aws_vxc.foo", "b_end.0.type", "private"),
				),
			},
		},
	})
	if vxc.ProductUid != vxcUpdated.ProductUid {
		t.Errorf("TestAccMegaportAwsVxc_basic: expected the VXC to be updated but the resource ids differ")
	}
}
