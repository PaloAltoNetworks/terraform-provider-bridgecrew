package bridgecrew

import (
	"context"
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

func dataSourceIncidents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncidentRead,

		Schema: map[string]*schema.Schema{
			"incidents": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"incident_id": {
							Type:        schema.TypeString,
							Description: "Typically the bridgecrew incident id, e.g. \"BC_AWS_IAM_37\"",
							Computed:    true,
						},
						"title": {
							Type:        schema.TypeString,
							Description: "The short check title",
							Required:    true,
						},
						"category": {
							Type:        schema.TypeString,
							Description: "The check category",
							Required:    true,
						},
						"guideline": {
							Type:        schema.TypeString,
							Description: "The URL to the checks description documentation",
							Computed:    true,
						},
						"severity": {
							Type:        schema.TypeString,
							Description: "Severity of the Incident",
							Required:    true,
						},
						"constructive_title": {
							Type:        schema.TypeString,
							Description: "A title but verbose- hopefully",
							Optional:    true,
						},
						"iscustom": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_types": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:    schema.TypeString,
								Default: "",
							},
						},
						"benchmarks": {
							Type:        schema.TypeList,
							Description: "The compliance framework this check/incident is against",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"remediation_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Describes the auto-remediation that is available for this",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"runtime_remediation": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"params": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The parameters",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Required: true,
									},
									"warning": {
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

//goland:noinspection GoUnusedParameter
func dataSourceIncidentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	configure := m.(ProviderConfig)
	var diagnostics diag.Diagnostics

	limit, offset := 100, 0
	var data []map[string]interface{}

	for {
		path := fmt.Sprintf("/incidents?limit=%d&offset=%d", limit, offset)
		params := RequestParams{"%s" + path, "v2", "POST"}

		client, req, tempDiagnostics, done := authClient(params, configure, nil)
		diagnostics = tempDiagnostics

		if done {
			return diagnostics
		}

		req.Header.Add("Content-Type", "application/json")
		r, err := client.Do(req)

		if err != nil {
			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed at client.Do %s \n", err.Error()),
			})
			return diagnostics
		}

		//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
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

		partialData := typedjson.Maps("data")
		data = append(data, partialData...)

		hasNext := typedjson.Bool("hasNext")
		if !hasNext {
			break
		}
		offset += limit
	}

	flatIncidents := flattenIncidentData(&data)

	if err := d.Set("incidents", flatIncidents); err != nil {
		log.Fatal(reflect.TypeOf(data))
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenIncidentData(incidents *[]map[string]interface{}) []interface{} {
	if incidents != nil {
		ois := make([]interface{}, len(*incidents))
		for i, Incident := range *incidents {
			oi := make(map[string]interface{})

			oi["incident_id"] = Incident["incidentId"]
			oi["title"] = Incident["title"]
			oi["constructive_title"] = Incident["constructiveTitle"]
			oi["severity"] = Incident["severity"]
			oi["category"] = Incident["category"]
			oi["guideline"] = Incident["guideline"]
			oi["iscustom"] = Incident["isCustom"]
			oi["provider"] = Incident["provider"]
			oi["resource_types"] = Incident["resourceTypes"]
			if keyExists(Incident, "benchmarks") {
				oi["benchmarks"] = Incident["benchmarks"]
			}

			if keyExists(Incident, "remediationIds") {
				oi["remediation_ids"] = Incident["remediationIds"]
			}

			if keyExists(Incident, "runtimeRemediation") {
				var remediations []interface{}

				remediateData := Incident["runtimeRemediation"].([]interface{})
				if len(remediateData) > 0 {
					for _, element := range remediateData {
						account := make(map[string]interface{})
						temp := element.(map[string]interface{})
						account["id"] = temp["id"]
						account["warning"] = temp["warning"]
						account["description"] = temp["description"]
						if keyExists(temp, "params") {
							account["params"], _ = temp["params"].([]string)
						}

						remediations = append(remediations, account)
					}

					oi["runtime_remediation"] = remediations
				}
			}

			ois[i] = oi
		}
		return ois
	}

	return make([]interface{}, 0)
}
