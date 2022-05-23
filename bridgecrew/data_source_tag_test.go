package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTag(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTag(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_tag.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tag.test", "name"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tag.test", "description"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tag.test", "definition"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tag.test", "createdby"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tag.test", "creationdate"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tag.test", "isenabled"),
					//resource.TestCheckResourceAttrSet("data.bridgecrew_tag.test", "tagruleootbid"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tag.test", "repositories.0.name"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tag.test", "candoactions"),
				),
			},
		},
	})
}

func testAccDataSourceTag() string {
	return `
	data "bridgecrew_tag" "test" {
		id = "756a0450-f260-4275-9875-30ade35bea27"
	}`
}
