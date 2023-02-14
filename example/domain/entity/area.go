package entity

import (
	"github.com/ericnts/cherry/base"
	"github.com/ericnts/cherry/example/domain/po"
)

type Area struct {
	base.EventEntity
	po.Area
}
