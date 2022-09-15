package bridgecrew

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/karlseguin/typed"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	//nolint
	err := yaml.Unmarshal([]byte(s), &y)

	return s, err
}

// VerifyReturn Looks at the return object from the Platform
func VerifyReturn(body []byte) (*Result, diag.Diagnostics, bool) {
	newResults := &Result{}
	err := json.Unmarshal(body, newResults)

	if err != nil {
		errStr := errors.New("platform failed to return ID")
		return nil, diag.FromErr(errStr), true
	}
	return newResults, nil, false
}

// CastToStringList is a helper to work with conversion of types
// If there's a better way (most likely)?
func CastToStringList(temp []interface{}) ([]string, bool) {
	var versions []string
	if temp != nil {
		for _, version := range temp {
			if version != nil {
				versions = append(versions, version.(string))
			}
		}
	} else {
		log.Print("Cast from Nil")
		return nil, true
	}
	return versions, false
}

// highlight is just to help with manual debugging, so you can find the lines
//
//goland:noinspection SpellCheckingInspection
func highlight(myPolicy interface{}) {
	log.Print("XXXXXXXXXXX")
	log.Print(myPolicy)
	log.Print("XXXXXXXXXXX")
}

func setNotNil(typedjson typed.Typed, d *schema.ResourceData, diags diag.Diagnostics, item string, toset string) diag.Diagnostics {
	var err error
	if typedjson[item] != nil {
		err = d.Set(toset, strings.ToLower(typedjson[item].(string)))
	}
	diags = LogAppendError(err, diags)
	return diags
}
