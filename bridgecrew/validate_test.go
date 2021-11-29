package bridgecrew

import (
	"testing"
)

func TestValidateOperator(t *testing.T) {
	validValues := []string{
		"equals", "not_equals", "regex_match", "not_reqex_match", "greater_than", "greater_than_or_equal",
		"less_than_or_equal", "less_than", "exists", "not_exists", "contains", "not_contains", "starting_with",
		"not_starting_with", "ending_with", "not_ending_with",
	}

	for _, s := range validValues {
		_, errors := ValidateOperator(s, "operator")
		if len(errors) > 0 {
			t.Fatalf("%q policy operator should have been valid: %v", s, errors)
		}
	}

	invalidValues := []string{
		"arse",
		"NOT_ENABLED",
		"white space",
		"EQUALS",
	}

	for _, s := range invalidValues {
		_, errors := ValidateOperator(s, "operator")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid policy operator", s)
		}
	}
}

func TestValidateCloudProvider(t *testing.T) {
	validValues := []string{
		"aws",
		"gcp",
		"linode",
		"azure",
		"oci",
		"alicloud",
		"digitalocean",
	}

	for _, s := range validValues {
		_, errors := ValidateCloudProvider(s, "cloud_provider")
		if len(errors) > 0 {
			t.Fatalf("%q cloud provider should have been valid: %v", s, errors)
		}
	}

	invalidValues := []string{
		"AWS",
		"arse",
		"NOT_ENABLED",
		"white space",
		"EQUALS",
	}

	for _, s := range invalidValues {
		_, errors := ValidateCloudProvider(s, "cloud_provider")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid cloud provider", s)
		}
	}
}

func TestValidateSeverity(t *testing.T) {
	validValues := []string{
		"critical",
		"high",
		"low",
		"medium",
	}

	for _, s := range validValues {
		_, errors := ValidateSeverity(s, "severity")
		if len(errors) > 0 {
			t.Fatalf("%q severity should have been valid: %v", s, errors)
		}
	}

	invalidValues := []string{
		"CRITICAL",
		"arse",
		"NOT_ENABLED",
		"white space",
		"EQUALS",
	}

	for _, s := range invalidValues {
		_, errors := ValidateSeverity(s, "severity")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid severity", s)
		}
	}
}

func TestValidateCategory(t *testing.T) {
	validValues := []string{
		"logging",
		"elasticsearch",
		"general",
		"storage",
		"encryption",
		"networking",
		"monitoring",
		"kubernetes",
		"serverless",
		"backup_and_recovery",
		"iam",
		"secrets",
		"public",
	}

	for _, s := range validValues {
		_, errors := ValidateCategory(s, "category")
		if len(errors) > 0 {
			t.Fatalf("%q category should have been valid: %v", s, errors)
		}
	}

	invalidValues := []string{
		"CRITICAL",
		"arse",
		"NOT_ENABLED",
		"white space",
		"EQUALS",
	}

	for _, s := range invalidValues {
		_, errors := ValidateCategory(s, "category")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid category", s)
		}
	}
}

func TestValidateGuidelines(t *testing.T) {
	validValues := []string{
		"A really long Guideline should be over 50 characters",
	}

	for _, s := range validValues {
		_, errors := ValidateGuidelines(s, "guideline")
		if len(errors) > 0 {
			t.Fatalf("%q guideline should have been valid: %v", s, errors)
		}
	}

	invalidValues := []string{
		"Terse",
		"is",
		"beauty",
	}

	for _, s := range invalidValues {
		_, errors := ValidateGuidelines(s, "guideline")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid guideline", s)
		}
	}
}

func TestValidatePolicyTitle(t *testing.T) {
	validValues := []string{
		"A descriptive title has 20 chars",
	}

	for _, s := range validValues {
		_, errors := ValidatePolicyTitle(s, "title")
		if len(errors) > 0 {
			t.Fatalf("%q title should have been valid: %v", s, errors)
		}
	}

	invalidValues := []string{
		"Brevity",
		"is",
		"Wit",
	}

	for _, s := range invalidValues {
		_, errors := ValidatePolicyTitle(s, "title")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid title", s)
		}
	}
}

func TestValidPolicyJSON(t *testing.T) {
	validValues := []string{
		"{\"name\":\"John\", \"age\":30, \"car\":null}",
		"{\"provider\":\"AWS\",\"id\":\"james_AWS_1620660945849\",\"title\":\"new policy\",\"descriptiveTitle\":null,\"constructiveTitle\":null,\"severity\":\"CRITICAL\",\"category\":\"General\",\"resourceTypes\":[\"aws_api_gateway_api_key\"]}"}
	for _, s := range validValues {
		_, errors := ValidPolicyJSON(s, "json")
		if len(errors) > 0 {
			t.Fatalf("%q json should have been valid: %v", s, errors)
		}
	}

	invalidValues := []string{
		"",
		"james",
	}

	for _, s := range invalidValues {
		_, errors := ValidatePolicyTitle(s, "json")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid json", s)
		}
	}
}
