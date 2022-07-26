package bridgecrew

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTag() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTagRead,
		Schema: map[string]*schema.Schema{
			// Input.
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag ID",
			},

			// Output.
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
	}
}

//goland:noinspection GoUnusedParameter,GoLinter,GoLinter
func dataSourceTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)

	request := "%s/tag-rules/" + id

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

	Tag := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&Tag)

	if err != nil {
		if err.Error() == "EOF" {
			temp := fmt.Sprintf("no data found for id: %s", id)
			err = errors.New(temp)
			log.Print(temp)
		} else {
			log.Println("Failed to parse data")
		}
		return diag.FromErr(err)
	}

	if err := flattenTag(Tag, d, id); err != nil {
		return err
	}

	return diagnostics
}

func flattenTag(tag map[string]interface{}, d *schema.ResourceData, id string) diag.Diagnostics {
	name := tag["name"].(string)
	description := tag["description"].(string)
	candoactions := tag["canDoActions"].(bool)

	if tag["createdBy"] != nil {
		createdby := tag["createdBy"].(string)
		if err := d.Set("createdby", createdby); err != nil {
			return diag.FromErr(err)
		}
	}

	if tag["definition"] != nil {
		u, err := json.Marshal(tag["definition"])
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("definition", string(u)); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(tag["repositories"].([]interface{})) != 0 {
		repositories := tag["repositories"].([]interface{})
		processed := make([]interface{}, len(repositories))
		for i, repo := range repositories {
			myrepo := repo.(map[string]interface{})
			oi := make(map[string]interface{})
			oi["id"] = myrepo["id"].(string)
			oi["name"] = myrepo["name"].(string)
			oi["source"] = myrepo["source"].(string)
			oi["owner"] = myrepo["owner"].(string)
			oi["repo"] = myrepo["repo"].(string)
			oi["defaultbranch"] = myrepo["defaultBranch"].(string)
			processed[i] = oi
		}

		if err := d.Set("repositories", processed); err != nil {
			return diag.FromErr(err)
		}
	}

	if tag["tagRuleOOTBId"] != nil {
		tagruleootbid := tag["tagRuleOOTBId"].(string)
		if err := d.Set("tagruleootbid", tagruleootbid); err != nil {
			return diag.FromErr(err)
		}
	}

	creationdate := tag["creationDate"].(string)
	isenabled := tag["isEnabled"].(bool)
	if err := d.Set("description", description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("candoactions", candoactions); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creationdate", creationdate); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("isenabled", isenabled); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return nil
}
