package bridgecrew

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProviderConfig vars for endpoint and API key
type ProviderConfig struct {
	URL         string
	Token       string
	AccessKeyID string
	SecretKey   string
	Prisma      string
}

// Provider main object
// goland:noinspection GoDeprecation
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
			"accesskeyid": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("PRISMA_ACCESS_KEY_ID", nil),
				Description: "Access key for Prisma",
			},
			"secretkey": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("PRISMA_SECRET_KEY", nil),
				Description: "Secret Key for Prisma",
			},
			"prisma": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PRISMA_API_URL", nil),
				Description: "URL for the Prisma, if set overrides the URL",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bridgecrew_policy":         resourcePolicy(),
			"bridgecrew_simple_policy":  resourceSimplePolicy(),
			"bridgecrew_complex_policy": resourceComplexPolicy(),
			"bridgecrew_tag":            resourceTag(),
			//"bridgecrew_enforcement_rule": resourceEnforcementRule(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bridgecrew_repositories":         dataSourceRepositories(),
			"bridgecrew_repository_branches":  dataSourceRepositoryBranches(),
			"bridgecrew_suppressions":         dataSourceSuppressions(),
			"bridgecrew_policies":             dataSourcePolicies(),
			"bridgecrew_authors":              dataSourceAuthors(),
			"bridgecrew_apitokens":            dataSourceApitokens(),
			"bridgecrew_apitokens_customer":   dataSourceApitokensByCustomer(),
			"bridgecrew_integrations":         dataSourceIntegrations(),
			"bridgecrew_users":                dataSourceUsers(),
			"bridgecrew_incidents":            dataSourceIncidents(),
			"bridgecrew_incidents_info":       dataSourceIncidentsInfo(),
			"bridgecrew_incidents_preset":     dataSourceIncidentsPreset(),
			"bridgecrew_organisation":         dataSourceOrganisation(),
			"bridgecrew_mappings":             dataSourceMappings(),
			"bridgecrew_tag":                  dataSourceTag(),
			"bridgecrew_tags":                 dataSourceTags(),
			"bridgecrew_justifications":       dataSourceJustifications(),
			"bridgecrew_enforcement_rules":    dataSourceEnforcementRules(),
			"bridgecrew_enforcement_rule":     dataSourceEnforcementRule(),
			"bridgecrew_enforcement_accounts": dataSourceEnforcementAccounts(),
		},
	}
}

// providerConfigure parses the config into the Terraform provider metaobject
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	token := d.Get("token").(string)
	accessKeyID := d.Get("accesskeyid").(string)
	secretKey := d.Get("secretkey").(string)
	prisma := d.Get("prisma").(string)
	return newProvider(url, token, accessKeyID, secretKey, prisma)
}

// newProviderClient is a factory for creating ProviderClient structs
func newProvider(url, token, accessKeyID, secretKey, prisma string) (ProviderConfig, error) {
	p := ProviderConfig{
		URL:         url,
		Token:       token,
		AccessKeyID: accessKeyID,
		SecretKey:   secretKey,
		Prisma:      prisma,
	}

	return p, nil
}
