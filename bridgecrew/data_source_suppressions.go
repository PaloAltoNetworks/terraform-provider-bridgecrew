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

func dataSourceSuppressions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSuppressionRead,
		Schema: map[string]*schema.Schema{
			"suppressions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"suppressiontype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"creationdate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policyid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"accountid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"resourceid": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceSuppressionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := RequestParams{"%s/suppressions", "v1", "GET"}

	configure := m.(ProviderConfig)
	client, req, diagnostics, done, err := authClient(params, configure)

	if done {
		return diagnostics
	}

	r, err := client.Do(req)

	if err != nil {
		log.Fatal("Failed at client.Do")
	}

	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	Suppressions := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&Suppressions)

	if err != nil {
		log.Fatal("Failed to parse data")
	}

	flatRepos := flattenSuppressionData(&Suppressions)

	if err := d.Set("suppressions", flatRepos); err != nil {
		log.Fatal(reflect.TypeOf(Suppressions))
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenSuppressionData(Suppressions *[]map[string]interface{}) []interface{} {
	if Suppressions != nil {
		ois := make([]interface{}, len(*Suppressions))

		for i, Suppression := range *Suppressions {
			oi := make(map[string]interface{})
			oi["suppressiontype"] = Suppression["suppressionType"]
			oi["creationdate"] = Suppression["creationDate"]
			oi["id"] = Suppression["id"]
			oi["policyid"] = Suppression["policyId"]
			oi["comment"] = Suppression["comment"]

			var myresources []interface{}
			if Suppression["resources"] != nil {
				resources := Suppression["resources"].([]interface{})

				for _, element := range resources {

					account := element.(map[string]interface{})

					myaccount := make(map[string]interface{})
					myaccount["accountid"] = account["accountId"].(string)
					myaccount["resourceid"] = account["resourceId"]

					myresources = append(myresources, myaccount)
				}

				oi["resources"] = myresources
			}

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
