/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package gojsonata

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngineInit(t *testing.T) {
	modules := CachedModules()
	fmt.Printf("cached modules %+v\n", modules)
	assert.Equal(t, 7, len(modules), "number of initial cache should contain 7 modules")
	assert.Contains(t, modules, "traceur-runtime", "traceur-runtime should be one of the cached modules")
	assert.Contains(t, modules, "jsonata", "jsonata should be one of the cached modules")
	assert.Contains(t, modules, "parser", "parser should be one of the cached modules")
	assert.Contains(t, modules, "functions", "functions should be one of the cached modules")
	assert.Contains(t, modules, "signature", "signature should be one of the cached modules")
	assert.Contains(t, modules, "utils", "utils should be one of the cached modules")
	assert.Contains(t, modules, "datetime", "datetime should be one of the cached modules")
}

func TestSimpleExpr(t *testing.T) {

	value, err := RunScript(
		`var data = {
		  example: [
		    {value: 4},
		    {value: 7},
		    {value: 13}
		  ]
	    };
	  
	    var expression = jsonata('$sum(example.value)');
	    expression.evaluate(data);`)
	require.NoError(t, err, "simple expression should not result in error")
	result, err := value.ToInteger()
	require.NoError(t, err, "result should be an integer")
	assert.Equal(t, int64(24), result, "sum of values should be 24")
}

func TestCallScript(t *testing.T) {
	result, err := CallScriptFile("./vehicle.js", "result")
	require.NoError(t, err, "load vehicle.js file should not be error")
	assert.Equal(t, vehicleResult, result.String(), "result should match content of vehicleResult")
}

func TestTransform(t *testing.T) {
	source := make(map[string]interface{})
	err := json.Unmarshal([]byte(vehicleData), &source)
	require.NoError(t, err, "unmarshal vehicle data should not be error")
	result, err := Transform(source, vehicleExpr)
	require.NoError(t, err, "call transform function should not be error")
	data, err := json.Marshal(result)
	require.NoError(t, err, "json marshal should not be error")
	fmt.Printf("result value: %s\n", string(data))
	assert.Equal(t, vehicleResult, string(data), "result should match content of vehicleResult")
}

func TestTransformArray(t *testing.T) {
	source := make(map[string]interface{})
	err := json.Unmarshal([]byte(vehicleData), &source)
	require.NoError(t, err, "unmarshal vehicle data should not be error")
	result, err := Transform(source, vehicleExpr)
	require.NoError(t, err, "call transform function should not be error")

	array := result.([]map[string]interface{})
	fmt.Printf("result array type %T, first item: %+v\n", result, array[0])
	assert.Equal(t, "CHEVROLET", array[0]["make"].(string), "first vehicle make should be 'CHEVROLET'")
}

var vehicleData = `{
    "_msgid": "146e84d7.4e38eb",
    "topic": "",
    "payload": {
        "@type": "PagedVehicleResult",
        "index": 0,
        "limit": 200,
        "count": 175,
        "total": 175,
        "vehicle": [
        {
            "@id": 942297,
            "vin": {
                "@readOnly": true,
                "value": "1GNSK3EC1FRXXXXXX"
            },
            "label": "E M660",
            "color": "WHITE",
            "make": "CHEVROLET",
            "model": "TAHOE",
            "deviceSerialNumber": {
                "@readOnly": true,
                "value": "50X00XXX8X"
            },
            "year": 2015,
            "odometer": {
                "@units": "MILES",
                "@timestamp": "2018-07-16T11:51:17Z",
                "@readOnly": "TRUE",
                "value": 26775.18
            },
            "engineRunTime": {
                "@isTracked": "NO",
                "@readOnly": "TRUE",
                "value": "PT0.000S"
            },
            "licensePlate": {
                "@state": "NY",
                "@country": "US",
                "value": "AXXXXX"
            },
            "trackableItemType": "VEHICLE",
            "fuelType": "GASOLINE",
            "createdTimestamp": {
                "@readOnly": true,
                "value": "2018-06-28T15:05:12Z"
            },
            "modifiedTimestamp": {
                "@readOnly": true,
                "value": "2018-06-28T15:05:12Z"
            }
        },
        {
            "@id": 763803,
            "vin": {
                "@readOnly": true,
                "value": "1FTFX1EF9EFXXXXXX"
            },
            "label": "LW04",
            "color": "WHITE",
            "make": "FORD",
            "model": "F-150",
            "deviceSerialNumber": {
                "@readOnly": true,
                "value": "5X147XX5XX"
            },
            "year": 2014,
            "odometer": {
                "@units": "MILES",
                "@timestamp": "2018-11-02T14:57:56Z",
                "@readOnly": "TRUE",
                "value": 22425.99
            },
            "engineRunTime": {
                "@isTracked": "NO",
                "@readOnly": "TRUE",
                "value": "PT0.000S"
            },
            "licensePlate": {
                "@state": "NY",
                "@country": "US",
                "value": "XB3XXX"
            },
            "trackableItemType": "VEHICLE",
            "fuelType": "GASOLINE",
            "createdTimestamp": {
                "@readOnly": true,
                "value": "2016-12-06T17:58:22Z"
            },
            "modifiedTimestamp": {
                "@readOnly": true,
                "value": "2016-12-06T17:58:22Z"
            }
        },
        {
            "@id": 763800,
            "vin": {
                "@readOnly": true,
                "value": "XXT0X20FPXXX0XXX9"
            },
            "label": "LWB1",
            "color": "YELLOW",
            "make": "CAT",
            "model": "BACKHOE",
            "deviceSerialNumber": {
                "@readOnly": true,
                "value": "50XX6X7X4X"
            },
            "year": 2000,
            "odometer": {
                "@units": "MILES",
                "@timestamp": "2018-11-05T15:49:10Z",
                "@readOnly": "TRUE",
                "value": 5013.63
            },
            "engineRunTime": {
                "@isTracked": "NO",
                "@readOnly": "TRUE",
                "value": "PT0.000S"
            },
            "licensePlate": {
                "@state": "NY",
                "@country": "US",
                "value": "XJ8XXX"
            },
            "trackableItemType": "VEHICLE",
            "fuelType": "GASOLINE",
            "createdTimestamp": {
                "@readOnly": true,
                "value": "2016-12-06T17:55:35Z"
            },
            "modifiedTimestamp": {
                "@readOnly": true,
                "value": "2016-12-06T17:55:35Z"
            }
        }
        ]
    }
}`

var vehicleExpr = `payload.vehicle.{
    "id": ` + "`@id`," + `
    "label": label,
    "vin": vin.value,
    "serial": deviceSerialNumber.value,
    "make": make,
    "model": model,
    "year": year,
    "license": licensePlate.` + "`@state`" + ` & ' ' & licensePlate.value
  }`

var vehicleResult = `[{"id":942297,"label":"E M660","license":"NY AXXXXX","make":"CHEVROLET","model":"TAHOE","serial":"50X00XXX8X","vin":"1GNSK3EC1FRXXXXXX","year":2015},{"id":763803,"label":"LW04","license":"NY XB3XXX","make":"FORD","model":"F-150","serial":"5X147XX5XX","vin":"1FTFX1EF9EFXXXXXX","year":2014},{"id":763800,"label":"LWB1","license":"NY XJ8XXX","make":"CAT","model":"BACKHOE","serial":"50XX6X7X4X","vin":"XXT0X20FPXXX0XXX9","year":2000}]`
