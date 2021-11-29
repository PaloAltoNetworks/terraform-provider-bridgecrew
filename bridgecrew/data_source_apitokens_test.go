package bridgecrew

import (
	"fmt"
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
	return fmt.Sprintf(
		`
data "bridgecrew_apitokens" "test" {
}`)
}

func TestAccAPITokensDataSource_basic(t *testing.T) {
	//resourceName := "bridgecrew_c.test"
	//dataSourceName := "data.bridgecrew_apitokens.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAPITokens(),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}
