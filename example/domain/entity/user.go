package entity

import (
	"github.com/ericnts/cherry/base"
	"github.com/ericnts/cherry/example/domain/po"
)

type User struct {
	base.EventEntity
	po.User
}
