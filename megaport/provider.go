package megaport

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

type Config struct {
	Client *api.Client
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"MEGAPORT_TOKEN",
				}, nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{},

		DataSourcesMap: map[string]*schema.Resource{
			"megaport_locations":          dataSourceLocations(),
			"megaport_megaports":          dataSourceMegaports(),
			"megaport_internet_exchanges": dataSourceInternetExchanges(),
		},

		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			client := api.NewClient(api.EndpointProduction)
			if v, ok := d.GetOk("token"); ok {
				client.SetToken(v.(string))
			}
			return &Config{
				Client: client,
			}, nil
		},
	}
}
