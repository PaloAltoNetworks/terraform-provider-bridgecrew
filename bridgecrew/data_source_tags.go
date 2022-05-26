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

func dataSourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTagsRead,
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified by",
						},
						"definition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"createdby": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag Author",
						},
						"creationdate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timestamp",
						},
						"isenabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Tag is enabled",
						},
						"tagruleootbid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "tagRuleOOTBId",
						},
						"repositories": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of repositories",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"owner": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"repo": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defaultbranch": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"candoactions": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Can Do Actions",
						},
					},
				},
			},
		},
	}
}

//goland:noinspection GoUnusedParameter,GoLinter,GoLinter
func dataSourceTagsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	request := "%s/tag-rules/"
	params := RequestParams{request, "v1", "GET"}

	configure := m.(ProviderConfig)
	client, req, diagnostics, done := authClient(params, configure, nil)

	if done {
		return diagnostics
	}

	r, err := client.Do(req)

	if err != nil {
		log.Fatal("Failed at client.Do")
	}

	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	Tags := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&Tags)

	if err != nil {
		log.Panic("Failed to parse data")
	}

	flatTags := flattenTags(&Tags)

	if err := d.Set("tags", flatTags); err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Tags didn't set %s \n", err.Error()),
		})
		return diagnostics
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenTags(Tags *[]map[string]interface{}) []interface{} {
	if Tags != nil {
		ois := make([]interface{}, len(*Tags))

		for i, Tag := range *Tags {
			oi := make(map[string]interface{})
			oi["id"] = Tag["id"]
			oi["name"] = Tag["name"]
			oi["description"] = Tag["description"]

			if Tag["definition"] == nil {
				oi["definition"] = Tag["definition"]
			} else {
				mydefinitionjson, _ := json.Marshal(Tag["definition"])
				oi["definition"] = string(mydefinitionjson)
			}

			oi["createdby"] = Tag["createdBy"]
			oi["creationdate"] = Tag["creationDate"]
			oi["isenabled"] = Tag["isEnabled"].(bool)
			oi["tagruleootbid"] = Tag["tagRuleOOTBId"]

			if len(Tag["repositories"].([]interface{})) != 0 {
				repositories := Tag["repositories"].([]interface{})
				processed := make([]interface{}, len(repositories))
				for i, repo := range repositories {
					myrepo := repo.(map[string]interface{})
					nuoi := make(map[string]interface{})
					nuoi["id"] = myrepo["id"].(string)
					nuoi["name"] = myrepo["name"].(string)
					nuoi["source"] = myrepo["source"].(string)
					nuoi["owner"] = myrepo["owner"].(string)
					nuoi["repo"] = myrepo["repo"].(string)
					nuoi["defaultbranch"] = myrepo["defaultBranch"].(string)
					processed[i] = nuoi
				}
				oi["repositories"] = processed
			} else {
				oi["repositories"] = Tag["repositories"]
			}
			oi["candoactions"] = Tag["canDoActions"]

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
