package mate

import (
	"context"
	"errors"
	"fmt"
	"github.com/ericnts/cherry/middleware"
	"github.com/ericnts/config"
	"github.com/ericnts/log"
	"github.com/ericnts/orm"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"sync"
)

var (
	App         *Application
	prepareOnce sync.Once
)

func init() {
	App = new(Application)
	App.Engine = gin.New()
	App.routers = make([]routerItem, 0, 10)
	App.prepares = make(map[PrepareLevel][]func(application *Application), 3)
	App.prepares[PLHeight] = make([]func(application *Application), 0, 3)
	App.prepares[PLMiddle] = make([]func(application *Application), 0, 3)
	App.prepares[PLLow] = make([]func(application *Application), 0, 3)
	App.poolMap = make(map[ResourceType]*resourcePool)
	App.poolMap[RTInfrastructure] = newResourcePool(RTInfrastructure)
	App.poolMap[RTRepository] = newResourcePool(RTRepository)
	App.poolMap[RTFactory] = newResourcePool(RTFactory)
	App.poolMap[RTService] = newResourcePool(RTService)
	App.poolMap[RTApplication] = newResourcePool(RTApplication)
	App.poolMap[RTController] = newResourcePool(RTController)
}

type routerItem struct {
	middlewares []gin.HandlerFunc
	controllers map[string]reflect.Type
}

type Application struct {
	Engine      *gin.Engine
	RouterGroup *gin.RouterGroup
	middlewares []gin.HandlerFunc
	routers     []routerItem
	poolMap     map[ResourceType]*resourcePool
	prepares    map[PrepareLevel][]func(*Application)
	afters      []func(*Application)
}

func (a *Application) DB() *gorm.DB {
	return orm.DB
}

func (a *Application) DoPrepare() {
	prepareOnce.Do(func() {
		for i := range a.prepares[PLHeight] {
			a.prepares[PLHeight][i](a)
		}
		for i := range a.prepares[PLMiddle] {
			a.prepares[PLMiddle][i](a)
		}
		for i := range a.prepares[PLLow] {
			a.prepares[PLLow][i](a)
		}
	})
}

func (a *Application) Run(handlers ...gin.HandlerFunc) {
	a.DoPrepare()

	a.Engine.Use(gin.Recovery(), middleware.CorsMiddleware)
	a.Engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	baseRouter := a.Engine.Group(config.Options.HttpPrefix)
	if config.Options.Swagger {
		baseRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("doc.json")))
		defer log.Infof("swagger: http://localhost:%v/%s/swagger/index.html", config.Options.HttpPort, config.Options.Name)
	}
	if config.Options.Pprof {
		pprof.RouteRegister(baseRouter)
	}
	a.RouterGroup = baseRouter.Group("/", handlers...)

	for _, routerItem := range a.routers {
		currentRouter := a.RouterGroup.Group("/", routerItem.middlewares...)
		for path, cType := range routerItem.controllers {
			for i := 0; i < cType.NumMethod(); i++ {
				fun := cType.Method(i)
				method, methodPath := parseMethodName(fun.Name)
				if len(method) == 0 {
					continue
				}
				work := NewWorker(context.TODO())
				if _, ok := a.Get(work, RTController, cType); !ok {
					log.Panicf("未找到Controller %s.%s", cType.Elem().PkgPath(), cType.Elem().Name())
				}
				work.Free()

				currentType := cType
				currentRouter.Handle(method, path+methodPath, func(context *gin.Context) {
					work := NewWorker(context)
					rElement, ok := a.Get(work, RTController, currentType)
					if !ok {
						panic(errors.New("获取控制器失败"))
					}
					fun.Func.Call([]reflect.Value{rElement.Value})
					work.Free()
				})
			}
		}
	}
	go func() {
		if err := a.Engine.Run(fmt.Sprintf(":%v", config.Options.HttpPort)); err != nil {
			log.WithError(err).Info("http服务启动失败")
			panic(err)
		}
	}()

	for i := range a.afters {
		a.afters[i](a)
	}
}

func (a *Application) Get(work Worker, rType ResourceType, cType reflect.Type) (*ResourceElement, bool) {
	pool := a.poolMap[rType]
	if pool == nil {
		return emptyResourceElement, false
	}
	element, ok := pool.get(cType)
	if !ok {
		return emptyResourceElement, false
	}
	element.setWork(work)
	return element, true
}

func (a *Application) GetAnonymous(work Worker, cType reflect.Type) (*ResourceElement, bool) {
	for _, pool := range a.poolMap {
		element, ok := pool.get(cType)
		if !ok {
			continue
		}
		element.setWork(work)
		return element, true
	}
	return emptyResourceElement, false
}

func (a *Application) Put(element *ResourceElement) {
	for rType, pool := range a.poolMap {
		if element.resourceType == rType {
			pool.put(element)
		}
	}
}

func (a *Application) Prepare(level PrepareLevel, f func(*Application)) {
	a.prepares[level] = append(a.prepares[level], f)
}

func (a *Application) After(f func(*Application)) {
	a.afters = append(a.afters, f)
}

func (a *Application) BindController(path string, f interface{}, handlers ...gin.HandlerFunc) {
	a.Prepare(PLMiddle, func(app *Application) {
		outType, err := parsePoolFunc(f)
		if err != nil {
			log.Fatalf("The binding function is incorrect, %v : %s", f, err.Error())
		}
		app.poolMap[RTController].bind(false, f)
		app.routers = append(app.routers, routerItem{middlewares: handlers, controllers: map[string]reflect.Type{path: outType}})
	})
}

func (a *Application) BindControllers(controllerMap map[string]interface{}, handlers ...gin.HandlerFunc) {
	a.Prepare(PLMiddle, func(app *Application) {
		controllers := make(map[string]reflect.Type, len(controllerMap))
		for path, f := range controllerMap {
			outType, err := parsePoolFunc(f)
			if err != nil {
				log.Fatalf("The binding function is incorrect, %v : %s", f, err.Error())
			}
			app.poolMap[RTController].bind(false, f)
			controllers[path] = outType
		}
		app.routers = append(app.routers, routerItem{middlewares: handlers, controllers: controllers})
	})
}

func (a *Application) BindApplication(f interface{}) {
	a.poolMap[RTApplication].bind(false, f)
}

func (a *Application) BindService(f interface{}) {
	a.poolMap[RTService].bind(false, f)
}

func (a *Application) BindFactory(f interface{}) {
	a.poolMap[RTFactory].bind(false, f)
}

func (a *Application) BindRepository(f interface{}) {
	a.poolMap[RTRepository].bind(false, f)
}

func (a *Application) BindInfrastructure(single bool, com interface{}) {
	a.poolMap[RTInfrastructure].bind(single, com)
}
