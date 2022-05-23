package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOrganisation(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganisation(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_organisation.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_organisation.test", "organisation"),
				),
			},
		},
	})
}

func testAccDataSourceOrganisation() string {
	return `
	data "bridgecrew_organisation" "test" {
	}`
}
