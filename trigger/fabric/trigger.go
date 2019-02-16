package fabric

// check into git: https://github.com/yxuco/flogo-components.git
// add trigger in flogo-web: https://github.com/yxuco/flogo-components/trigger/fabric

import (
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// Create a new logger
var log = logger.GetLogger("trigger-dovetail-fabric")

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
		fn := handler.GetStringSetting("function")
		log.Info("init function:", fn)
		_, ok := TriggerMap[fn]
		if ok {
			log.Warnf("function %s used by multiple trigger handlers, only the last handler is effective")
		}
		TriggerMap[fn] = t
		args, ok := handler.GetSetting("args")
		if ok {
			log.Infof("init args: %T, %+v", args, args)
		} else {
			log.Info("init args not set")
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

// Invoke starts the trigger and invokes the action registered in the handler
func (t *Trigger) Invoke() error {
	log.Info("fabric.Trigger invoke")

	for _, handler := range t.handlers {
		fn := handler.GetStringSetting("function")
		log.Info("invoke function:", fn)
		args, ok := handler.GetSetting("args")
		if ok {
			log.Infof("invoke args: %T, %+v", args, args)
		} else {
			log.Info("invoke args not set")
		}
	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *Trigger) Stop() error {
	// stop the trigger
	return nil
}
