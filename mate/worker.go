package mate

import (
	"context"
	"errors"
	"github.com/ericnts/cherry/base"
	"github.com/ericnts/cherry/util"
	"github.com/ericnts/orm"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session/v3"
	"gorm.io/gorm"
	"reflect"
	"sync"
)

var (
	workerType            = reflect.TypeOf((*Worker)(nil)).Elem()
	_              Worker = (*worker)(nil)
	TransactionKey        = "local_transaction_db"
)

func NewWorker(c context.Context) *worker {
	if c == nil {
		c = context.TODO()
	}
	bus := NewBus()
	bus.Set("tradingID", util.CreateUUID())
	return &worker{
		bus:     bus,
		context: c,
	}
}

type Worker interface {
	Context() context.Context
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Remove(key string)
	DB() *gorm.DB
	Session() (session.Store, error)
	Bus() *Bus
	Transaction(func() error) error
	AddEvent(event base.Event)
	AddFreeElement(element *ResourceElement)
	Free()
}

type worker struct {
	context      context.Context
	store        sync.Map
	bus          *Bus
	events       []base.Event
	freeElements []*ResourceElement
}

func (w *worker) Session() (session.Store, error) {
	ctx, ok := w.Context().(*gin.Context)
	if !ok {
		return nil, errors.New("current context is not gin.Context")
	}
	return session.Start(ctx, ctx.Writer, ctx.Request)
}

func (w *worker) DB() *gorm.DB {
	get, ok := w.Get(TransactionKey)
	if ok {
		return get.(*gorm.DB)
	}
	return orm.DB.WithContext(w.context)
}

func (w *worker) Transaction(f func() error) error {
	_, ok := w.Get(TransactionKey)
	if ok {
		return f()
	}
	return orm.DB.WithContext(w.context).Transaction(func(tx *gorm.DB) error {
		w.Set(TransactionKey, tx)
		defer w.Remove(TransactionKey)
		err := f()
		if err != nil {
			return err
		}
		w.pubEvent()
		return nil
	})
}

func (w *worker) AddEvent(event base.Event) {
	event.GetIdentity()
	w.events = append(w.events, event)
}

func (w *worker) pubEvent() {
	for _, pubEvent := range w.events {
		pubEvent.SetPrototypes(w.Bus().Clone())
	}
	w.events = nil
}

func (w *worker) Set(key string, value interface{}) {
	w.store.Store(key, value)
}

func (w *worker) Get(key string) (interface{}, bool) {
	return w.store.Load(key)
}

func (w *worker) Remove(key string) {
	w.store.Delete(key)
}

func (w *worker) Context() context.Context {
	return w.context
}

func (w *worker) AddFreeElement(element *ResourceElement) {
	w.freeElements = append(w.freeElements, element)
}

func (w *worker) Free() {
	for i := range w.freeElements {
		App.Put(w.freeElements[i])
	}
}

func (w *worker) Bus() *Bus {
	return w.bus
}
