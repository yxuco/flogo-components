package fabricop

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	//	"github.com/hyperledger/fabric/core/chaincode/shim"
	fabtrigger "github.com/yxuco/flogo-components/trigger/fabric"
)

// Create a new logger
//var log = shim.NewLogger("activity-tibco-fabricop")
var log = logger.GetLogger("activity-tibco-fabricop")

const (
	ivOperation = "operation"
	ivKey       = "key"
	ivData      = "data"
	ivFilter    = "filter"
	ovResult    = "result"
)

// FabActivity is used to execute a hyperledger fabric operation
type FabActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new new FabActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabActivity) Eval(ctx activity.Context) (done bool, err error) {

	// check operation type, UI cannot config settings, so use input args
	//if op, ok := ctx.GetSetting(sOperation); ok {
	//	log.Infof("perform operation: %s", op.(string))
	//}

	// check input args
	op := ctx.GetInput(ivOperation)
	log.Debugf("perform operation: %s", op)
	key := ctx.GetInput(ivKey)
	log.Debugf("input key: %s", key)
	data := ctx.GetInput(ivData)
	log.Debugf("input data: %+v", data)
	filter := ctx.GetInput(ivFilter)
	log.Debugf("input filter: %+v", filter)

	// get chaincode stub
	stub, err := GetData("$flow."+fabtrigger.FabricStub, ctx)
	if err != nil {
		log.Errorf("failed to get stub: %+v", err)
	} else {
		log.Infof("fetched stub of type %T: %+v", stub, stub)
	}

	// set output
	ctx.SetOutput(ovResult, "done")
	return true, nil
}

// GetData resolves and returns data from the flow's context, unmarshals JSON string to map[string]interface{}.
// The name to Resolve is a valid output attribute of a flogo activity, e.g., `activity[app_16].value` or `$flow.content`,
// which is shown in normal flogo mapper as, e.g., "{{$flow.content}}"
func GetData(toResolve string, context activity.Context) (value interface{}, err error) {
	actionCtx := context.ActivityHost()
	log.Debugf("fabricop context data: %+v", actionCtx.WorkingData())
	actValue, err := actionCtx.GetResolver().Resolve(toResolve, actionCtx.WorkingData())
	return actValue, err
}
