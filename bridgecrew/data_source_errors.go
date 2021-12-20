package bridgecrew

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceErrors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceErrorRead,
		Schema: map[string]*schema.Schema{
			"errors": {
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

func dataSourceErrorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := RequestParams{"%s/errors/gitBlameAuthors", "v1", "GET"}

	configure := m.(ProviderConfig)
	client, req, diagnostics, done, err := authClient(params, configure)

	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed at authClient %s \n", err.Error()),
		})
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

	errors := make([]map[string]interface{}, 0)
	var strerrors string

	err = json.NewDecoder(r.Body).Decode(&strerrors)

	//TODO:
	// errors method actually requires parameters to be supplied
	// fullRepoName: repository, sourceType
	// Published example lacks information on these

	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to parse data %s \n", err.Error()),
		})

		return diagnostics
	}

	flatErrors := flattenErrorData(&errors)

	if err := d.Set("errors", flatErrors); err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Set Errors failed %s \n", err.Error()),
		})
		return diagnostics
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenErrorData(Errors *[]map[string]interface{}) []interface{} {
	if Errors != nil {
		ois := make([]interface{}, len(*Errors), len(*Errors))

		for i, Error := range *Errors {
			oi := make(map[string]interface{})
			log.Print(Error)
			//todo
			/*oi["id"] = Repository["id"]
			oi["repository"] = Repository["repository"]
			oi["source"] = Repository["source"]
			oi["owner"] = Repository["owner"]
			oi["creationdate"] = Repository["creationDate"]
			oi["defaultbranch"] = Repository["defaultBranch"]
			oi["ispublic"] = Repository["ispublic"]*/

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
