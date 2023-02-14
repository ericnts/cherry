package infrastructure

import (
	"github.com/ericnts/cherry"
	"github.com/ericnts/cherry/example/domain/entity"
	_ "github.com/ericnts/cherry/example/infrastructure/persistence"
	"github.com/ericnts/cherry/mate"
	"github.com/ericnts/orm"
)

func init() {
	cherry.Prepare(AutoMigrate)
}

func AutoMigrate(app *mate.Application) {
	err := orm.DB.AutoMigrate(
		new(entity.User),
		new(entity.Area),
	)
	if err != nil {
		panic(err)
	}
}
