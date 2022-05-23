package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceIncidentsInfo(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIncidentsInfo(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "total"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "encryption.0.encrypted"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "encryption.0.noencryption"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "encryption.0.unencrypted"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "reachability.0.noreachability"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "reachability.0.private"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "reachability.0.public"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "status.0.open"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "status.0.passed"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "status.0.suppressed"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "traced.0.nottraced"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_info.test", "traced.0.traced"),
				),
			},
		},
	})
}

func testAccDataSourceIncidentsInfo() string {
	return `
	data "bridgecrew_incidents_info" "test" {
	}`
}
