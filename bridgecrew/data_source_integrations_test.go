package bridgecrew

import (
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_integrations.test", "integrations.0.id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_integrations.test", "integrations.0.enable"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_integrations.test", "integrations.0.integration_details"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_integrations.test", "integrations.0.params"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_integrations.test", "integrations.0.type"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_integrations.test", "integrations.0.sf_execution_name"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_integrations.test", "integrations.0.status"),
				),
			},
		},
	})
}

func testAccDataSourceIntegrations() string {
	return `
	data "bridgecrew_integrations" "test" {
	}`
}
