resource "bridgecrew_complex_policy" "new" {
  cloud_provider = "aws"
  title          = "my second test ${random_string.complex.id}"
  severity       = "critical"
  category       = "logging"
  frameworks     = ["Terraform"]

  conditionquery {
    and {
      or {
        resource_types = ["aws_instance"]
        cond_type      = "attribute"
        attribute      = "instance_type"
        operator       = "equals"
        value          = "t2.micro"
      }
      or {
        resource_types = ["aws_instance"]
        cond_type      = "attribute"
        attribute      = "instance_type"
        operator       = "equals"
        value          = "t3.micro"
      }
    }

    and {
      resource_types = ["aws_instance"]
      cond_type      = "attribute"
      attribute      = "name"
      operator       = "not_equals"
      value          = "jimbo"
    }

  }

  guidelines = "This should explain a lot more, in fact im padding this out to at least 50 characters"

  // although benchmarks take a free text this is total ***, as it needs to be an existing benchmark as
  // does the version, and that more like a category than anything
  benchmarks {
    cis_aws_v12 = ["1.1", "2.1"]
    cis_aws_v13 = ["1.3", "2.4"]
  }

}

resource "random_string" "complex" {
  length  = 8
  special = false
}
