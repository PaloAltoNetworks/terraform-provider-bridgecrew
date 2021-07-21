package bridgecrew

import (
	"context"
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"
	"time"
	"log"
	"os"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRepositoryRead,
		Schema: map[string]*schema.Schema{
			"repositories": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"repository": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"defaultbranch": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ispublic": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creationdate": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

    api := os.Getenv("BRIDGECREW_API")
    
    if api == "" {
        log.Fatal("BRIDGECREW_API is missing")
    }

    // Create a Bearer string by appending string access token
    var bearer = "Bearer " + api

	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repositories", "https://www.bridgecrew.cloud/api/v1"), nil)

	if err != nil {
		log.Fatal("Failed at http")
		return diag.FromErr(err)
	}

	log.Print("Passed http Request")

    // add authorization header to the req
    req.Header.Add("Authorization", bearer)
	
	log.Print("Added Header")
	
	r, err := client.Do(req)
	
	log.Print("Queried")

	if err != nil {
		log.Fatal("Failed at client.Do")
		return diag.FromErr(err)
	}
	defer r.Body.Close()
	
	log.Print("All data obtained")
	repositories := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&repositories)
	
	log.Print("Decoded data")

	if err != nil {
		log.Fatal("Failed to parse data")
		return diag.FromErr(err)
	}
    
	log.Print(repositories)

	if err := d.Set("repositories", repositories); err != nil {
		log.Fatal(reflect.TypeOf(repositories))
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
