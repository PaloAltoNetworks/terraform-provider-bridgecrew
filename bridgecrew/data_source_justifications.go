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

func dataSourceJustifications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJustificationsRead,
		Schema: map[string]*schema.Schema{
			"policyid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"accounts": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"justifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customer": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"suppression_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"violation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceJustificationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	policyid := d.Get("policyid").(string)
	accounts := d.Get("accounts").([]interface{})

	var query string
	for i, account := range accounts {
		query = query + "accounts=" + account.(string)
		if i < (len(accounts) - 1) {
			query = query + "&"
		}
	}

	url := "%s/suppressions/" + policyid + "/justifications?" + query

	params := RequestParams{url, "v1", "GET"}

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

	Justifications := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&Justifications)

	if err != nil {
		log.Fatal("Failed to parse data")
	}

	flatJustice := flattenJustificationsData(&Justifications)

	if err := d.Set("justifications", flatJustice); err != nil {
		log.Fatal(reflect.TypeOf(Justifications))
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenJustificationsData(Justifications *[]map[string]interface{}) []interface{} {
	if Justifications != nil {
		ois := make([]interface{}, len(*Justifications))

		for i, Justify := range *Justifications {
			oi := make(map[string]interface{})

			oi["customer"] = Justify["customer"]
			oi["id"] = Justify["id"]
			oi["date"] = Justify["date"].(float64)
			oi["owner"] = Justify["owner"]
			oi["comment"] = Justify["comment"]
			oi["suppression_type"] = Justify["suppressionType"]
			oi["violation_id"] = Justify["violationId"]
			oi["origin"] = Justify["origin"]
			oi["active"] = Justify["active"].(bool)
			oi["type"] = Justify["type"]

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
