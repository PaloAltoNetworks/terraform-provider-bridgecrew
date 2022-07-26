package bridgecrew

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/karlseguin/typed"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
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
			"frameworks": {
				Type:        schema.TypeList,
				Description: "Which IAC framework is this policy targeting.",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
			"file": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "This is the name of the YAML policy file.",
				ValidateFunc: ValidateIsYAMLFile,
			},
			"source_code_hash": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "By providing the source code hash change to the YAML file can be caught and the resource updated.",
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

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	myPolicy, err := setPolicy(d)

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

	newResults, d2, done := VerifyReturn(body)
	if done {
		return d2
	}

	d.SetId(newResults.Policy)

	resourcePolicyRead(ctx, d, m)

	return diags
}

func setPolicy(d *schema.ResourceData) (Policy, error) {
	myPolicy := Policy{}
	myBenchmark, err := setBenchmark(d)

	if err == nil {
		myPolicy.Benchmarks = myBenchmark
	}

	filename, hasFilename := d.GetOk("file")

	// if the filename is set then this is a yaml policy
	if hasFilename {
		file, _ := filepath.Abs(filename.(string))
		code, err := loadFileContent(file)
		if err != nil {
			return myPolicy, fmt.Errorf("unable to load %q: %w", filename.(string), err)
		}
		myPolicy.Code = string(code)
	} else {
		return myPolicy, fmt.Errorf("filename not set")
	}

	myPolicy.Provider = d.Get("cloud_provider").(string)
	myPolicy.Frameworks, _ = CastToStringList(d.Get("frameworks").([]interface{}))

	return myPolicy, nil
}

//goland:noinspection GoUnusedParameter
func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	if typedjson["frameworks"] != nil {
		err = d.Set("frameworks", typedjson["frameworks"])
		diags = LogAppendError(err, diags)
	}

	return diags
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	policyID := d.Id()
	if policyChange(d) {
		myPolicy, err := setPolicy(d)

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

		_, d2, done := VerifyReturn(body)

		if done {
			return d2
		}

		err = d.Set("last_updated", time.Now().Format(time.RFC850))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourcePolicyRead(ctx, d, m)
}

func policyChange(d *schema.ResourceData) bool {
	return d.HasChange("cloud_provider") ||
		d.HasChange("benchmarks") ||
		d.HasChange("source_code_hash") ||
		d.HasChange("file") ||
		d.HasChange("frameworks")
}

//goland:noinspection GoUnusedParameter
func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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
