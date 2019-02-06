//go:generate go run $GOPATH/src/github.com/TIBCOSoftware/flogo-lib/flogo/gen/gen.go $GOPATH
package main

import (
    // Default Go packages
	"context"
	"os"
	"strconv"

    // The Flogo Log activity
    "github.com/TIBCOSoftware/flogo-contrib/activity/log"
    // The Flogo REST trigger
    "github.com/TIBCOSoftware/flogo-contrib/trigger/rest"
    // Core packages for the Flogo engine
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/flogo"
    "github.com/TIBCOSoftware/flogo-lib/logger"
)

var (
	httpport = os.Getenv("HTTPPORT")
)

func main() {
	// Create a new Flogo app
	app := appBuilder()

	e, err := flogo.NewEngine(app)

	if err != nil {
		logger.Error(err)
		return
	}

	engine.RunEngine(e)
}

func appBuilder() *flogo.App {
	app := flogo.NewApp()

	// Convert the HTTPPort to an integer
	port, err := strconv.Atoi(httpport)
	if err != nil {
		logger.Error(err)
	}

	// Register the HTTP trigger
	trg := app.NewTrigger(&rest.RestTrigger{}, map[string]interface{}{"port": port})
	trg.NewFuncHandler(map[string]interface{}{"method": "GET", "path": "/cheesecake/:name"}, Handler)

	return app
}

// Handler is the function that gets executed when the engine receives a message
func Handler(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	// Get the name from the path
	name := inputs["pathParams"].Value().(map[string]string)["name"]

	// Log, using the Flogo log activity
	// There are definitely better ways to do this with Go, but we want to show how to use activities
	in := map[string]interface{}{"message": name, "flowInfo": "true", "addToFlow": "true"}
	_, err := flogo.EvalActivity(&log.LogActivity{}, in)
	if err != nil {
		return nil, err
	}

    // Set the result message
    var cheesecake string
    switch name {
    case "retgits":
        cheesecake = "Likes all cheesecakes"
    case "flynn":
        cheesecake = "Prefers some nectar!"
    default:
        cheesecake = "Plain cheesecake is the best"
    }

	// The return message is a map[string]*data.Attribute which we'll have to construct
	response := make(map[string]interface{})
	response["name"] = name
	response["cheesecake"] = cheesecake

	ret := make(map[string]*data.Attribute)
	ret["code"], _ = data.NewAttribute("code", data.TypeInteger, 200)
	ret["data"], _ = data.NewAttribute("data", data.TypeAny, response)

	return ret, nil
}