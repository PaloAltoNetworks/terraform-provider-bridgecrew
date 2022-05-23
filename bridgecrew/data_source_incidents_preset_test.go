package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceIncidentsPreset(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIncidentsPreset(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_preset.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_preset.test", "presets.0.counter"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_preset.test", "presets.0.description"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_preset.test", "presets.0.id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_preset.test", "presets.0.isselected"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_incidents_preset.test", "presets.0.name"),
				),
			},
		},
	})
}

func testAccDataSourceIncidentsPreset() string {
	return `
	data "bridgecrew_incidents_preset" "test" {
	}`
}
