package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUsers(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUsers(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_users.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_users.test", "users.0.customername"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_users.test", "users.0.email"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_users.test", "users.0.lastmodified"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_users.test", "users.0.role"),
				),
			},
		},
	})
}

func testAccUsers() string {
	return `
	data "bridgecrew_users" "test" {
	}`
}
