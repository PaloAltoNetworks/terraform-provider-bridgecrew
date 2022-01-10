package bridgecrew

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProviderConfig vars for endpoint and API key
type ProviderConfig struct {
	URL   string
	Token string
}

//Provider main object
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Default:     "https://www.bridgecrew.cloud",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BRIDGECREW_URL", nil),
				Description: "URL for the Bridgecrew Platform",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("BRIDGECREW_API", nil),
				Description: "API Token for Bridgecrew",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bridgecrew_policy":         resourcePolicy(),
			"bridgecrew_simple_policy":  resourceSimplePolicy(),
			"bridgecrew_complex_policy": resourceComplexPolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bridgecrew_repositories":        dataSourceRepositories(),
			"bridgecrew_repository_branches": dataSourceRepositoryBranches(),
			"bridgecrew_suppressions":        dataSourceSuppressions(),
			"bridgecrew_policies":            dataSourcePolicies(),
			"bridgecrew_errors":              dataSourceErrors(),
			"bridgecrew_apitokens":           dataSourceApitokens(),
			"bridgecrew_apitokens_customer":  dataSourceApitokensByCustomer(),
			"bridgecrew_integrations":        dataSourceIntegrations(),
			"bridgecrew_users":               dataSourceUsers(),
			"bridgecrew_incidents":           dataSourceIncidents(),
			"bridgecrew_incidents_info":      dataSourceIncidentsInfo(),
			"bridgecrew_incidents_preset":    dataSourceIncidentsPreset(),
		},
	}
}

// providerConfigure parses the config into the Terraform provider meta-object
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	token := d.Get("token").(string)
	return newProvider(url, token)
}

// newProviderClient is a factory for creating ProviderClient structs
func newProvider(url, token string) (ProviderConfig, error) {
	p := ProviderConfig{
		URL:   url,
		Token: token,
	}

	return p, nil
}
