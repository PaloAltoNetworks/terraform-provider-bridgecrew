package bridgecrew

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrganisation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganisationRead,

		Schema: map[string]*schema.Schema{
			"organisation": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

//goland:noinspection GoUnusedParameter
func dataSourceOrganisationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := RequestParams{"%s/organization", "v1", "GET"}

	configure := m.(ProviderConfig)
	client, req, diagnostics, done := authClient(params, configure, nil)

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

	myinfo := strings.Trim(string(body), "\"")

	check := d.Set("organisation", myinfo)
	diagnostics = LogAppendError(check, diagnostics)

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}
