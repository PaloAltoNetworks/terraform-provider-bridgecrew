package bridgecrew

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/karlseguin/typed"
)

func dataSourcePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyRead,

		Schema: map[string]*schema.Schema{
			"policies": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_provider": {
							Type:        schema.TypeString,
							Computed:    false,
							Required:    true,
							Description: "The name of the Cloud Provider",
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"descriptive_title": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"constructive_title": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Required: true,
						},
						"category": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_types": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:    schema.TypeString,
								Default: "",
							},
						},
						"accountsdata": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"repository": {
										Required: true,
										Type:     schema.TypeString,
									},
									"amounts": {
										Required: true,
										Type:     schema.TypeMap,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"lastupdatedate": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"guideline": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iscustom": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"condition_query": {
							Type:     schema.TypeString,
							Required: true,
						},
						"benchmarks": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"benchmark": {
										Type:     schema.TypeString,
										Required: true,
									},
									"version": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"createdby": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

//goland:noinspection GoUnusedParameter
func dataSourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := RequestParams{"%s/policies/table/data", "v1", "GET"}

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

	body, _ := ioutil.ReadAll(r.Body)
	typedjson, err := typed.Json(body)

	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed at unmarshalling with typed %s \n", err.Error()),
		})
		return diagnostics
	}

	var data = typedjson.Maps("data")

	flatPolicies := flattenPolicyData(&data)

	if err := d.Set("policies", flatPolicies); err != nil {
		log.Fatal(reflect.TypeOf(data))
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenPolicyData(Policies *[]map[string]interface{}) []interface{} {
	if Policies != nil {
		ois := make([]interface{}, len(*Policies))
		for i, Policy := range *Policies {
			oi := make(map[string]interface{})

			oi["cloud_provider"] = Policy["provider"] // AWS
			oi["id"] = Policy["id"]                   // james_AWS_1620660945849
			oi["title"] = Policy["title"]             // new policy
			oi["descriptive_title"] = Policy["descriptiveTitle"]
			oi["constructive_title"] = Policy["constructiveTitle"]
			oi["severity"] = Policy["severity"] //CRITICAL
			oi["category"] = Policy["category"] // General
			oi["guideline"] = Policy["guideline"]
			oi["iscustom"] = Policy["isCustom"]

			var accounts []interface{}

			accountsData := Policy["accountsData"].(map[string]interface{})
			if len(accountsData) > 0 {
				for key, element := range accountsData {
					account := make(map[string]interface{})
					account["repository"] = key
					temp := element.(map[string]interface{})
					account["amounts"] = temp["amounts"]
					account["lastupdatedate"] = temp["lastUpdateDate"]
					accounts = append(accounts, account)
				}

				oi["accountsdata"] = accounts
			}

			var marks []interface{}

			benchmarks := Policy["benchmarks"].(map[string]interface{})
			for key, element := range benchmarks {
				bench := make(map[string]interface{})
				bench["benchmark"] = key
				bench["version"] = element
				marks = append(marks, bench)
			}

			oi["benchmarks"] = marks

			// this was hard... dynamic schema!
			u, err := json.Marshal(Policy["conditionQuery"])
			if err != nil {
				panic(err)
			}

			oi["condition_query"] = string(u)
			oi["resource_types"] = Policy["resourceTypes"]
			oi["createdby"] = Policy["createdBy"]
			oi["code"] = Policy["code"]
			ois[i] = oi
		}
		return ois
	}

	return make([]interface{}, 0)
}
