package fabric

// check into git: https://github.com/yxuco/flogo-components.git
// add trigger in flogo-web: https://github.com/yxuco/flogo-components/trigger/fabric

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	sFunction = "function"
	sArgs     = "args"
	oData     = "data"
	rResult   = "result"
	cStub     = "chaincode-stub"
)

// Create a new logger
var logger = shim.NewLogger("fabric-invoke")

// TriggerMap maps 'function' name in trigger handler setting to the trigger,
// so we can lookup trigger by chaincode function name
var triggerMap = map[string]*Trigger{}

// GetTrigger returns the cached trigger for a specified function name;
// return false in the second value if no trigger is cached for the specified fn
func GetTrigger(fn string) (*Trigger, bool) {
	trig, ok := triggerMap[fn]
	return trig, ok
}

// TriggerFactory Fabric Trigger factory
type TriggerFactory struct {
	metadata *trigger.Metadata
}

// NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &TriggerFactory{metadata: md}
}

// New Creates a new trigger instance for a given id
func (t *TriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &Trigger{metadata: t.metadata, config: config}
}

// Trigger is a stub for your Trigger implementation
type Trigger struct {
	metadata *trigger.Metadata
	config   *trigger.Config
	handlers []*trigger.Handler
}

// Initialize implements trigger.Init.Initialize
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	t.handlers = ctx.GetHandlers()
	for _, handler := range t.handlers {
		fn := handler.GetStringSetting(sFunction)
		logger.Info("init function:", fn)
		_, ok := triggerMap[fn]
		if ok {
			logger.Warningf("function %s used by multiple trigger handlers, only the last handler is effective", fn)
		}
		triggerMap[fn] = t
		args, ok := handler.GetSetting(sArgs)
		if ok {
			logger.Infof("init args: %T, %+v", args, args)
		} else {
			logger.Info("init args not set")
		}
	}
	return nil
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *Trigger) Start() error {
	return nil
}

// Invoke starts the trigger and invokes the action registered in the handler,
// and returns result as JSON string
func (t *Trigger) Invoke(stub shim.ChaincodeStubInterface, fn string, args []string) (string, error) {
	logger.Debugf("fabric.Trigger invoke fn %s with args %+v", fn, args)
	for _, handler := range t.handlers {
		if f := handler.GetStringSetting(sFunction); f != fn {
			logger.Warningf("handler function %s is different from requested function %s", f, fn)
			continue
		}

		// construct trigger data
		var names []string
		if argNames, ok := handler.GetSetting(sArgs); ok {
			names, ok = argNames.([]string)
		}
		data := prepareFlowData(names, args)
		if data == nil {
			return "", fmt.Errorf("failed to prepare trigger data from input %+v", args)
		}

		// debug flow data
		triggerData, _ := json.Marshal(data)
		logger.Debugf("trigger output data: %s", string(triggerData))

		flowData := make(map[string]interface{})
		flowData[oData] = data
		flowData[cStub] = stub

		// execute flogo flow
		results, err := handler.Handle(context.Background(), flowData)
		if err != nil {
			return "", err
		}
		if len(results) != 0 {
			if dataAttr, ok := results[rResult]; ok {
				replyData := dataAttr.Value()
				if s, ok := replyData.(string); ok {
					return s, nil
				}
				logger.Infof("flogo flow returned data type %T is not a string: %+v", replyData, replyData)
			} else {
				logger.Infof("flogo flow result does not contain attribute %s", rResult)
			}
		}
		logger.Info("flogo flow did not return any data")
		return "", nil
	}
	logger.Infof("no flogo handler is activated for function %s", fn)
	return "", nil
}

func prepareFlowData(names, values []string) interface{} {
	if names == nil || len(names) == 0 {
		// construct array of objects
		var result []interface{}
		for _, v := range values {
			result = append(result, unmarshalString(v))
		}
		return result
	}

	// convert array to object with name-values
	result := make(map[string]interface{})
	for i, v := range values {
		var key string
		if len(names)-1 >= i {
			key = names[i]
		} else {
			// make up name if value list is longer than name list
			key = "arg" + strconv.Itoa(i)
		}
		result[key] = unmarshalString(v)
	}
	return result
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

// Stop implements trigger.Trigger.Start
func (t *Trigger) Stop() error {
	// stop the trigger
	return nil
}
