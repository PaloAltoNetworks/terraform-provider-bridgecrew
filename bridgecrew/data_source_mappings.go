package bridgecrew

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMappings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMappingsRead,
		Schema: map[string]*schema.Schema{
			"guidelines": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check": {
							Type:     schema.TypeString,
							Required: true,
						},
						"guideline": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"idmappings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bcid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"checkovid": {
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
func dataSourceMappingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	params := RequestParams{"%s/guidelines", "v1", "GET"}

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

	body, _ := ioutil.ReadAll(r.Body)
	Mappings := make(map[string]interface{})

	err = json.Unmarshal(body, &Mappings)

	if err != nil {
		log.Print("Failed to parse data")
		return LogAppendError(err, diagnostics)
	}

	err = setGuidelines(Mappings, d)

	if err != nil {
		diagnostics = LogAppendError(err, diagnostics)
	}

	err = setLookups(Mappings, d)

	if err != nil {
		diagnostics = LogAppendError(err, diagnostics)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func setLookups(mappings map[string]interface{}, d *schema.ResourceData) error {
	var lookups = mappings["idMapping"].(map[string]interface{})
	flatLookups := make([]interface{}, 0)
	for y := range lookups {
		oi := make(map[string]interface{})
		oi["checkovid"] = lookups[y].(string)
		oi["bcid"] = y
		flatLookups = append(flatLookups, oi)
	}

	if err := d.Set("idmappings", flatLookups); err != nil {
		return err
	}

	return nil
}

func setGuidelines(mappings map[string]interface{}, d *schema.ResourceData) error {
	var guidelines = mappings["guidelines"].(map[string]interface{})
	flatMaps := make([]interface{}, 0)
	for i := range guidelines {
		oi := make(map[string]interface{})
		oi["check"] = i
		oi["guideline"] = guidelines[i].(string)
		flatMaps = append(flatMaps, oi)
	}

	if err := d.Set("guidelines", flatMaps); err != nil {
		return err
	}
	return nil
}
