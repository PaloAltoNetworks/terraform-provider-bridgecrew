package bridgecrew

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/karlseguin/typed"
)

func resourceTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagCreate,
		ReadContext:   resourceTagRead,
		UpdateContext: resourceTagUpdate,
		DeleteContext: resourceTagDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tag rule description",
			},
			"definition": {
				MaxItems:    1,
				Required:    true,
				Type:        schema.TypeList,
				Description: "Tag Definition",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_groups": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Tag group",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag group name",
									},
									"tags": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeMap,
													Required: true,
													Elem: &schema.Schema{
														Type:     schema.TypeString,
														Required: true,
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
			},
			"isenabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Tag is enabled",
			},
			"tagruleootbid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "tagRuleOOTBId",
			},
			"createdby": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Record created by",
			},
			"candoactions": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "candoactions",
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
			"repositories": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "A list of repository IDs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTagCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	myTag, err := setTag(d)

	if err != nil {
		return diag.FromErr(err)
	}

	jsPolicy, err := json.Marshal(myTag)
	if err != nil {
		return diag.FromErr(err)
	}

	payload := strings.NewReader(string(jsPolicy))
	params := RequestParams{"%s/tag-rules", "v1", "POST"}
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

	resourceTagRead(ctx, d, m)

	return diags
}

func setTag(d *schema.ResourceData) (Tag, error) {
	myTag := Tag{}
	myTag.Name = d.Get("name").(string)
	temp := d.Get("definition").([]interface{})

	definition := temp[0].(map[string]interface{})
	groups := definition["tag_groups"].([]interface{})
	myGroup := TagGroup{}

	for _, group := range groups {
		test := group.(map[string]interface{})
		myGroup.Name = test["name"].(string)
		rawvalues := test["tags"].([]interface{})
		var myTags []Tags
		for _, rawvalue := range rawvalues {
			myTagged := Tags{}
			value := rawvalue.(map[string]interface{})
			myTagged.Name = value["name"].(string)
			myTagged.Value = value["value"].(map[string]interface{})
			myTags = append(myTags, myTagged)
		}

		myGroup.Tags = myTags
	}

	myTag.Definition.TagGroups = append(myTag.Definition.TagGroups, myGroup)
	repos := d.Get("repositories").([]interface{})

	for _, repo := range repos {
		myTag.Repositories = append(myTag.Repositories, repo.(string))
	}

	myTag.Description = d.Get("description").(string)
	myTag.IsEnabled = d.Get("isenabled").(bool)
	myTag.TagRuleOOTBId = d.Get("tagruleootbid").(string)
	return myTag, nil
}

//goland:noinspection GoUnusedParameter
func resourceTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tagID := d.Id()

	params := RequestParams{"%s/tag-rules/" + tagID, "v1", "GET"}
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

	diags = setNotNil(typedjson, d, diags, "name", "name")
	diags = setNotNil(typedjson, d, diags, "description", "description")
	diags = setNotNil(typedjson, d, diags, "tagRuleOOTBid", "tagruleootbid")
	diags = setNotNil(typedjson, d, diags, "createdBy", "createdby")
	diags = setNotNil(typedjson, d, diags, "creationDate", "creationdate")

	if typedjson["isEnabled"] != nil {
		err = d.Set("isenabled", typedjson["isEnabled"].(bool))
		diags = LogAppendError(err, diags)
	}

	if typedjson["canDoActions"] != nil {
		err = d.Set("candoactions", typedjson["canDoActions"].(bool))
		diags = LogAppendError(err, diags)
	}

	if reflect.TypeOf(typedjson["definition"]).String() != "[]interface {}" {
		var temp []interface{}
		definition := typedjson["definition"]
		temp = append(temp, definition)
		err = d.Set("definition", temp)
	} else {
		err = d.Set("definition", typedjson["definition"])
	}
	diags = LogAppendError(err, diags)

	if reflect.TypeOf(typedjson["repositories"]).String() != "[]interface {}" {
		diags = setNotNil(typedjson, d, diags, "repositories", "repositories")
	} else {
		var repositories []string
		temp := typedjson["repositories"].([]interface{})
		for _, repo := range temp {
			temp := repo.(map[string]interface{})
			repositories = append(repositories, temp["id"].(string))
		}
		err = d.Set("repositories", repositories)
		diags = LogAppendError(err, diags)
	}

	return diags
}

func resourceTagUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	tagID := d.Id()
	if tagChange(d) {
		myTag, err := setTag(d)

		if err != nil {
			return diag.FromErr(err)
		}

		jsPolicy, err := json.Marshal(myTag)
		if err != nil {
			return diag.FromErr(err)
		}

		payload := strings.NewReader(string(jsPolicy))

		params := RequestParams{"%s/tag-rules/" + tagID, "v1", "PUT"}
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

		//no such field
		err = d.Set("last_updated", time.Now().Format(time.RFC850))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceTagRead(ctx, d, m)
}

func tagChange(d *schema.ResourceData) bool {
	return d.HasChange("definition") ||
		d.HasChange("repositories") ||
		d.HasChange("isenabled") ||
		d.HasChange("tagruleootbid") ||
		d.HasChange("name") ||
		d.HasChange("description")
}

//goland:noinspection GoUnusedParameter
func resourceTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	tagID := d.Id()
	configure := m.(ProviderConfig)
	params := RequestParams{"%s/tag-rules/tagRuleId" + tagID, "v1", "DELETE"}
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
