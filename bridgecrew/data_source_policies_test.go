package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePolicies(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicies(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_policies.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_policies.test", "policies.0.cloud_provider"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_policies.test", "policies.0.title"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_policies.test", "policies.0.severity"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_policies.test", "policies.0.category"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_policies.test", "policies.0.resource_types.0"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_policies.test", "policies.0.guideline"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_policies.test", "policies.0.createdby"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_policies.test", "policies.0.frameworks.0"),
				),
			},
		},
	})
}

func testAccDataSourcePolicies() string {
	return `
	data "bridgecrew_policies" "test" {
	}`
}
