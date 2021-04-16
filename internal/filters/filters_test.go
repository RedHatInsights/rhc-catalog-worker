package filters

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapFilterBadValue(t *testing.T) {
	f := Value{Data: "Not a valid expression"}
	body := `{"id": 100, "name": "Fred Flintstone", "age": 56, "state": "NY"}`
	var jsonBody map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(body)))
	decoder.UseNumber()
	err := decoder.Decode(&jsonBody)
	if err != nil {
		t.Error("Error decoding JSON")
	}

	_, err = f.Apply(jsonBody)
	if err == nil {
		t.Error("Parsing did not fail")
	}
}

func TestMapFilter(t *testing.T) {
	f := Value{Data: `{"id":"id", "name":"name"}`}
	body := `{"id": 100, "name": "Fred Flintstone", "age": 56, "state": "NY"}`
	var jsonBody map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(body)))
	decoder.UseNumber()
	err := decoder.Decode(&jsonBody)
	if err != nil {
		t.Error("Error decoding JSON")
	}
	result, err := f.Apply(jsonBody)
	if err != nil {
		t.Error(err)
	}
	if len(result) != 2 {
		t.Error("Key Length didn't match 2")
	}

	if result["name"].(string) != "Fred Flintstone" {
		t.Error("Name didn't match Fred Flintstone")
	}

	value, err := result["id"].(json.Number).Int64()
	assert.NoError(t, err)
	if value != 100 {
		t.Error("id didn't match 100")
	}

	if _, found := result["age"]; found {
		t.Error("age should be missing")
	}
}

func TestMapWithArrayFilter(t *testing.T) {
	f := Value{Data: "results[].{catalog_id:id, name:name}", ReplaceResults: true}
	body := `{"count": 2, "results":[{"id": 100, "name": "Fred", "age": 56}, {"id": 200, "name": "Barney", "state": "NY"}]}`
	var jsonBody map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(body)))
	decoder.UseNumber()
	err := decoder.Decode(&jsonBody)
	assert.NoError(t, err)

	result, err := f.Apply(jsonBody)
	assert.NoError(t, err)

	value, err := result["count"].(json.Number).Int64()
	assert.NoError(t, err)
	if value != 2 {
		t.Error("count didn't match 2")
	}

	x := result["results"].([]interface{})
	item := x[0].(map[string]interface{})
	if item["name"].(string) != "Fred" {
		t.Error("Name didn't match Fred")
	}

}

func TestFilterStringValue(t *testing.T) {
	f := Value{}
	f.Parse("results[].{catalog_id:id, url:url,created:created,name:name, modified:modified, playbook:playbook}")
	if !f.ReplaceResults {
		t.Error("Results should be replaced")
	}
}

func TestFilterMapStringValue(t *testing.T) {
	f := Value{}
	v := map[string]interface{}{"id": "id", "url": "url", "description": "description", "name": "name", "playbook": "playbook"}
	f.Parse(v)
	if f.ReplaceResults {
		t.Error("Results should not be replaced")
	}
}
