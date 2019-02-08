package jsmapper

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestPrepareMapper(t *testing.T) {
	mapperExpr := `{
		"port": {{$env[HTTPPORT]}},
		"event": {{$flow.event}},
		"street": {{$flow.content.address}}.street,
		"message": {{$activity[app_2].message}},
		"vehicle": {{$activity[app_16].value}}.vehicle.{
			"make": make,
			"year": year,
			"vin": vin.value
		}
	}`
	expr, attrs := prepareMapper(mapperExpr)
	fmt.Println(expr)
	fmt.Printf("data tags: %+v\n", attrs)
}

func TestTransform(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivMapexpr, "jsonata")

	act.Eval(tc)

	value := tc.GetOutput(ovValue).(string)

	assert.NotNil(t, value, "output value should be set")
	assert.Equal(t, "jsonata", value, "not equal")
}
