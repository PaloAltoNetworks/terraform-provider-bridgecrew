package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMappings(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMappings(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_mappings.test", "guidelines.0.check"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_mappings.test", "guidelines.0.guideline"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_mappings.test", "idmappings.0.bcid"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_mappings.test", "idmappings.0.checkovid"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_mappings.test", "id"),
				),
			},
		},
	})
}

func testAccDataSourceMappings() string {
	return `
	data "bridgecrew_mappings" "test" {
	}`
}
