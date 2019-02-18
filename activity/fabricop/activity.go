package fabricop

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//	fabtrigger "github.com/yxuco/flogo-components/trigger/fabric"
)

// Create a new logger
var logger = shim.NewLogger("activity-tibco-fabricop")

const (
	sOperation = "operation"
	ivKey      = "key"
	ivData     = "data"
	ivFilter   = "filter"
	ovResult   = "result"

	fTxID   = "$flow.txID"
	fTxTime = "$flow.txTime"
	fStub   = "$flow.chaincode-stub"
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

	// check operation type
	if op, ok := ctx.GetSetting(sOperation); ok {
		logger.Infof("perform operation: %s", op.(string))
	}

	// check input args
	key := ctx.GetInput(ivKey)
	logger.Debugf("input key: %s", key)
	data := ctx.GetInput(ivData)
	logger.Debugf("input data: %+v", data)
	filter := ctx.GetInput(ivFilter)
	logger.Debugf("input filter: %+v", filter)

	// get chaincode stub
	stub, err := GetData("$flow.chaincode_stab", ctx)
	if err != nil {
		logger.Errorf("failed to get stub: %+v", err)
	} else {
		logger.Infof("fetched stub of type %T: %+v", stub, stub)
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
	logger.Debugf("fabricop context data: %+v", actionCtx.WorkingData())
	actValue, err := actionCtx.GetResolver().Resolve(toResolve, actionCtx.WorkingData())
	return actValue, err
}
