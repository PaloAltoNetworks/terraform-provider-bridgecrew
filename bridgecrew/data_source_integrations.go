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

func dataSourceIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIntegrationRead,

		Schema: map[string]*schema.Schema{
			"integrations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"integration_details": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"params": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sf_execution_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

//goland:noinspection GoUnusedParameter
func dataSourceIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := RequestParams{"%s/integrations", "v1", "GET"}

	configure := m.(ProviderConfig)
	client, req, diagnostics, done := authClient(params, configure, nil)

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

	flatIntegrations := flattenIntegrationData(&data)

	if err := d.Set("integrations", flatIntegrations); err != nil {
		log.Fatal(reflect.TypeOf(data))
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenIntegrationData(integrations *[]map[string]interface{}) []interface{} {
	if integrations != nil {
		ois := make([]interface{}, len(*integrations))
		for i, Integration := range *integrations {
			oi := make(map[string]interface{})
			oi["enable"] = Integration["enable"]
			oi["id"] = Integration["id"]
			oi["sf_execution_name"] = Integration["sf_execution_name"]
			oi["status"] = Integration["status"]
			oi["type"] = Integration["type"]

			// cheat and use json
			jsoned, _ := json.Marshal(Integration["params"])
			oi["params"] = string(jsoned)

			// and again
			jsdetails, _ := json.Marshal(Integration["integration_details"])
			oi["integration_details"] = string(jsdetails)
			ois[i] = oi
		}
		return ois
	}

	return make([]interface{}, 0)
}
