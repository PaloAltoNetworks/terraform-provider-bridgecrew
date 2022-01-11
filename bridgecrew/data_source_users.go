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
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,

		Schema: map[string]*schema.Schema{
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accounts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"lastmodified": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"customername": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := RequestParams{"%s/manage/users", "v1", "GET"}

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
	var jsonMap []map[string]interface{}
	err = json.Unmarshal(body, &jsonMap)

	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed at unmarshalling with typed %s \n", err.Error()),
		})
		return diagnostics
	}

	flatUsers := flattenUserData(&jsonMap)

	if err := d.Set("users", flatUsers); err != nil {
		log.Fatal(reflect.TypeOf(jsonMap))
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenUserData(Users *[]map[string]interface{}) []interface{} {
	if Users != nil {
		ois := make([]interface{}, len(*Users))
		for i, User := range *Users {
			oi := make(map[string]interface{})
			oi["role"] = User["role"]
			oi["email"] = User["email"]
			oi["lastmodified"] = User["last_modified"]
			oi["accounts"] = User["accounts"]
			oi["customername"] = User["customer_name"]
			ois[i] = oi
		}
		return ois
	}

	return make([]interface{}, 0)
}
