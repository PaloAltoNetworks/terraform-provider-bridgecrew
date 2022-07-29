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

func dataSourceEnforcementRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnforcementRulesRead,
		Schema: map[string]*schema.Schema{
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creationdate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"createdby": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mainrule": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"editable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"codecategories": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"supply_chain": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"soft_fail_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hard_fail_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"comments_bot_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"secrets": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"soft_fail_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hard_fail_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"comments_bot_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"iac": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"soft_fail_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hard_fail_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"comments_bot_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"images": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"soft_fail_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hard_fail_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"comments_bot_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"open_source": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"soft_fail_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hard_fail_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"comments_bot_threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"repositories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"accountsnotinmainrule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"account_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

//goland:noinspection GoUnusedParameter,GoLinter,GoLinter
func dataSourceEnforcementRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	request := "%s/enforcement-rules/"

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

	Enforcement := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&Enforcement)

	if err != nil {
		if err.Error() == "EOF" {
			err = errors.New("no data found")
		} else {
			log.Println("Failed to parse data")
		}
		return diag.FromErr(err)
	}

	if err := flattenEnforcementRules(Enforcement, d); err != nil {
		return err
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diagnostics
}

func flattenEnforcementRules(enforcement map[string]interface{}, d *schema.ResourceData) diag.Diagnostics {
	var themrules []interface{}
	if enforcement != nil {
		if enforcement["rules"] != nil {
			rules := enforcement["rules"].([]interface{})
			therule := make(map[string]interface{})
			for _, rule := range rules {
				myrule := rule.(map[string]interface{})
				therule["id"] = myrule["id"].(string)
				therule["creationdate"] = myrule["creationDate"].(string)
				therule["name"] = myrule["name"].(string)
				therule["createdby"] = myrule["createdBy"].(string)
				therule["mainrule"] = myrule["mainRule"].(bool)
				therule["editable"] = myrule["editable"].(bool)
				therule["repositories"] = myrule["repositories"].([]interface{})

				codecategories := make([]interface{}, 0)
				mycat := make(map[string]interface{})
				mycode := myrule["codeCategories"].(map[string]interface{})

				supplies := setcategories(mycode, "SUPPLY_CHAIN")
				secrets := setcategories(mycode, "SECRETS")
				iac := setcategories(mycode, "IAC")
				images := setcategories(mycode, "IMAGES")
				opensource := setcategories(mycode, "OPEN_SOURCE")

				mycat["supply_chain"] = supplies
				mycat["secrets"] = secrets
				mycat["iac"] = iac
				mycat["images"] = images
				mycat["open_source"] = opensource

				codecategories = append(codecategories, mycat)
				therule["codecategories"] = codecategories
				themrules = append(themrules, therule)
			}
			if err := d.Set("rules", themrules); err != nil {
				return diag.FromErr(err)
			}
		}

		if enforcement["accountsNotInMainRule"] != nil {
			accounts := enforcement["accountsNotInMainRule"].([]interface{})
			Repos := make([]interface{}, len(accounts))
			for _, account := range accounts {
				theNot := make(map[string]interface{})
				myAccount := account.(map[string]interface{})
				theNot["account_name"] = myAccount["accountName"].(string)
				theNot["account_id"] = myAccount["accountId"].(string)
				Repos = append(Repos, theNot)
			}
			if err := d.Set("accountsnotinmainrule", Repos); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func setcategories(mycode map[string]interface{}, node string) []interface{} {
	secrets := make([]interface{}, 0)
	secretchain := make(map[string]interface{})
	mysecrets := mycode[node].(map[string]interface{})
	secretchain["soft_fail_threshold"] = mysecrets["softFailThreshold"].(string)
	secretchain["hard_fail_threshold"] = mysecrets["hardFailThreshold"].(string)
	secretchain["comments_bot_threshold"] = mysecrets["commentsBotThreshold"].(string)
	secrets = append(secrets, secretchain)
	return secrets
}
