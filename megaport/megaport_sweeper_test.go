package megaport

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sharedClientForRegion returns a common provider client configured for the
// specified region, in this case an environment
func sharedClientForRegion(region string) (interface{}, error) {
	endpoint := ""
	switch strings.ToLower(region) {
	case "staging":
		endpoint = api.EndpointStaging
	case "production":
		endpoint = api.EndpointProduction
	default:
		return nil, fmt.Errorf("Unknown region %q", region)
	}
	t := os.Getenv("MEGAPORT_TOKEN")
	if t == "" {
		return nil, fmt.Errorf("Must set the environment variable MEGAPORT_TOKEN")
	}
	client := api.NewClient(endpoint)
	client.Token = t
	return client, nil
}
