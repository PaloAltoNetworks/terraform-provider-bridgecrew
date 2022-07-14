package bridgecrew

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
			//"pcseverity": {
			//	Type:        schema.TypeString,
			//	Optional:    true,
			//	Description: "PRISMA severity category allows you to indicate importance and this value can determine build or PR failure in the platform.",
			//	//ValidateFunc: ValidateSeverity,
			//},
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
				Optional: true,
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

	var diags diag.Diagnostics

	myPolicy, err := setSimplePolicy(d)

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
	//myPolicy.PRISMASeverity = d.Get("pcseverity").(string)
	myPolicy.Title = d.Get("title").(string)
	myPolicy.Guidelines = d.Get("guidelines").(string)
	myPolicy.Frameworks, _ = CastToStringList(d.Get("frameworks").([]interface{}))

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
			myResources, _ = CastToStringList(temp["resource_types"].([]interface{}))
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
		myItem.Cisawsv12, _ = CastToStringList(s["cis_aws_v12"].([]interface{}))
		myItem.Cisawsv13, _ = CastToStringList(s["cis_aws_v13"].([]interface{}))
		myItem.Cisazurev11, _ = CastToStringList(s["cis_azure_v11"].([]interface{}))
		myItem.Cisazurev12, _ = CastToStringList(s["cis_azure_v12"].([]interface{}))
		myItem.Cisazurev13, _ = CastToStringList(s["cis_azure_v13"].([]interface{}))
		myItem.Cisgcpv11, _ = CastToStringList(s["cis_gcp_v11"].([]interface{}))
		myItem.Ciskubernetesv15, _ = CastToStringList(s["cis_kubernetes_v15"].([]interface{}))
		myItem.Ciskubernetesv16, _ = CastToStringList(s["cis_kubernetes_v16"].([]interface{}))
		myItem.Cisdockerv11, _ = CastToStringList(s["cis_docker_v11"].([]interface{}))
		myItem.Ciseksv11, _ = CastToStringList(s["cis_eks_v11"].([]interface{}))
		myItem.Cisgkev11, _ = CastToStringList(s["cis_gke_v11"].([]interface{}))
		return myItem, nil
	}

	return myItem, errors.New("no benchmark data")
}

//goland:noinspection GoUnusedParameter
func resourceSimplePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	diags = setNotNil(typedjson, d, diags, "provider", "cloud_provider")
	diags = setNotNil(typedjson, d, diags, "file", "file")
	diags = setNotNil(typedjson, d, diags, "title", "title")
	diags = setNotNil(typedjson, d, diags, "severity", "severity")
	diags = setNotNil(typedjson, d, diags, "category", "category")

	if typedjson["frameworks"] != nil {
		err = d.Set("frameworks", typedjson["frameworks"])
		diags = LogAppendError(err, diags)
	}

	if typedjson["guideline"] != nil {
		err = d.Set("guidelines", typedjson["guideline"])
		diags = LogAppendError(err, diags)
	}

	//myconditions should be an array it currently a map
	//hence this fudge
	//todo: once you start passing around condition arrays
	//this can go
	myConditions := make([]interface{}, 1)
	myConditions[0] = typedjson["conditionQuery"]
	err = d.Set("conditions", myConditions)
	diags = LogAppendError(err, diags)

	return diags
}

func resourceSimplePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policyID := d.Id()
	configure := m.(ProviderConfig)
	params := RequestParams{"%s/policies/" + policyID, "v1", "DELETE"}
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
