package jsmapper

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/pkg/errors"
	"github.com/yxuco/flogo-components/activity/jsmapper/gojsonata"
)

var log = logger.GetLogger("activity-tibco-jsmapper")

const (
	ivMapexpr   = "mapexpr"
	ivSerialize = "serializeOutput"
	ovValue     = "value"

	envTag      = "env["
	flowTag     = "flow."
	activityTag = "activity["
)

// JsMapActivity is used to map data by evaluate JSONata script
type JsMapActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new JsMapActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &JsMapActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *JsMapActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *JsMapActivity) Eval(context activity.Context) (done bool, err error) {

	mapexpr := context.GetInput(ivMapexpr).(string)
	log.Debug("Mapper expression: ", mapexpr)

	doSerialize := context.GetInput(ivSerialize).(bool)
	log.Debug("Serialize output: ", doSerialize)

	// convert expr to JSONata expression, and extract names of required flow resources
	expr, attrs := prepareMapper(mapexpr)
	log.Debug("JSONata expr: ", expr)
	log.Debugf("required flow data: %+v", attrs)

	if len(attrs) == 0 {
		// no input data is referenced, so set output to a constant expr
		log.Warn("no valid input source is specified, return mapping expression as the result")
		context.SetOutput(ovValue, expr)
		return true, nil
	}

	// fetch all required flow resources to construct input data source as JSON string
	source, err := collectSource(attrs, context)
	if err != nil {
		log.Errorf("%+v", err)
		return false, err
	}
	log.Debugf("constructed source: %+v", source)

	// srcJSON, err := json.Marshal(source)
	// if err != nil {
	// 	log.Errorf("failed to marshal source JSON %+v", err)
	// 	return false, err
	// }
	// log.Debug("source data: ", string(srcJSON))

	// Transform source by applying JSONata expression
	value, err := gojsonata.Transform(source, expr)
	if err != nil {
		log.Errorf("failed JSONata transformation %+v", err)
		return false, err
	}

	if doSerialize {
		if result, err := json.Marshal(value); err == nil {
			// serialized JSON string
			log.Infof("JSON result: %+v\n", string(result))
			context.SetOutput(ovValue, string(result))
			return true, nil
		}
	}
	log.Infof("Transformation result: %+v\n", value)
	context.SetOutput(ovValue, value)
	return true, nil
}

func collectSource(attrs []string, context activity.Context) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	if envData, err := collectEnvData(attrs); err != nil {
		log.Errorf("failed to collect env data: %+v", err)
	} else {
		if envData != nil {
			data["env"] = envData
		}
	}
	if flowData, err := collectFlowData(attrs, context); err != nil {
		log.Errorf("failed to collect flow data: %+v", err)
	} else {
		if flowData != nil {
			data["flow"] = flowData
		}
	}
	activityData, err := collectActivityData(attrs, context)
	if err != nil {
		log.Errorf("failed to collect activity data: %+v", err)
		return nil, err
	}
	if activityData != nil {
		data["activity"] = activityData
	}
	return data, nil
}

func collectActivityData(attrs []string, context activity.Context) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	for _, a := range attrs {
		if strings.HasPrefix(a, activityTag) {
			// parse tag of format 'activity[key].name'
			keyEnd := strings.Index(a, "]")
			if keyEnd < 0 {
				return nil, errors.Errorf("invalid activity tag %s", a)
			}
			key := strings.TrimSpace(a[9:keyEnd])
			name := strings.TrimSpace(a[keyEnd+2:])

			// fetch data from flow context, and store it in map[string]interface{} for JSON
			value, err := GetData("$"+a, context)
			if err != nil {
				log.Warnf("error fetch data for %s: %+v", a, err)
				continue
			}
			activityMap, ok := data[key]
			if !ok {
				// create map to add a new activity
				activityMap = make(map[string]interface{})
				data[key] = activityMap
			}
			// add value of the attribute of 'name' to the activity of 'key'
			vMap := activityMap.(map[string]interface{})
			if s, ok := value.(string); ok {
				// value is a string, so try to unmarshal JSON
				js := unmarshalString(s)
				vMap[name] = js
			} else {
				vMap[name] = value
			}
		}
	}
	if len(data) > 0 {
		return data, nil
	}
	// no flow data is collected
	return nil, nil
}

func collectFlowData(attrs []string, context activity.Context) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	for _, a := range attrs {
		if strings.HasPrefix(a, flowTag) {
			key := strings.TrimSpace(a[5:])
			value, err := GetData("$"+a, context)
			if err != nil {
				log.Warnf("error fetch data for %s: %+v", a, err)
				continue
			}
			if s, ok := value.(string); ok {
				// value is a string, so try to unmarshal JSON
				js := unmarshalString(s)
				data[key] = js
			} else {
				data[key] = value
			}
		}
	}
	if len(data) > 0 {
		return data, nil
	}
	// no flow data is collected
	return nil, nil
}

// unmarshalString returns unmarshaled object if input is a valid JSON object or array,
// or returns the input string if it is not a valid JSON format
func unmarshalString(data string) interface{} {
	s := strings.TrimSpace(data)
	if strings.HasPrefix(s, "[") {
		var result []interface{}
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			return s
		}
		return result
	}
	if strings.HasPrefix(s, "{") {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			return s
		}
		return result
	}
	return s
}

func collectEnvData(attrs []string) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	for _, a := range attrs {
		if strings.HasPrefix(a, envTag) {
			tagEnd := strings.Index(a, "]")
			if tagEnd < 0 {
				return nil, errors.Errorf("invalid env tag %s", a)
			}
			env := strings.TrimSpace(a[4:tagEnd])
			value, ok := os.LookupEnv(env)
			if ok {
				data[env] = value
			}
		}
	}
	if len(data) > 0 {
		return data, nil
	}
	// no env is collected
	return nil, nil
}

// GetData resolves and returns data from the flow's context, unmarshals JSON string to map[string]interface{}.
// The name to Resolve is a valid output attribute of a flogo activity, e.g., `activity[app_16].value` or `$flow.content`,
// which is shown in normal flogo mapper as, e.g., "{{$flow.content}}"
func GetData(toResolve string, context activity.Context) (value interface{}, err error) {
	actionCtx := context.ActivityHost()
	actValue, err := actionCtx.GetResolver().Resolve(toResolve, actionCtx.WorkingData())
	return actValue, err
}

// prepareMapper parses mapperExpr to identify flog tags, and replace flog tags with corresponding JSONata expression.
// it returns the resulting JSONata expression and list of unique flogo tags for required source data.
func prepareMapper(mapperExpr string) (expr string, attrs []string) {
	attrMap := make(map[string]interface{})
	var buffer bytes.Buffer
	prefix, tag, suffix := splitNextTag(mapperExpr)
	tagReplacer := strings.NewReplacer("[", ".", "]", "")
	for suffix != "" {
		// store unique tag for constructing source data
		altTag := ""
		if tag != "" {
			if _, ok := attrMap[tag]; !ok {
				attrMap[tag] = nil
			}
			altTag = tagReplacer.Replace(tag)
		}
		buffer.WriteString(prefix)
		if altTag != "" {
			buffer.WriteString(altTag)
		}
		prefix, tag, suffix = splitNextTag(suffix)
	}
	buffer.WriteString(prefix)

	// collect unique flogo tags
	attrs = []string{}
	for k := range attrMap {
		attrs = append(attrs, k)
	}
	return buffer.String(), attrs
}

func splitNextTag(mapperExpr string) (prefix, tag, suffix string) {
	foundTag := false
	tagStart := strings.Index(mapperExpr, "{{$")
	for !foundTag {
		if tagStart < 0 {
			// no more tags
			return mapperExpr, "", ""
		}
		if hasFlogoTag(mapperExpr[tagStart+3:]) {
			foundTag = true
		} else {
			// not a Flogo tag, so find the next tag
			nextStart := strings.Index(mapperExpr[tagStart+3:], "{{$")
			if nextStart < 0 {
				tagStart = nextStart
			} else {
				tagStart += nextStart + 3
			}
		}
	}
	tagEnd := strings.Index(mapperExpr[tagStart:], "}}")
	if tagEnd < 0 {
		// no tag end, return original although it could be error
		return mapperExpr, "", ""
	}
	tagEnd += tagStart
	tag = mapperExpr[tagStart+3 : tagEnd]

	// verify that tag contains at least one dot
	propTag, tagExtra := validateFlogoTag(tag)
	if propTag == "" {
		// invalid tag, ignore it with blank tag
		if tagEnd+2 < len(mapperExpr) {
			return mapperExpr[0 : tagEnd+2], "", mapperExpr[tagEnd+2:]
		}
		return mapperExpr, "", ""
	}

	// remove quotes around the Flogo tag if they exist
	if tagStart > 0 && tagEnd+3 < len(mapperExpr) {
		if mapperExpr[tagStart-1:tagStart] == "\"" &&
			mapperExpr[tagEnd+2:tagEnd+3] == "\"" {
			tagStart--
			tagEnd++
		}
	}
	return mapperExpr[0:tagStart], propTag, tagExtra + mapperExpr[tagEnd+2:]
}

// valiateFlogoTag verifies that proper tag contains only top-level attribute
// return proper tag and move second-level nested attribute to tagExtra
// return blank tag if input tag is not a valid flogo tag
func validateFlogoTag(tag string) (propTag, tagExtra string) {
	// validate env tag $env[VAR]
	if strings.HasPrefix(tag, envTag) {
		if strings.Index(tag, "]") == len(tag)-1 {
			// valid proper tag
			return tag, ""
		}
		return "", ""
	}

	// verify that tag contains at least one dot
	dotPos := strings.Index(tag, ".")
	if dotPos < 0 {
		// invalid tag
		return "", ""
	}

	// keep only the first "." attribute in tag, move the rest to suffix
	tagExtra = ""
	propTag = tag
	extraPos := strings.Index(tag[dotPos+1:], ".")
	if extraPos > 0 {
		tagExtra = tag[dotPos+extraPos+1:]
		propTag = tag[0 : dotPos+extraPos+1]
	}

	// validate activity tag format: $activity[A].prop
	if strings.HasPrefix(propTag, activityTag) {
		if strings.Index(propTag, "].") <= 0 {
			// invalid action tag
			return "", ""
		}
	}
	return propTag, tagExtra
}

// hasFlogoTag returns true if an expression starts with "flow.", "env." or "activity."
func hasFlogoTag(mapperExpr string) bool {
	if strings.HasPrefix(mapperExpr, activityTag) {
		return true
	}
	if strings.HasPrefix(mapperExpr, flowTag) {
		return true
	}
	if strings.HasPrefix(mapperExpr, envTag) {
		return true
	}
	return false
}
