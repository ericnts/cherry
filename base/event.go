package base

import "github.com/ericnts/cherry/util"

type Event interface {
	Topic() string
	GetPrototypes() map[string]interface{}
	SetPrototypes(map[string]interface{})
	GetIdentity() string
	SetIdentity(identity string)
}

type NormalEvent struct {
	Identity   string
	Prototypes map[string]interface{}
}

func (e *NormalEvent) GetPrototypes() map[string]interface{} {
	return e.Prototypes
}

func (e *NormalEvent) SetPrototypes(prototypes map[string]interface{}) {
	e.Prototypes = prototypes
}

func (e *NormalEvent) GetIdentity() string {
	if len(e.Identity) == 0 {
		e.Identity = util.CreateUUID()
	}
	return e.Identity
}

func (e *NormalEvent) SetIdentity(identity string) {
	e.Identity = identity
}
