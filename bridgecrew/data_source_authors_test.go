package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAuthors(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAuthors(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_authors.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_authors.test", "authors.0"),
				),
			},
		},
	})
}

func testAccDataSourceAuthors() string {
	return `
	data "bridgecrew_authors" "test" {
		fullreponame="JamesWoolfenden/test"
		sourcetype="Github"
	}`
}
