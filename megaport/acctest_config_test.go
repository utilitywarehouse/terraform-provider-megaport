package megaport

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	testAccConfigTemplates = &template.Template{}
)

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	r := make(map[string]interface{}, len(a))
	for k, v := range a {
		r[k] = v
	}
	for k, v := range b {
		r[k] = v
	}
	return r
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

type testAccConfig struct {
	Config string
	Name   string
	Step   int
}

func (c testAccConfig) log() {
	l := strings.Split(c.Config, "\n")
	for i := range l {
		l[i] = "      " + l[i]
	}
	fmt.Printf("+++ CONFIG %q (step %d):\n%s\n", c.Name, c.Step, strings.Join(l, "\n"))
}

func newTestAccConfig(name string, values map[string]interface{}, step int) (*testAccConfig, error) {
	var (
		t   *template.Template
		err error
		cfg = &strings.Builder{}
	)
	t = testAccConfigTemplates.Lookup(name)
	if t == nil {
		t, err = testAccNewConfig(name)
		if err != nil {
			return nil, err
		}
	}
	if err := t.Execute(cfg, values); err != nil {
		return nil, err
	}
	return &testAccConfig{Config: cfg.String(), Name: name, Step: step}, nil
}
