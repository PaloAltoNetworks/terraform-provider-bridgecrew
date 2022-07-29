resource "bridgecrew_enforcement_rule" "example" {
  name = "examplea"

  repositories {
    account_id   = "ed6fcaa9-86ba-46b8-93be-58cda8b74fd9"
    account_name = "james_cli_repo/alicloud"
  }

  repositories {
    account_id   = "b03361bc-e9f9-4108-91ef-ae150f8d12c6"
    account_name = "JamesWoolfenden/ansible-target-multi"
  }

  code_categories {
    supply_chain {
      soft_fail_threshold    = "HIGH"
      hard_fail_threshold    = "MEDIUM"
      comments_bot_threshold = "LOW"
    }
    secrets {
      soft_fail_threshold    = "HIGH"
      hard_fail_threshold    = "HIGH"
      comments_bot_threshold = "LOW"
    }
    iac {
      soft_fail_threshold    = "HIGH"
      hard_fail_threshold    = "HIGH"
      comments_bot_threshold = "LOW"
    }
    images {
      soft_fail_threshold    = "HIGH"
      hard_fail_threshold    = "HIGH"
      comments_bot_threshold = "LOW"
    }
    open_source {
      soft_fail_threshold    = "HIGH"
      hard_fail_threshold    = "HIGH"
      comments_bot_threshold = "LOW"
    }
  }
}
