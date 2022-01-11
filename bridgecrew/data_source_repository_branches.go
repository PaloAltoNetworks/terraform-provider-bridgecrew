package bridgecrew

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRepositoryBranches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRepositoryBranchRead,
		Schema: map[string]*schema.Schema{
			"repositoriesbranches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Resource{},
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceRepositoryBranchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	target := d.Get("target")
	params := RequestParams{"%s/repositories/branches" + "/" + target.(string), "v1", "GET"}

	//todo endpoint doesnt work like this

	configure := m.(ProviderConfig)
	client, req, diagnostics, done, err := authClient(params, configure)

	if done {
		return diagnostics
	}

	r, err := client.Do(req)

	if err != nil {
		log.Print("Failed at client.Do")
		return diag.FromErr(err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	log.Print("All data obtained")
	repositoriesbranches := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&repositoriesbranches)

	//todo this actually needs the target repository

	log.Print("Decoded data")
	log.Print(r.Body)

	if err != nil {
		log.Fatal("Failed to parse data")
	}

	log.Print(repositoriesbranches)
	flatBranch := flattenBranchData(&repositoriesbranches)

	if err := d.Set("repositories", flatBranch); err != nil {
		log.Fatal(reflect.TypeOf(repositoriesbranches))
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diagnostics
}

func flattenBranchData(Repositories *[]map[string]interface{}) []interface{} {
	if Repositories != nil {
		ois := make([]interface{}, len(*Repositories))

		//for i, Repository := range *Repositories {
		//	oi := make(map[string]interface{})
		//oi["name"] = Repository["name"]
		//oi["creationdate"] = Repository["creationDate"]
		//oi["defaultbranch"] = Repository["defaultBranch"]

		//	ois[i] = oi
		//}

		return ois
	}

	return make([]interface{}, 0)
}
