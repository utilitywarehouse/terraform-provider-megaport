package megaport

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

var (
	testAccConfigTemplates = &template.Template{}
)

func testAccCheckResourceExists(n string, o interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cfg := testAccProvider.Meta().(*Config)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckResourceExists: cannot find %q", n)
		}
		switch t := o.(type) {
		case *api.Product:
			v, err := cfg.Client.GetPort(rs.Primary.ID)
			if err != nil {
				return err
			}
			*(o.(*api.Product)) = *v
		case *api.ProductAssociatedVxc:
			v, err := cfg.Client.GetCloudVxc(rs.Primary.ID)
			if err != nil {
				return err
			}
			*(o.(*api.ProductAssociatedVxc)) = *v
		default:
			return fmt.Errorf("testAccCheckResourceExists: not implemented, cannot check %q of type %s", n, t)
		}
		return nil
	}
}

func testAccCheckResourceDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*Config)
	for n, rs := range s.RootModule().Resources {
		if strings.Split(n, ".")[0] == "data" {
			continue
		}
		switch rs.Type {
		case "megaport_port":
			v, err := cfg.Client.GetPort(rs.Primary.ID)
			if err != nil {
				return err
			}
			if v != nil && !isResourceDeleted(v.ProvisioningStatus) {
				return fmt.Errorf("testAccCheckResourceDestroy: %q (%s) has not been destroyed", n, rs.Primary.ID)
			}
		case "megaport_aws_vxc":
			v, err := cfg.Client.GetCloudVxc(rs.Primary.ID)
			if err != nil {
				return err
			}
			if v != nil && !isResourceDeleted(v.ProvisioningStatus) {
				return fmt.Errorf("testAccCheckResourceDestroy: %q (%s) has not been destroyed", n, rs.Primary.ID)
			}
		default:
			return fmt.Errorf("testAccCheckResourceDestroy: not implemented, cannot check %q (%s)", n, rs.Primary.ID)
		}
	}
	return nil
}

func testAccNewConfig(name string) (*template.Template, error) {
	config := ""
	if err := filepath.Walk(filepath.Join("../examples/", name), func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		r, err := filepath.Match("*.tf", f.Name())
		if err != nil {
			return err
		}
		if r {
			c, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			config = config + string(c)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	t, err := testAccConfigTemplates.New(name).Parse(config)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func testAccGetConfig(name string, values map[string]interface{}) (string, error) {
	var (
		t   *template.Template
		err error
		cfg = &strings.Builder{}
	)
	t = testAccConfigTemplates.Lookup(name)
	if t == nil {
		t, err = testAccNewConfig(name)
		if err != nil {
			return "", err
		}
	}
	if err := t.Execute(cfg, values); err != nil {
		return "", err
	}
	return cfg.String(), nil
}

func testAccLogConfig(cfg string) {
	l := strings.Split(cfg, "\n")
	for i := range l {
		l[i] = "      " + l[i]
	}
	fmt.Printf("+++ CONFIG:\n%s\n", strings.Join(l, "\n"))
}
