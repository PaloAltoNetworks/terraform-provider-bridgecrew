package bridgecrew

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnforcementAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnforcementAccountsRead,
		Schema: map[string]*schema.Schema{
			"accounts": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
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
func dataSourceEnforcementAccountsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	request := "%s/enforcement-rules/accounts"

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

	Enforcement := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&Enforcement)

	if err != nil {
		if err.Error() == "EOF" {
			temp := "no data found"
			err = errors.New(temp)
			log.Print(temp)
		} else {
			log.Println("Failed to parse data")
		}
		return diag.FromErr(err)
	}

	if err := flattenEnforcementAccounts(Enforcement, d); err != nil {
		return err
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diagnostics
}

func flattenEnforcementAccounts(Enforcement []map[string]interface{}, d *schema.ResourceData) diag.Diagnostics {
	var accounts []interface{}
	if Enforcement != nil {
		for _, account := range Enforcement {
			myaccount := make(map[string]interface{})
			myaccount["account_id"] = account["accountId"]
			myaccount["account_name"] = account["accountName"]
			myaccount["source"] = account["source"]
			accounts = append(accounts, myaccount)
		}

		if err := d.Set("accounts", accounts); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
