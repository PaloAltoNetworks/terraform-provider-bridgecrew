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

func dataSourceApitokens() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApitokensRead,
		Schema: map[string]*schema.Schema{
			"apitokens": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"userid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"createdon": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

//goland:noinspection GoUnusedParameter,GoLinter,GoLinter
func dataSourceApitokensRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := RequestParams{"%s/api-tokens", "v1", "GET"}

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

	Apitokens := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&Apitokens)

	if err != nil {
		log.Fatal("Failed to parse data")
	}

	flatRepos := flattenApitokensData(&Apitokens)

	if err := d.Set("apitokens", flatRepos); err != nil {
		log.Fatal(reflect.TypeOf(Apitokens))
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenApitokensData(Apitokens *[]map[string]interface{}) []interface{} {
	if Apitokens != nil {
		ois := make([]interface{}, len(*Apitokens))

		for i, Apitoken := range *Apitokens {
			oi := make(map[string]interface{})

			oi["uuid"] = Apitoken["uuid"]
			oi["createdon"] = Apitoken["createdOn"]
			oi["description"] = Apitoken["description"]
			oi["alias"] = Apitoken["alias"]
			oi["userid"] = Apitoken["userId"]

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
