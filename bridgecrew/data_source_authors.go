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

func dataSourceAuthors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthorsRead,
		Schema: map[string]*schema.Schema{
			"fullreponame": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full Repo Name",
			},
			"sourcetype": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Repository Source Type",
				ValidateFunc: ValidateRepository,
			},
			"authors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// goland:noinspection GoUnusedParameter
func dataSourceAuthorsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	fullreponame := d.Get("fullreponame").(string)
	sourceType := d.Get("sourcetype").(string)
	request := "%s/errors/gitBlameAuthors?fullRepoName=" + fullreponame + "&sourceType=" + sourceType
	params := RequestParams{request, "v1", "GET"}

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

	nubody, err := ioutil.ReadAll(r.Body)

	temp := strings.ReplaceAll(string(nubody), "\"", "")
	temp = strings.ReplaceAll(strings.Replace(temp, "]", "", -1), "[", "")

	Authors := strings.Split(temp, ",")

	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to parse data %s \n", err.Error()),
		})

		return diagnostics
	}

	if err := d.Set("authors", Authors); err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Set Errors failed %s \n", err.Error()),
		})
		return diagnostics
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}
