package bridgecrew

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/karlseguin/typed"
)

func resourceSimplePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSimplePolicyCreate,
		ReadContext:   resourceSimplePolicyRead,
		UpdateContext: resourceSimplePolicyUpdate,
		DeleteContext: resourceSimplePolicyDelete,
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
			"conditions": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
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
							Description:  "The logic Operator",
							ValidateFunc: ValidateOperator,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The Value to Check",
							Required:    true,
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

func resourceSimplePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := &http.Client{Timeout: 60 * time.Second}
	var diags diag.Diagnostics

	myPolicy, err := setSimplePolicy(d)

	if err != nil {
		return diag.FromErr(err)
	}

	jsPolicy, err := json.Marshal(myPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	configure := m.(ProviderConfig)
	url := configure.URL + "/api/v1/policies"
	payload := strings.NewReader(string(jsPolicy))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", configure.Token)

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

	newResults, d2, fail := VerifyReturn(err, body)
	if fail {
		return d2
	}

	d.SetId(newResults.Policy)

	resourceSimplePolicyRead(ctx, d, m)

	return diags
}

func setSimplePolicy(d *schema.ResourceData) (simplePolicy, error) {
	myPolicy := simplePolicy{}
	myBenchmark, err := setBenchmark(d)

	if err == nil {
		myPolicy.Benchmarks = myBenchmark
	}

	myPolicy.Category = d.Get("category").(string)

	conditions, err := setConditions(d)

	//Don't set if not set
	if err != nil {
		return myPolicy, fmt.Errorf("unable set conditions %q", err)
	}
	myPolicy.Conditions = conditions[0]

	myPolicy.Provider = d.Get("cloud_provider").(string)
	myPolicy.Severity = d.Get("severity").(string)
	myPolicy.Title = d.Get("title").(string)
	myPolicy.Guidelines = d.Get("guidelines").(string)
	myPolicy.Frameworks = CastToStringList(d.Get("frameworks").([]interface{}))

	return myPolicy, nil
}

func setConditions(d *schema.ResourceData) ([]Conditions, error) {
	conditions := make([]Conditions, 0, 1)

	myConditions := d.Get("conditions").([]interface{})

	if len(myConditions) > 0 {
		for _, myCondition := range myConditions {
			temp := myCondition.(map[string]interface{})
			var Condition Conditions
			Condition.Value = temp["value"].(string)
			Condition.CondType = temp["cond_type"].(string)
			Condition.Attribute = temp["attribute"].(string)
			Condition.Operator = temp["operator"].(string)

			var myResources []string
			myResources = CastToStringList(temp["resource_types"].([]interface{}))
			Condition.ResourceTypes = myResources

			conditions = append(conditions, Condition)
		}
	} else {
		return nil, errors.New("no Conditions Set")
	}

	return conditions, nil
}

func setBenchmark(d *schema.ResourceData) (Benchmark, error) {

	_, data := d.GetOk("benchmarks")
	var myItem Benchmark

	if data {
		myBenchmark := (d.Get("benchmarks").(*schema.Set)).List()

		s := myBenchmark[0].(map[string]interface{})
		myItem.Cisawsv12 = CastToStringList(s["cis_aws_v12"].([]interface{}))
		myItem.Cisawsv13 = CastToStringList(s["cis_aws_v13"].([]interface{}))
		myItem.Cisazurev11 = CastToStringList(s["cis_azure_v11"].([]interface{}))
		myItem.Cisazurev12 = CastToStringList(s["cis_azure_v12"].([]interface{}))
		myItem.Cisazurev13 = CastToStringList(s["cis_azure_v13"].([]interface{}))
		myItem.Cisgcpv11 = CastToStringList(s["cis_gcp_v11"].([]interface{}))
		myItem.Ciskubernetesv15 = CastToStringList(s["cis_kubernetes_v15"].([]interface{}))
		myItem.Ciskubernetesv16 = CastToStringList(s["cis_kubernetes_v16"].([]interface{}))
		myItem.Cisdockerv11 = CastToStringList(s["cis_docker_v11"].([]interface{}))
		myItem.Ciseksv11 = CastToStringList(s["cis_eks_v11"].([]interface{}))
		myItem.Cisgkev11 = CastToStringList(s["cis_gke_v11"].([]interface{}))
		return myItem, nil
	}

	return myItem, errors.New("no benchmark data")
}

func resourceSimplePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 60 * time.Second}

	policyID := d.Id()

	configure := m.(ProviderConfig)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/policies/%s", configure.URL, policyID), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// add authorization header to the req
	req.Header.Add("authorization", configure.Token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

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

	d.Set("cloud_provider", strings.ToLower(typedjson["provider"].(string)))
	d.Set("title", typedjson["title"].(string))
	d.Set("severity", strings.ToLower(typedjson["severity"].(string)))
	d.Set("category", strings.ToLower(typedjson["category"].(string)))
	d.Set("frameworks", typedjson["frameworks"])

	err = d.Set("guidelines", typedjson["guideline"])
	if err != nil {
		return diag.FromErr(err)
	}

	//myconditions should be an array it currently a map
	//hence this fudge
	//todo: once you start passing around condition arrays
	//this can go
	myConditions := make([]interface{}, 1)
	myConditions[0] = typedjson["conditionQuery"]
	err = d.Set("conditions", myConditions)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	return diags
}

func resourceSimplePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := &http.Client{Timeout: 60 * time.Second}

	policyID := d.Id()
	if simplepolicyChange(d) {
		myPolicy, err := setSimplePolicy(d)

		if err != nil {
			return diag.FromErr(err)
		}

		jsPolicy, err := json.Marshal(myPolicy)
		if err != nil {
			return diag.FromErr(err)
		}

		configure := m.(ProviderConfig)

		payload := strings.NewReader(string(jsPolicy))
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/policies/%s", configure.URL, policyID), payload)

		if err != nil {
			return diag.FromErr(err)
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("authorization", configure.Token)

		res, err := client.Do(req)
		if err != nil {
			return diag.FromErr(err)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return diag.FromErr(err)
		}

		_, d2, fail := VerifyReturn(err, body)
		if fail {
			return d2
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}
	return resourceSimplePolicyRead(ctx, d, m)
}

func simplepolicyChange(d *schema.ResourceData) bool {
	return d.HasChange("conditions") ||
		d.HasChange("cloud_provider") ||
		d.HasChange("title") ||
		d.HasChange("severity") ||
		d.HasChange("category") ||
		d.HasChange("guidelines") ||
		d.HasChange("benchmarks") ||
		d.HasChange("frameworks")
}

func resourceSimplePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 60 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policyID := d.Id()
	configure := m.(ProviderConfig)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/policies/%s", configure.URL, policyID), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", configure.Token)

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
