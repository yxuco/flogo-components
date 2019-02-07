package jsmapper

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("activity-tibco-jsmapper")

const (
	ivSource  = "source"
	ivMapexpr = "mapexpr"

	ovValue = "value"
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

	// TODO: process either string or JSON data
	source := context.GetInput(ivSource)
	switch v := source.(type) {
	case string:
		log.Info("Input string data:", v)
	case map[string]interface{}:
		log.Infof("Input JSON object %T: %+v", v, v)
	default:
		log.Infof("Input other data type %T, %+v", v, v)
	}
	mapexpr := context.GetInput(ivMapexpr).(string)
	log.Info("Mapper expression:", mapexpr)

	actionCtx := context.ActivityHost()
	log.Infof("ActivitHost %T: %+v", actionCtx, actionCtx)
	wd := actionCtx.WorkingData()
	log.Infof("WorkingData %T: %+v", wd, wd)

	switch v := wd.(type) {
	case *data.FixedScope:
		log.Infof("FixedScope %T: %+v", v, v)
	default:
		log.Infof("working data scope type %T: %+v", v, v)
	}

	switch v := context.(type) {
	default:
		log.Infof("Context %T: %+v", v, v)
	}
	actValue, err := data.GetBasicResolver().Resolve("$activity[app_16].value", wd)
	if err != nil {
		log.Errorf("failed to resolve %+v", err)
	} else {
		log.Infof("resolved flow data %T: %+v", actValue, actValue)
	}

	// TODO: transform data here
	value := source
	context.SetOutput(ovValue, value)
	return true, nil
}
