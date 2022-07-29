resource "bridgecrew_policy" "new" {
  count          = 1
  cloud_provider = "aws"

  //CloudFomration isnt optional for some reason and if you dont put it in well the api adds it.
  frameworks = ["CloudFormation", "Terraform"]

  file             = "./terraform/policy/policy.yaml"
  source_code_hash = filesha256("${path.module}/policy/policy.yaml")

  // although benchmarks take a free text this is total ***, as it needs to be an existing benchmark as
  // does the version, and that more like a category than anything
  benchmarks {
    cis_aws_v12 = ["1.1", "2.1", "3.1"]
    //cis_aws_v13 = ["1.3", "2.4"]
  }

}
