package bridgecrew

import (
	"reflect"
	"testing"
)

func TestCheckYAMLString(t *testing.T) {
	var err error
	var actual string

	validYaml := `---
abc:
  def: 123
  xyz:
    -
      a: "ホリネズミ"
      b: "1"
`

	actual, err = CheckYAMLString(validYaml)
	if err != nil {
		t.Fatalf("Expected not to throw an error while parsing YAML, but got: %s", err)
	}

	// We expect the same YAML string back
	if actual != validYaml {
		t.Fatalf("Got:\n\n%s\n\nExpected:\n\n%s\n", actual, validYaml)
	}

	invalidYaml := `abc: [`

	actual, err = CheckYAMLString(invalidYaml)
	if err == nil {
		t.Fatalf("Expected to throw an error while parsing YAML, but got: %s", err)
	}

	// We expect the invalid YAML to be shown back to us again.
	if actual != invalidYaml {
		t.Fatalf("Got:\n\n%s\n\nExpected:\n\n%s\n", actual, invalidYaml)
	}
}

func TestHighlight(t *testing.T) {
	highlight("For Coverage")
}

func TestCastToStringList(t *testing.T) {

	expected := []string{"first", "second", "third"}
	names := make([]interface{}, len(expected))
	for i, s := range expected {
		names[i] = s
	}

	actual, _ := CastToStringList(names)

	// We expect the same YAML string back
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Got:\n\n%s\n\nExpected:\n\n%s\n", actual, expected)
	}
}
