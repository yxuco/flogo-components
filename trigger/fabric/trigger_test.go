package fabric

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/stretchr/testify/assert"
)

func getJSONMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
  "id": "fabric-trigger",
  "settings": {
    "setting": "somevalue"
  },
  "handlers": [
    {
      "settings": {
        "handler_setting": "somevalue"
      },
      "action" {
	     "id": "test_action"
      }
    }
  ]
}`

func TestCreate(t *testing.T) {

	// New factory
	md := trigger.NewMetadata(getJSONMetadata())
	f := NewFactory(md)

	if f == nil {
		t.Fail()
	}

	// New Trigger
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	trg := f.New(&config)

	if trg == nil {
		t.Fail()
	}
}

func TestUnmarshal(t *testing.T) {
	s := `[{"n1":"v1"}, {"n2":"v2"}]`
	result := unmarshalString(s)
	array, ok := result.([]interface{})
	assert.True(t, ok, "result should be array of map[string]interface{}")
	assert.Equal(t, 2, len(array), "result array length should be 2")

	s = `["n1", "n2"]`
	result = unmarshalString(s)
	sArray, ok := result.([]interface{})
	assert.True(t, ok, "result should be array of interface{}")
	assert.Equal(t, 2, len(sArray), "result array length should be 2")

	s = `["n1", {"n2":"v2"}]`
	result = unmarshalString(s)
	sArray, ok = result.([]interface{})
	assert.True(t, ok, "result should be array of interface{}")
	assert.Equal(t, 2, len(sArray), "result array length should be 2")
	obj, ok := sArray[1].(map[string]interface{})
	assert.True(t, ok, "second array item should be of type may[string]interface{}")
	assert.Equal(t, "v2", obj["n2"].(string), "obj should contain value 'v2'")
}
