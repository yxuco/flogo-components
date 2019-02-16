package fabric

// TriggerMap maps 'function' name in trigger handler setting to the trigger,
// so we can lookup trigger by chaincode function name
var TriggerMap map[string]*Trigger

func init() {
	TriggerMap = make(map[string]*Trigger)
}
