package bridgecrew

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIncidentsInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncidentInfoRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"open": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"passed": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"suppressed": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"traced": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"traced": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"nottraced": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"encryption": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encrypted": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"unencrypted": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"noencryption": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"reachability": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"private": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"noreachability": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceIncidentInfoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := RequestParams{"%s/incidents/info", "v2", "POST"}

	configure := m.(ProviderConfig)
	client, req, diagnostics, done, err := authClient(params, configure)

	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed at authClient %s \n", err.Error()),
		})
	}

	if done {
		return diagnostics
	}

	r, err := client.Do(req)

	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed at client.Do %s \n", err.Error()),
		})
		return diagnostics
	}

	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	myinfo := make(map[string]interface{})
	err = json.Unmarshal(body, &myinfo)

	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to parse data %s \n", err.Error()),
		})

		return diagnostics
	}

	test := myinfo["data"].(map[string]interface{})
	d.Set("status", flattenStatus(test["status"].(map[string]interface{})))
	d.Set("traced", flattenTraced(test["traced"].(map[string]interface{})))
	d.Set("encryption", flattenEncryption(test["encryption"].(map[string]interface{})))
	d.Set("reachability", flattenReachability(test["reachability"].(map[string]interface{})))
	d.Set("total", test["total"])

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenEncryption(configuration map[string]interface{}) []interface{} {
	if configuration == nil {
		return []interface{}{}
	}

	return []interface{}{map[string]interface{}{
		"encrypted":    configuration["encrypted"],
		"unencrypted":  configuration["unencrypted"],
		"noencryption": configuration["noEncryption"],
	}}
}

func flattenReachability(configuration map[string]interface{}) []interface{} {
	if configuration == nil {
		return []interface{}{}
	}

	return []interface{}{map[string]interface{}{
		"public":         configuration["public"],
		"private":        configuration["private"],
		"noreachability": configuration["noReachability"],
	}}
}

func flattenStatus(configuration map[string]interface{}) []interface{} {
	if configuration == nil {
		return []interface{}{}
	}

	return []interface{}{map[string]interface{}{
		"open":       configuration["open"],
		"passed":     configuration["passed"],
		"suppressed": configuration["suppressed"],
	}}
}

func flattenTraced(configuration map[string]interface{}) []interface{} {
	if configuration == nil {
		return []interface{}{}
	}

	return []interface{}{map[string]interface{}{
		"traced":    configuration["traced"],
		"nottraced": configuration["notTraced"],
	}}
}
