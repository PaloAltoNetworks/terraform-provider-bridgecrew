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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_repository_branches.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repository_branches.test", "branches.0.creationdate"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repository_branches.test", "branches.0.defaultbranch"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repository_branches.test", "branches.0.name"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repository_branches.test", "branches.0.defaultbranch"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repository_branches.test", "reponame"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repository_branches.test", "repoowner"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repository_branches.test", "source"),
				),
			},
		},
	})
}

func testAccDataSourceRepositoryBranches() string {
	return `
	data "bridgecrew_repository_branches" "test" {
	 reponame="terraform-aws-cassandra"
         repoowner="JamesWoolfenden"
	}`
}
