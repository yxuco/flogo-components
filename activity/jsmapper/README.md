---
title: Jsmapper
weight: 
---

# Jsmapper
This activity allows you to map JSON data by specifying a [JSONata expression](http://docs.jsonata.org/overview.html).

## Installation
### Flogo Web
On "Add Activity dialog", select "Install New Activity", then paste the following URL: https://github.com/yxuco/flogo-components/tree/master/activity/jsmapper

### Flogo CLI
```bash
flogo install github.com/yxuco/flogo-components/activity/jsmapper
```

## Dependencies
The implementation of this activity depends on the Flogo libs, i.e.,
```bash
go get -u github.com/TIBCOSoftware/flogo-lib/...
go get -u github.com/TIBCOSoftware/flogo-contrib/...
go get -u github.com/TIBCOSoftware/flogo-cli/...
go get -u github.com/TIBCOSoftware/flogo/...
```

Besides, this activity also depends on the following Go packages:
```bash
go get -u github.com/pkg/errors
go get -u github.com/robertkrimen/otto
go get -u github.com/stretchr/testify
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "mapexpr",
      "type": "string",
      "required": true
    },
    {
      "name": "serializeOutput",
      "type": "boolean"
    }
  ],
  "output": [
    {
      "name": "value",
      "type": "string"
    }
  ]
}
```

## Settings
| Setting          | Required | Description |
|:-----------------|:---------|:------------|
| mapexpr          | True     | JSONata expression for data transformation |
| serializeOutput  | False    | Serialize the result to JSON string if it is true |
| value            |          | Transformation result, or the expr itself if it does not contain any source tag |

## Mapping Examples
The sample application [transform-app](https://github.com/yxuco/flogo-components/tree/master/apps/transform-app) contains the following activity that transforms output data from other activities in the same flow.

```json
{
  "id": "jsmapper_18",
  "name": "JSON Mapper",
  "description": "JSONata Mapper Activity",
  "activity": {
    "ref": "github.com/yxuco/flogo-components/activity/jsmapper",
    "input": {
      "mapexpr": "jsonata expression",
      "serializeOutput": false
    }
  }
}
```
The `mapexpr` is a JSONata expression containing tags of flogo data, including envirionment variables, flow properties, and/or output data from other flogo activities reachable in the same flow.  For example, `transform-app` contains the following expression:
```
{
  "port": {{$env[HTTPPORT]}},
  "event": {{$flow.event}},
  "street": {{$flow.content}}.address.street,
  "message": {{$activity[log_2].message}},
  "vehicle": {{$activity[app_16].value}}.vehicle.{
    "make": make,
    "year": year,
    "vin": vin.value
  }
}
```
The above expression includes the following data items from the flogo flow:
* An environment variable `HTTPPORT`, which could be specified when launching the application;
* 2 flow properties `event` and `content`, which could be set by mappings in the trigger.  The `content` in this example is a JSON string that contains a street attribute in an address;
* A simple string attribute `message` from the activity `log_2`;
* A complex JSON attribute `value` from the activity `app_16`, which contains multiple `vehicle` objects.  The details of the vehicles are shown in [inventory.json](https://github.com/yxuco/flogo-components/tree/master/apps/transform-app/inventory.json).

The output of this transformation is as follows:

```json
{
  "event": "list",
  "message": "Received event list",
  "port": "8080",
  "street": "3033 Hillsview Ave",
  "vehicle": [
    {
      "make": "CHEVROLET",
      "vin": "1GNSK3EC1FRXXXXXX",
      "year": 2015
    },
    {
      "make": "FORD",
      "vin": "1FTFX1EF9EFXXXXXX",
      "year": 2014
    },
    {
      "make": "CAT",
      "vin": "XXT0X20FPXXX0XXX9",
      "year": 2000
    }
  ]
}
```

When `serializeOutput` is set to `true`, the transformation result will be returned as a serialized JSON string.  However, Flogo handles the result as JSON object even if it is serialized.  For example, in a Flogo activity following a `jsmapper`, you can access an attribute in the mapper result by e.g., `$activity[jsmapper_18].value.vehicle[0].make`, even if the output value of `jsmapper_18` is serialized to a string.

## Transformation Expression
In the transformation expression, Flogo data sources must be tagged in the form of `{{$...}}`.  On the Flogo UI, you can enter these tags by clicking an item in the list of "Available Data".  You can then add double curly brackets around the data tag.

Except for these data tags, the expression should match [JSONata](http://jsonata.org/) specification.  The expression can use any of the JSONata [functions](http://docs.jsonata.org/string-functions), e.g., for string manipulation, and aggregation, etc.

You can develop and test JSONata expressions at http://try.jsonata.org/.

This transformer activity does not support Golang functions at present.  However, if it is necessary, it can be extended to support Golang functions, as well as other custom JavaScript functions besides the core functions provided by JSONata.

## Future Improvement
The editor for the JSONata expression can be improved, so it can at least support multi-line display.  It would be even better if Flogo UI can generate the JSONata expression by visual drag-and-drop.