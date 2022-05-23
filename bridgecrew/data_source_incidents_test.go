package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceIncidents(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIncidents(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents.test", "incidents.0.provider"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents.test", "incidents.0.guideline"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents.test", "incidents.0.incident_id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents.test", "incidents.0.iscustom"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents.test", "incidents.0.severity"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents.test", "incidents.0.title"),
				),
			},
		},
	})
}

func testAccDataSourceIncidents() string {
	return `
	data "bridgecrew_incidents" "test" {
	}`
}
