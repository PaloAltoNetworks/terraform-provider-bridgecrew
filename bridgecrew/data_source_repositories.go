package bridgecrew

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRepositoryRead,
		Schema: map[string]*schema.Schema{
			"repositories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creationdate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Required: true,
						},
						"defaultbranch": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ispublic": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func dataSourceRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := RequestParams{"%s/repositories", "v1", "GET"}

	configure := m.(ProviderConfig)
	client, req, diagnostics, done, err := authClient(params, configure)

	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed at authClient %s \n", err.Error()),
		})
		return diagnostics
	}

	if done {
		return diagnostics
	}

	r, err := client.Do(req)
	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed at client.Do %s \n", err.Error()),
		})
		return diagnostics
	}

	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	repositories := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&repositories)

	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to parse data %s \n", err.Error()),
		})
		return diagnostics
	}

	flatRepos := flattenRepositoryData(&repositories)

	if err := d.Set("repositories", flatRepos); err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Repositories didn't set %s \n", err.Error()),
		})
		return diagnostics
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenRepositoryData(Repositories *[]map[string]interface{}) []interface{} {
	if Repositories != nil {
		ois := make([]interface{}, len(*Repositories), len(*Repositories))

		for i, Repository := range *Repositories {
			oi := make(map[string]interface{})
			oi["id"] = Repository["id"]
			oi["repository"] = Repository["repository"]
			oi["source"] = Repository["source"]
			oi["owner"] = Repository["owner"]
			oi["creationdate"] = Repository["creationDate"]
			oi["defaultbranch"] = Repository["defaultBranch"]
			oi["ispublic"] = Repository["ispublic"]

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
