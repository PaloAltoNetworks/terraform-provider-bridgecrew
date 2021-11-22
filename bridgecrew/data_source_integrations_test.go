package bridgecrew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceIntegrations(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIntegrations(),
			},
		},
	})
}

func testAccDataSourceIntegrations() string {
	return fmt.Sprintf(
		`
data "bridgecrew_integrations" "test" {
}`)
}
