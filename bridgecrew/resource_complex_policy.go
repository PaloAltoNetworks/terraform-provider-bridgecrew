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
	"github.com/karlseguin/typed"
)

func resourceComplexPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComplexPolicyCreate,
		ReadContext:   resourceComplexPolicyRead,
		UpdateContext: resourceComplexPolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Schema: map[string]*schema.Schema{
			"cloud_provider": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Computed:     false,
				Description:  "The Cloud provider this is for e.g. - aws, gcp, azure.",
				Required:     true,
				ValidateFunc: ValidateCloudProvider,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				Description:  "The title of the check, needs to be longer than 20 chars - an effort to ensure detailed names.",
				ValidateFunc: ValidatePolicyTitle,
			},
			"severity": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Severity category allows you to indicate importance and this value can determine build or PR failure in the platform.",
				ValidateFunc: ValidateSeverity,
			},
			"frameworks": {
				Type:        schema.TypeList,
				Description: "Which IAC framework is this policy targeting.",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"category": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Check category for grouping similar checks.",
				ValidateFunc: ValidateCategory,
			},
			"guidelines": {
				Type:     schema.TypeString,
				Required: true,
				Description: "A detailed description helps you understand why the check was written and should include details on how " +
					"to fix the violation. The field must more than 50 chars in it, to encourage detail.",
				ValidateFunc: ValidateGuidelines,
			},
			"conditionquery": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "The actual query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"and": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Conditions captures the actual check logic. Do not add resource_types and an or statement in the same block",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_types": {
										Type:        schema.TypeList,
										Description: "The resource type",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"cond_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "",
									},
									"attribute": {
										Type:        schema.TypeString,
										Description: "The field that you want the condition on",
										Optional:    true,
									},
									"operator": {
										Type:         schema.TypeString,
										Description:  "The logic operator",
										Optional:     true,
										ValidateFunc: ValidateOperator,
									},
									"value": {
										Description: "The value to check against",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"or": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Conditions captures the actual check logic",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource_types": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "The resource type",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"cond_type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"attribute": {
													Type:        schema.TypeString,
													Description: "The field that you want the condition on",
													Required:    true,
												},
												"operator": {
													Type:         schema.TypeString,
													Required:     true,
													Description:  "The logic operator",
													ValidateFunc: ValidateOperator,
												},
												"value": {
													Type:        schema.TypeString,
													Description: "The Value to check",
													Required:    true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"benchmarks": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "This associates the check to one or many compliance frameworks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cis_azure_v11": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_azure_v12": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_azure_v13": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_aws_v12": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_aws_v13": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_kubernetes_v15": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_kubernetes_v16": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_gcp_v11": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_gke_v11": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_docker_v11": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_eks_v11": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceComplexPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	myPolicy, err := setComplexPolicy(d)

	if err != nil {
		return diag.FromErr(err)
	}

	jsPolicy, err := json.Marshal(myPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	payload := strings.NewReader(string(jsPolicy))
	params := RequestParams{"%s/policies", "v1", "POST"}
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

	newResults, d2, fail := VerifyReturn(body)
	if fail {
		return d2
	}

	d.SetId(newResults.Policy)

	resourceComplexPolicyRead(ctx, d, m)

	return diags
}

func setComplexPolicy(d *schema.ResourceData) (complexPolicy, error) {
	myPolicy := complexPolicy{}
	myBenchmark, err := setBenchmark(d)

	if err == nil {
		myPolicy.Benchmarks = myBenchmark
	}

	myPolicy.Category = d.Get("category").(string)

	conditionQuery, err := setComplexConditions(d)

	//Don't set if not set
	if err != nil {
		return myPolicy, fmt.Errorf("unable set conditions %q", err)
	}
	myPolicy.ConditionQuery = conditionQuery
	myPolicy.Provider = d.Get("cloud_provider").(string)
	myPolicy.Severity = d.Get("severity").(string)
	myPolicy.Title = d.Get("title").(string)
	myPolicy.Guidelines = d.Get("guidelines").(string)
	myPolicy.Frameworks, _ = CastToStringList(d.Get("frameworks").([]interface{}))

	return myPolicy, nil
}

func setComplexConditions(d *schema.ResourceData) (ConditionQuery, error) {
	var conditionQuery ConditionQuery

	query := d.Get("conditionquery").(*schema.Set)

	if query.Len() > 0 {
		myQuery := query.List()
		duff := myQuery[0].(map[string]interface{})
		TheAnds := duff["and"].([]interface{})

		var conditions []Conditions
		var TheOrs Conditions

		for _, myCondition := range TheAnds {
			temp := myCondition.(map[string]interface{})

			var Condition Conditions

			if len(temp["or"].([]interface{})) > 0 {
				//have to have some way of ignoring root vars if you choose to make an or statement
				log.Print("Or value set, ignoring other vars in block")

				ors := temp["or"].([]interface{})
				var orConditions []Or

				//Process all the Or conditions
				for _, anor := range ors {
					localor := anor.(map[string]interface{})

					var orCondition Or
					orCondition.Value = localor["value"].(string)
					orCondition.CondType = localor["cond_type"].(string)
					orCondition.Attribute = localor["attribute"].(string)
					orCondition.Operator = localor["operator"].(string)

					orCondition.ResourceTypes, _ = CastToStringList(localor["resource_types"].([]interface{}))
					orConditions = append(orConditions, orCondition)
				}
				TheOrs.Or = orConditions
			} else {
				//process the condition if just the root Condition
				Condition.Value = temp["value"].(string)
				Condition.CondType = temp["cond_type"].(string)
				Condition.Attribute = temp["attribute"].(string)
				Condition.Operator = temp["operator"].(string)

				Condition.ResourceTypes, _ = CastToStringList(temp["resource_types"].([]interface{}))
				conditions = append(conditions, Condition)
			}
		}

		//append the Ors as complete coniditions
		conditions = append(conditions, TheOrs)

		//add the whole condition to the resource
		conditionQuery.Ands = conditions
	} else {
		return conditionQuery, errors.New("no Conditions Set")
	}

	return conditionQuery, nil
}

func resourceComplexPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	policyID := d.Id()

	params := RequestParams{"%s/policies/" + policyID, "v1", "GET"}
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

	body, _ := ioutil.ReadAll(r.Body)
	typedjson, err := typed.Json(body)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("cloud_provider", strings.ToLower(typedjson["provider"].(string)))
	diags = LogAppendError(err, diags)

	err = d.Set("title", typedjson["title"].(string))
	diags = LogAppendError(err, diags)

	err = d.Set("severity", strings.ToLower(typedjson["severity"].(string)))
	diags = LogAppendError(err, diags)

	err = d.Set("category", strings.ToLower(typedjson["category"].(string)))
	diags = LogAppendError(err, diags)

	err = d.Set("frameworks", typedjson["frameworks"])
	diags = LogAppendError(err, diags)

	err = d.Set("guidelines", typedjson["guideline"])
	diags = LogAppendError(err, diags)

	myConditions := make([]interface{}, 1)
	myConditions[0] = typedjson["conditionQuery"]
	err = d.Set("conditionquery", myConditions)
	diags = LogAppendError(err, diags)

	return diags
}

func resourceComplexPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	policyID := d.Id()
	if complexPolicyChange(d) {
		myPolicy, err := setComplexPolicy(d)

		if err != nil {
			return diag.FromErr(err)
		}

		jsPolicy, err := json.Marshal(myPolicy)
		if err != nil {
			return diag.FromErr(err)
		}

		payload := strings.NewReader(string(jsPolicy))

		params := RequestParams{"%s/policies/" + policyID, "v1", "PUT"}
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

		_, d2, fail := VerifyReturn(body)
		if fail {
			return d2
		}

		err = d.Set("last_updated", time.Now().Format(time.RFC850))

		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceComplexPolicyRead(ctx, d, m)
}

func complexPolicyChange(d *schema.ResourceData) bool {
	return d.HasChange("conditionquery") ||
		d.HasChange("cloud_provider") ||
		d.HasChange("title") ||
		d.HasChange("severity") ||
		d.HasChange("category") ||
		d.HasChange("guidelines") ||
		d.HasChange("benchmarks") ||
		d.HasChange("frameworks")
}
