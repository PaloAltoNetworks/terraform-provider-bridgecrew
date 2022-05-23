package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataApiTokens(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAPITokens(),
			},
		},
	})
}

func testAccDataAPITokens() string {
	return `
	data "bridgecrew_apitokens" "test" {
	}`
}

func TestAccAPITokensDataSource_basic(t *testing.T) {
	// resourceName := "bridgecrew_c.test"
	// dataSourceName := "data.bridgecrew_apitokens.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAPITokens(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_apitokens.test", "apitokens.0.alias"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_apitokens.test", "apitokens.0.createdon"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_apitokens.test", "apitokens.0.userid"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_apitokens.test", "apitokens.0.uuid"),
				),
			},
		},
	})
}
