package mate

import (
	"context"
	"github.com/ericnts/cherry/current"
	"github.com/ericnts/cherry/results"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session/v3"
	"gorm.io/gorm"
)

type Resource struct {
	Current Worker
}

func (r *Resource) Context() context.Context {
	return r.worker().Context()
}

func (r *Resource) Session() (session.Store, error) {
	return r.worker().Session()
}

func (r *Resource) Transaction(f func() error) error {
	return r.worker().Transaction(f)
}

func (r *Resource) DB() *gorm.DB {
	return r.worker().DB()
}

func (r *Resource) worker() Worker {
	if r.Current == nil {
		r.Current = NewWorker(context.TODO())
	}
	return r.Current
}

type GinResource struct {
	Resource
}

func (r *GinResource) Context() *gin.Context {
	return r.worker().Context().(*gin.Context)
}

func (r *GinResource) Bind(obj interface{}) error {
	return r.Context().Bind(obj)
}

func (r *GinResource) Query(key string) string {
	return r.Context().Query(key)
}

func (r *GinResource) Param(key string) string {
	return r.Context().Param(key)
}

func (r *GinResource) OfficeID() string {
	return current.OfficeID(r.Context())
}

func (r *GinResource) ID() string {
	return current.ID(r.Context())
}

func (r *GinResource) Err(err error) {
	results.Err(r.Context(), err)
}

func (r *GinResource) Success(data interface{}, page ...*results.Page) {
	results.Success(r.Context(), data, page...)
}
