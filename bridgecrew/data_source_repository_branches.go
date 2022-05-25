package bridgecrew

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRepositoryBranches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRepositoryBranchRead,
		Schema: map[string]*schema.Schema{
			"branches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creationdate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"defaultbranch": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"repoowner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reponame": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRepositoryBranchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	owner := d.Get("repoowner").(string)
	reponame := d.Get("reponame").(string)

	request := "%s/repositories/branches?repoOwner=" + owner + "&repoName=" + reponame
	params := RequestParams{request, "v1", "GET"}

	configure := m.(ProviderConfig)
	client, req, diagnostics, done := authClient(params, configure, nil)

	if done {
		return diagnostics
	}

	r, err := client.Do(req)

	if err != nil {
		log.Print("Failed at client.Do")
		return diag.FromErr(err)
	}

	// goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	repositoriesbranches := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&repositoriesbranches)

	if err != nil {
		log.Fatal("Failed to parse data")
	}

	flatBranch := flattenBranchData(&repositoriesbranches)

	if err := d.Set("branches", flatBranch); err != nil {
		log.Fatal(reflect.TypeOf(repositoriesbranches))
	}

	if err := d.Set("source", repositoriesbranches["source"].(string)); err != nil {
		log.Fatal(reflect.TypeOf(repositoriesbranches))
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenBranchData(repos *map[string]interface{}) []interface{} {
	if repos != nil {
		ois := make([]interface{}, len(*repos))
		temp := *repos
		branches := temp["branches"].([]interface{})
		for i, Repository := range branches {
			oi := make(map[string]interface{})
			scratch := Repository.(map[string]interface{})
			oi["name"] = scratch["name"].(string)
			oi["creationdate"] = scratch["creationDate"].(string)
			oi["defaultbranch"] = scratch["defaultBranch"].(bool)

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
