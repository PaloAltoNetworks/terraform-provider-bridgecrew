resource "bridgecrew_simple_policy" "new" {
  count          = 1
  cloud_provider = "aws"
  title          = "my second simple test ${count.index} ${lower(random_string.new.id)}"
  severity       = "critical"
  category       = "general"

  //api always adds cloudformation
  frameworks = ["Terraform", "CloudFormation"]


  conditions {
    resource_types = ["aws_s3_bucket"]
    cond_type      = "attribute"
    attribute      = "bucket"
    operator       = "equals"
    value          = "jimbo2"
  }

  guidelines = "This should explain a lot more, in fact im padding this out to at least 50 characters"

  // although benchmarks take a free text this is total ***, as it needs to be an existing benchmark as
  // does the version, and that more like a category than anything
  benchmarks {
    cis_aws_v12 = ["1.1", "2.1", "3.1"]
    cis_aws_v13 = ["1.3", "2.4"]
  }

}


resource "random_string" "new" {
  length  = 8
  special = false
}
