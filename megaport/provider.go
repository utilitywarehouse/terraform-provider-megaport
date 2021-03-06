package megaport

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/mutexkv"
)

var (
	megaportMutexKV = mutexkv.NewMutexKV()
)

type Config struct {
	Client *api.Client
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"MEGAPORT_TOKEN",
				}, nil),
			},
			"api_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"MEGAPORT_API_ENDPOINT",
				}, api.EndpointProduction),
				ValidateFunc: validation.IsURLWithHTTPS,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"megaport_port":        resourceMegaportPort(),
			"megaport_mcr":         resourceMegaportMcr(),
			"megaport_aws_vxc":     resourceMegaportAwsVxc(),
			"megaport_gcp_vxc":     resourceMegaportGcpVxc(),
			"megaport_private_vxc": resourceMegaportPrivateVxc(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"megaport_location":     dataSourceMegaportLocation(),
			"megaport_partner_port": dataSourceMegaportPartnerPort(),
			"megaport_port":         dataSourceMegaportPort(),
		},

		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			client := api.NewClient(d.Get("api_endpoint").(string))
			log.Printf("[INFO] Initialised megaport api client at %s", client.BaseURL)
			if v, ok := d.GetOk("token"); ok { // TODO: is it an error if not found?
				client.Token = v.(string)
			}
			return &Config{
				Client: client,
			}, nil
		},
	}
}
