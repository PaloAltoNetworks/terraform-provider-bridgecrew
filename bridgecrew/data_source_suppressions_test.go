package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSuppressions(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSuppressions(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_suppressions.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_suppressions.test", "suppressions.0.comment"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_suppressions.test", "suppressions.0.creationdate"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_suppressions.test", "suppressions.0.id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_suppressions.test", "suppressions.0.policyid"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_suppressions.test", "suppressions.0.suppressiontype"),
				),
			},
		},
	})
}

func testAccDataSourceSuppressions() string {
	return `
	data "bridgecrew_suppressions" "test" {
	}`
}
