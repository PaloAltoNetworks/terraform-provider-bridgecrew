package bridgecrew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceJustifications(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceJustifications(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "accounts.0"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "policyid"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "justifications.0.active"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "justifications.0.comment"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "justifications.0.customer"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "justifications.0.date"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "justifications.0.id"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "justifications.0.origin"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "justifications.0.owner"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "justifications.0.suppression_type"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "justifications.0.type"),
					resource.TestCheckResourceAttrSet("data.bridgecrew_justifications.test", "justifications.0.violation_id"),
				),
			},
		},
	})
}

func testAccDataSourceJustifications() string {
	return `
	data "bridgecrew_justifications" "test" {
  		policyid="james_aws_1643121179054"
  		accounts=["JamesWoolfenden/full-fast-fail", "JamesWoolfenden/terraform-aws-s3"]
	}`
}
