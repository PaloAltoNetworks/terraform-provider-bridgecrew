package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRepositories(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositories(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_repositories.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repositories.test", "repositories.0.repository"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repositories.test", "repositories.0.source"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repositories.test", "repositories.0.owner"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repositories.test", "repositories.0.defaultbranch"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repositories.test", "repositories.0.ispublic"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repositories.test", "repositories.0.runs"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repositories.test", "repositories.0.creationdate"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_repositories.test", "repositories.0.lastscandate"),
				),
			},
		},
	})
}

func testAccDataSourceRepositories() string {
	return `
	data "bridgecrew_repositories" "test" {
	}`
}
