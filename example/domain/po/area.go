package po

import "github.com/ericnts/cherry/base"

type Area struct {
	base.CascadePO

	Users []User `gorm:"many2many:user_area;" cherry:"auto"`

	Name string `gorm:"column:name;not null;comment:名称" cherry:"index:#Type_Name"`
	Code string `gorm:"column:code;not null;comment:编码" cherry:"index"`
	Type int    `gorm:"column:type;comment:类型"`
}

func (Area) TableName() string {
	return "areas"
}

func (p *Area) SetName(name string) {
	p.Name = name
	p.Update("name", name)
}
