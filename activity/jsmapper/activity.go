package jsmapper

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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

	log.Infof("Context %T: %+v", context, context)

	// TODO: transform data here
	value := source
	context.SetOutput(ovValue, value)
	return true, nil
}
