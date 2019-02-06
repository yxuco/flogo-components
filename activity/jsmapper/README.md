---
title: Jsmapper
weight: 
---

# Jsmapper
This activity allows you to map JSON data by specifying a JSONata expression.

## Installation
### Flogo Web
TODO
### Flogo CLI
```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/jsmapper
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "source",
      "type": "any",
      "required": true
    },
    {
      "name": "mapexpr",
      "type": "any"
    }
  ],
  "output": [
    {
      "name": "value",
      "type": "any"
    }
  ]
}
```

## Settings
| Setting  | Required | Description |
|:---------|:---------|:------------|
| source   | True     | Input JSON data, use normal flogo expression to construct it |         
| mapexpr  | False    | JSONata expression for data transformation |
| value    |          | Transformation result, or source data if expr is not specified |

## Mapping Examples
Query/translate JSON source data, it returns `value = "hello"`:

```json
{
  "id": "jsmapper_3",
  "name": "JSON Mapper",
  "description": "JSONata Mapper Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/jsmapper",
    "input": {
      "source": {"test": "hello"},
      "mapexpr": "test"
    }
  }
}
```
