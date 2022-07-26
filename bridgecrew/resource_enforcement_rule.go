package bridgecrew

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnforcementRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnforcementRuleCreate,
		ReadContext:   resourceEnforcementRuleRead,
		UpdateContext: resourceEnforcementRuleUpdate,
		DeleteContext: resourceEnforcementRuleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the exception rule",
			},
			"repositories": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accountid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"accountname": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"code_categories": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"supply_chain": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"soft_fail_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
									"hard_fail_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
									"comments_bot_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
								},
							},
						},
						"secrets": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"soft_fail_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
									"hard_fail_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
									"comments_bot_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
								},
							},
						},
						"iac": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"soft_fail_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
									"hard_fail_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
									"comments_bot_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
								},
							},
						},
						"images": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"soft_fail_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
									"hard_fail_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
									"comments_bot_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
								},
							},
						},
						"open_source": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"soft_fail_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
									"hard_fail_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
									"comments_bot_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: ValidateThreshold,
									},
								},
							},
						},
					},
				},
			},
			"createdby": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Record created by",
			},
			"creationdate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Record created date",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Record created modified",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceEnforcementRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	myRule, err := setEnforcementRule(d)

	if err != nil {
		return diag.FromErr(err)
	}

	jsPolicy, err := json.Marshal(myRule)
	if err != nil {
		return diag.FromErr(err)
	}

	payload := strings.NewReader(string(jsPolicy))
	params := RequestParams{"%s/enforcement-rules", "v1", "POST"}
	configure := m.(ProviderConfig)
	client, req, diagnostics, done := authClient(params, configure, payload)

	if done {
		return diagnostics
	}

	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	diagnostics, fail := CheckStatus(res)

	if fail {
		return diagnostics
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	newResults, d2, done := VerifyReturn(body)
	if done {
		return d2
	}

	d.SetId(newResults.ID)

	resourceEnforcementRuleRead(ctx, d, m)

	return diags
}

func setEnforcementRule(d *schema.ResourceData) (Rule, error) {
	myRule := Rule{}
	myRule.Name = d.Get("name").(string)

	if d.Get("repositories") != nil {
		repos := d.Get("repositories").(*schema.Set)
		if repos.Len() > 0 {
			myRepos := repos.List()
			for _, myRepo := range myRepos {
				var Repo Repo
				temp := myRepo.(map[string]interface{})
				if temp["accountid"] != nil {
					Repo.AccountID = temp["accountid"].(string)
				}
				if temp["accountname"] != nil {
					Repo.AccountName = temp["accountname"].(string)
				}

				myRule.Repositories = append(myRule.Repositories, Repo)
			}
		}
	}
	if d.Get("code_categories") != nil {
		codeCategories := d.Get("code_categories").(*schema.Set)
		temp := codeCategories.List()
		test1 := temp[0].(map[string]interface{})

		Secrets := SetToMap(test1, "secrets")
		myRule.CodeCategories.Secrets.CommentsBotThreshold = Secrets["comments_bot_threshold"].(string)
		myRule.CodeCategories.Secrets.SoftFailThreshold = Secrets["soft_fail_threshold"].(string)
		myRule.CodeCategories.Secrets.HardFailThreshold = Secrets["hard_fail_threshold"].(string)

		SupplyChain := SetToMap(test1, "supply_chain")
		myRule.CodeCategories.SupplyChain.CommentsBotThreshold = SupplyChain["comments_bot_threshold"].(string)
		myRule.CodeCategories.SupplyChain.SoftFailThreshold = SupplyChain["soft_fail_threshold"].(string)
		myRule.CodeCategories.SupplyChain.HardFailThreshold = SupplyChain["hard_fail_threshold"].(string)

		IAC := SetToMap(test1, "iac")
		myRule.CodeCategories.IAC.CommentsBotThreshold = IAC["comments_bot_threshold"].(string)
		myRule.CodeCategories.IAC.SoftFailThreshold = IAC["soft_fail_threshold"].(string)
		myRule.CodeCategories.IAC.HardFailThreshold = IAC["hard_fail_threshold"].(string)

		OpenSource := SetToMap(test1, "open_source")
		myRule.CodeCategories.OpenSource.CommentsBotThreshold = OpenSource["comments_bot_threshold"].(string)
		myRule.CodeCategories.OpenSource.SoftFailThreshold = OpenSource["soft_fail_threshold"].(string)
		myRule.CodeCategories.OpenSource.HardFailThreshold = OpenSource["hard_fail_threshold"].(string)

		Images := SetToMap(test1, "open_source")
		myRule.CodeCategories.Images.CommentsBotThreshold = Images["comments_bot_threshold"].(string)
		myRule.CodeCategories.Images.SoftFailThreshold = Images["soft_fail_threshold"].(string)
		myRule.CodeCategories.Images.HardFailThreshold = Images["hard_fail_threshold"].(string)
	}

	return myRule, nil
}

// SetToMap converts SchemaSet to Map
func SetToMap(test1 map[string]interface{}, item string) map[string]interface{} {
	secrets := (test1[item].(*schema.Set)).List()[0]
	test2 := secrets.(map[string]interface{})
	return test2
}

//goland:noinspection GoUnusedParameter
func resourceEnforcementRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	ID := d.Id()

	params := RequestParams{"%s/enforcement-rules/", "v1", "GET"}
	configure := m.(ProviderConfig)
	client, req, diagnostics, done := authClient(params, configure, nil)

	if done {
		return diagnostics
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	Enforcements := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&Enforcements)

	if err != nil {
		if err.Error() == "EOF" {
			err = errors.New("no data found")
		} else {
			log.Println("Failed to parse data")
		}
		return diag.FromErr(err)
	}

	log.Print(ID)
	// todo this doesnt work -really need for there to be a working read method
	// parsing through the enforcement list fails to return the recently create record

	return diags
}

func resourceEnforcementRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if enforcementRuleChange(d) {
		myRule, err := setEnforcementRule(d)

		if err != nil {
			return diag.FromErr(err)
		}

		jsPolicy, err := json.Marshal(myRule)
		if err != nil {
			return diag.FromErr(err)
		}

		payload := strings.NewReader(string(jsPolicy))

		params := RequestParams{"%s/enforcement-rules/", "v1", "PUT"}
		configure := m.(ProviderConfig)
		client, req, diagnostics, done := authClient(params, configure, payload)

		if done {
			return diagnostics
		}

		res, err := client.Do(req)
		if err != nil {
			return diag.FromErr(err)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return diag.FromErr(err)
		}

		_, d2, done := VerifyReturn(body)

		if done {
			return d2
		}

		// no such field
		err = d.Set("last_updated", time.Now().Format(time.RFC850))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceEnforcementRuleRead(ctx, d, m)
}

func enforcementRuleChange(d *schema.ResourceData) bool {
	return d.HasChange("codecategories") ||
		d.HasChange("name") || d.HasChange("repositories") ||
		d.HasChange("codecategories")

}

//goland:noinspection GoUnusedParameter
func resourceEnforcementRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	ID := d.Id()
	configure := m.(ProviderConfig)
	params := RequestParams{"%s/enforcement-rules/" + ID, "v1", "DELETE"}
	client, req, diagnostics, done := authClient(params, configure, nil)

	if done {
		return diagnostics
	}

	res, err := client.Do(req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed at client.Do %s \n", err.Error()),
		})
		return diags
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		highlight(string(body))
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to delete %s \n", err.Error()),
		})
		return diags
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
