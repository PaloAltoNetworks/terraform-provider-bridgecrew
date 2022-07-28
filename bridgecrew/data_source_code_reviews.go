package bridgecrew

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCodeReviews() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeReviewsRead,
		Schema: map[string]*schema.Schema{
			"codereviews": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"commit_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"run_id": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"git_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scan_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lastscandate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creationdate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"results": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"critical": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"high": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"medium": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"low": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
								},
							},
						},
						"pr": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"title": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pr_number": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"from_branch": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"into_branch": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enforcement_rule": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"supplychain": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"softfailthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"hardfailthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"commentsbotthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"secrets": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"softfailthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"hardfailthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"commentsbotthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"iac": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"softfailthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"hardfailthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"commentsbotthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"images": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"softfailthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"hardfailthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"commentsbotthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"opensource": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"softfailthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"hardfailthreshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"commentsbotthreshold": {
																Type:     schema.TypeString,
																Computed: true,
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
				},
			},
		},
	}
}

//goland:noinspection GoUnusedParameter,GoLinter,GoLinter
func dataSourceCodeReviewsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	request := "%s/development-pipeline/code-review/runs/data"

	params := RequestParams{request, "v1", "GET"}

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

	Reviews := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&Reviews)

	if err != nil {
		if err.Error() == "EOF" {
			temp := "no data found"
			err = errors.New(temp)
			log.Print(temp)
		} else {
			log.Println("Failed to parse data")
		}
		return diag.FromErr(err)
	}

	if err := flattenCodeReviews(Reviews, d); err != nil {
		return err
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diagnostics
}

func flattenCodeReviews(codereviews map[string]interface{}, d *schema.ResourceData) diag.Diagnostics {

	if codereviews != nil {
		data := codereviews["data"].([]interface{})
		reviews := make([]interface{}, 0)
		for _, item := range data {
			review := item.(map[string]interface{})
			myreview := make(map[string]interface{})
			myreview["id"] = review["id"].(float64)
			myreview["commit_id"] = review["commitId"].(string)
			myreview["git_user"] = review["gitUser"].(string)
			myreview["scan_status"] = review["scanStatus"].(string)
			myreview["customer_name"] = review["customerName"].(string)
			myreview["lastscandate"] = review["lastScanDate"].(string)
			myreview["run_id"] = review["runId"].(float64)
			myreview["creationdate"] = review["creationDate"].(string)
			myreview["repo_id"] = review["repo_id"].(string)

			if review["PR"] != nil {
				PR := review["PR"].(map[string]interface{})
				myprs := make([]interface{}, 0)
				mypr := make(map[string]interface{})
				mypr["from_branch"] = PR["fromBranch"].(string)
				mypr["into_branch"] = PR["intoBranch"].(string)
				mypr["pr_number"] = PR["prNumber"].(string)

				enforcement := make([]interface{}, 0)
				mycat := make(map[string]interface{})

				if PR["enforcementRule"] != nil {
					mycode := PR["enforcementRule"].(map[string]interface{})

					supplies := setcategories(mycode, "SUPPLY_CHAIN")
					secrets := setcategories(mycode, "SECRETS")
					iac := setcategories(mycode, "IAC")
					images := setcategories(mycode, "IMAGES")
					opensource := setcategories(mycode, "OPEN_SOURCE")

					mycat["supplychain"] = supplies
					mycat["secrets"] = secrets
					mycat["iac"] = iac
					mycat["images"] = images
					mycat["opensource"] = opensource

					enforcement = append(enforcement, mycat)
					mypr["enforcement_rule"] = enforcement
				}

				myprs = append(myprs, mypr)
				myreview["pr"] = myprs
			}

			if review["results"] != nil {
				results := review["results"].(map[string]interface{})
				myresults := make([]interface{}, 0)
				myresult := make(map[string]interface{})
				if results["CRITICAL"] != nil {
					myresult["critical"] = results["CRITICAL"].(float64)
				}
				if results["HIGH"] != nil {
					myresult["high"] = results["HIGH"].(float64)
				}
				if results["MEDIUM"] != nil {
					myresult["medium"] = results["MEDIUM"].(float64)
				}
				if results["LOW"] != nil {
					myresult["low"] = results["LOW"].(float64)
				}
				myresults = append(myresults, myresult)

				myreview["results"] = myresults
			}
			myreview["status"] = review["status"].(string)
			myreview["repository"] = review["repository"].(string)
			myreview["source_type"] = review["sourceType"].(string)
			myreview["organization"] = review["organization"].(string)
			reviews = append(reviews, myreview)
		}

		if err := d.Set("codereviews", reviews); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
