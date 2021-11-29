package bridgecrew

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"gopkg.in/yaml.v2"
)

// CheckYAMLString Takes a value containing YAML string and passes it through
// the YAML parser. Returns either a parsing
// error or original YAML string.
func CheckYAMLString(yamlString interface{}) (string, error) {
	var y interface{}

	if yamlString == nil || yamlString.(string) == "" {
		return "", nil
	}

	s := yamlString.(string)

	err := yaml.Unmarshal([]byte(s), &y)

	return s, err
}

// VerifyReturn Looks at the return object from the Platform
func VerifyReturn(err error, body []byte) (*Result, diag.Diagnostics, bool) {
	newResults := &Result{}
	err = json.Unmarshal([]byte(body), newResults)

	if err != nil {
		errStr := errors.New("platform failed to return ID")
		return nil, diag.FromErr(errStr), true
	}
	return newResults, nil, false
}

// CastToStringList is a helper to work with conversion of types
// If there's a better way (most likely)?
func CastToStringList(temp []interface{}) []string {
	var versions []string
	for _, version := range temp {
		versions = append(versions, version.(string))
	}
	return versions
}

// highlight is just to help with manual debugging, so you can find the lines
//goland:noinspection SpellCheckingInspection
func highlight(myPolicy interface{}) {
	log.Print("XXXXXXXXXXX")
	log.Print(myPolicy)
	log.Print("XXXXXXXXXXX")
}
