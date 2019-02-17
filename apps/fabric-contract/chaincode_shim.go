package main

import (
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	trigger "github.com/yxuco/flogo-components/trigger/fabric"
)

// Contract implements chaincode interface for invoking Flogo flows
type Contract struct {
}

var logger = shim.NewLogger("chaincode-shim")

// Init is called during chaincode instantiation to initialize any data,
// and also calls this function to reset or to migrate data.
func (t *Contract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (t *Contract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	logger.Debugf("invoke transaction fn=%s, args=%+v", fn, args)

	trig, ok := trigger.GetTrigger(fn)
	if !ok {
		return shim.Error(fmt.Sprintf("function %s is not implemented", fn))
	}
	result, err := trig.Invoke(stub, fn, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to execute transaction: %s, error: %+v", fn, err))
	}
	return shim.Success([]byte(result))
}

var (
	cp app.ConfigProvider
)

// main function starts up the chaincode in the container during instantiate
func main() {
	// configure flogo engine
	if cp == nil {
		// Use default config provider
		cp = app.DefaultConfigProvider()
	}

	app, err := cp.GetApp()
	if err != nil {
		fmt.Printf("failed to read Flogo app config: %+v\n", err)
		os.Exit(1)
	}

	e, err := engine.New(app)
	if err != nil {
		fmt.Printf("Failed to create flogo engine instance: %+v\n", err)
		os.Exit(1)
	}

	if err := e.Init(true); err != nil {
		fmt.Printf("Failed to initialize flogo engine: %+v\n", err)
		os.Exit(1)
	}

	if err := shim.Start(new(Contract)); err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
