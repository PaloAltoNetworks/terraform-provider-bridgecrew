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

func dataSourceIncidentsPreset() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncidentsPresetRead,
		Schema: map[string]*schema.Schema{
			"presets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isselected": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filters": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"search": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"encryption": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"sources": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"range": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"from": {
													Type:     schema.TypeString,
													Required: true,
												},
												"to": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"istraced": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"benchmarks": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"categories": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"reachability": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"sort": {
										Type:     schema.TypeString,
										Required: true,
									},
									"status": {
										Type:     schema.TypeString,
										Required: true,
									},
									"severities": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"tags": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"counter": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIncidentsPresetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := RequestParams{"%s/incidents/preset", "v2", "GET"}

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

	Presets := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&Presets)

	if err != nil {
		log.Fatal("Failed to parse data")
	}

	flatRepos := flattenIncidentsPresetData(&Presets)

	if err := d.Set("presets", flatRepos); err != nil {
		log.Fatal(reflect.TypeOf(Presets))
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenIncidentsPresetData(Presets *map[string]interface{}) []interface{} {
	if Presets != nil {
		data := (*Presets)["data"].([]interface{})
		ois := make([]interface{}, len(data), len(data))

		for i, Preset := range data {
			rawdata := Preset.(map[string]interface{})
			oi := make(map[string]interface{})
			oi["name"] = rawdata["name"].(string)
			oi["description"] = rawdata["description"].(string)
			oi["id"] = rawdata["id"].(string)
			oi["counter"] = rawdata["counter"].(float64)
			oi["isselected"] = rawdata["isSelected"].(bool)

			filters := make([]interface{}, 1, 1)
			if rawdata["filters"] != nil {
				raw := rawdata["filters"].(map[string]interface{})
				myfilter := make(map[string]interface{})

				if raw["sources"] != nil {
					myfilter["sources"] = raw["sources"]
				}

				if raw["encryption"] != nil {
					myfilter["encryption"] = raw["encryption"].(string)
				}

				if raw["search"] != nil {
					myfilter["search"] = raw["search"].(string)
				}

				if raw["range"] != nil {
					ranges := make([]interface{}, 1, 1)
					myrange := raw["range"].(map[string]interface{})
					ranges[0] = myrange
					myfilter["range"] = ranges
				}

				if raw["isTraced"] != nil {
					myfilter["istraced"] = raw["isTraced"].(bool)
				}

				if raw["categories"] != nil {
					myfilter["categories"] = raw["categories"]
				}

				if raw["reachability"] != nil {
					myfilter["reachability"] = raw["reachability"].(string)
				}

				if raw["benchmarks"] != nil {
					myfilter["benchmarks"] = raw["benchmarks"]
				}

				if raw["sort"] != nil {
					myfilter["sort"] = raw["sort"].(string)
				}

				if raw["status"] != nil {
					myfilter["status"] = raw["status"].(string)
				}

				if raw["severities"] != nil {
					myfilter["severities"] = raw["severities"]
				}

				if raw["tags"] != nil {
					myfilter["tags"] = raw["tags"]
				}

				filters[0] = myfilter
				oi["filters"] = filters
			}

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
