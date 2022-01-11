package bridgecrew

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

//RequestParams (parameters) for auth client
type RequestParams struct {
	path    string
	version string
	method  string
}

//use basic auth client
func authClient(params RequestParams, configure ProviderConfig) (*http.Client, *http.Request, diag.Diagnostics, bool, error) {

	var diags diag.Diagnostics
	api := configure.Token
	url := configure.URL
	prisma := configure.Prisma
	accessKey := configure.AccessKeyID
	secretKey := configure.SecretKey

	if prisma != "" {
		//check accessKey and secretKey aren't empty
		if accessKey == "" {
			log.Fatal("PRISMA_ACCESS_KEY_ID is missing")
		}

		if secretKey == "" {
			log.Fatal("PRISMA_SECRET_KEY is missing")
		}

		loginURL := prisma + "/login"

		var err diag.Diagnostics
		api, err = loginPrisma(accessKey, secretKey, loginURL)
		url = prisma + "/bridgecrew"

		if err != nil {
			diags = append(diags, err[0])
		}
	}

	baseurl := fmt.Sprintf(url + "/api/" + params.version)

	if api == "" {
		log.Fatal("BRIDGECREW_API is missing")
	}

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + api

	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest(params.method, fmt.Sprintf(params.path, baseurl), nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Get request failed %s \n", err.Error()),
		})
	}

	// add authorization header to the req
	req.Header.Add("authorization", bearer)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	return client, req, diags, false, err
}

// CheckStatus confirms returns codes are 200
func CheckStatus(res *http.Response) (diag.Diagnostics, bool) {
	if res.StatusCode != http.StatusOK {
		errStr := fmt.Errorf("Non-OK HTTP status: %d", res.StatusCode)
		return diag.FromErr(errStr), true
	}

	return nil, false
}

func loginPrisma(username string, password string, loginURL string) (string, diag.Diagnostics) {
	payload := strings.NewReader(fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\"}", username, password))
	req, _ := http.NewRequest("POST", loginURL, payload)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	rawToken, _ := ioutil.ReadAll(res.Body)

	mySecrets := make(map[string]interface{})
	err := json.Unmarshal(rawToken, &mySecrets)

	if err != nil {
		return "", diag.FromErr(err)
	}

	if mySecrets["message"] != "login_successful" {
		errStr := fmt.Errorf(mySecrets["message"].(string))
		return "", diag.FromErr(errStr)
	}

	return mySecrets["token"].(string), nil
}
