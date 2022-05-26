package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTags(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTags(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_tags.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tags.test", "tags.0.name"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tags.test", "tags.0.description"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tags.test", "tags.0.creationdate"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tags.test", "tags.0.isenabled"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tags.test", "tags.0.candoactions"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tags.test", "tags.0.tagruleootbid"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_tags.test", "tags.0.id"),

					resource.TestCheckNoResourceAttr("data.bridgecrew_tags.test", "tags.0.repositories"),
					resource.TestCheckResourceAttr("data.bridgecrew_tags.test", "tags.0.definition", ""),
					resource.TestCheckResourceAttr("data.bridgecrew_tags.test", "tags.0.createdby", ""),
				),
			},
		},
	})
}

func testAccDataSourceTags() string {
	return `
	data "bridgecrew_tags" "test" {
	}`
}
