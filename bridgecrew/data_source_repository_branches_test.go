package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRepositoryBranches(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryBranches(),
			},
		},
	})
}

func testAccDataSourceRepositoryBranches() string {
	return `
	data "bridgecrew_repository_branches" "test" {
	   target="cfngoat"
	}`
}
