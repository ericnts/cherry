package base

import (
	"github.com/ericnts/cherry/util"
	"reflect"
)

type Entity interface {
	PO

	GetPubEvents() (result []Event)
	RemoveAllPubEvent()
	AddSubEvent(event Event)
	GetSubEvents() (result []Event)
	RemoveAllSubEvent()
}

type EventEntity struct {
	pubEvents []Event `gorm:"-"`
	subEvents []Event `gorm:"-"`
}

func (e *EventEntity) AddPubEvent(event Event) {
	if reflect.ValueOf(event.GetIdentity()).IsZero() {
		event.SetIdentity(util.CreateUUID())
	}
	e.pubEvents = append(e.pubEvents, event)
}

func (e *EventEntity) GetPubEvents() (result []Event) {
	return e.pubEvents
}

func (e *EventEntity) RemoveAllPubEvent() {
	e.pubEvents = []Event{}
}

func (e *EventEntity) AddSubEvent(event Event) {
	e.subEvents = append(e.subEvents, event)
}

func (e *EventEntity) GetSubEvents() (result []Event) {
	return e.subEvents
}

func (e *EventEntity) RemoveAllSubEvent() {
	e.subEvents = []Event{}
}
