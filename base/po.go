package base

import (
	"github.com/ericnts/cherry/current"
	"github.com/ericnts/cherry/util"
	"gorm.io/gorm"
	"time"
)

type PO interface {
	TableName() string
	Update(name string, value interface{})
	HasChange(column string) bool
	IsChange() bool
	GetChanges() map[string]interface{}
}

type RPO interface {
	PO

	GetID() string
	SetID(id string)
}

type CPO interface {
	RPO

	ParentChanged() bool
	SetParentID(parentID string)
	SetParentIDs(parentIDs string)
	GetParentID() string
	GetParentIDs() string
}

type NormalPO struct {
	changes map[string]interface{} `gorm:"-"`
}

func (p *NormalPO) TableName() string {
	return ""
}

func (p *NormalPO) IsChange() bool {
	return len(p.changes) > 0
}

func (p *NormalPO) HasChange(column string) bool {
	_, ok := p.changes[column]
	return ok
}

func (p *NormalPO) GetChanges() map[string]interface{} {
	if p.changes == nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range p.changes {
		result[k] = v
	}
	p.changes = nil
	return result
}

func (p *NormalPO) Update(name string, value interface{}) {
	if p.changes == nil {
		p.changes = make(map[string]interface{})
	}
	p.changes[name] = value
}

type RecordPO struct {
	NormalPO

	ID        string         `gorm:"primarykey;<-:create;column:id;size:32;common:ID"` //主键
	UpdatedBy string         `gorm:"column:updated_by;size:32;comment:修改者ID"`          //修改者
	CreatedAt time.Time      `gorm:"<-:create;column:created_at;comment:创建时间"`         //创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;comment:修改时间"`                   //修改时间
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at;comment:删除时间"`             //删除时间
	Remark    string         `gorm:"column:remark;size:500;comment:描述信息"`              //描述信息
}

func (p *RecordPO) GetID() string {
	return p.ID
}

func (p *RecordPO) SetID(id string) {
	p.ID = id
}

func (p *RecordPO) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = util.CreateUUID()
	}
	if tx.Statement == nil {
		return nil
	}
	if tx.Statement.Context == nil {
		return nil
	}
	if p.UpdatedBy == "" {
		p.UpdatedBy = current.UserID(tx.Statement.Context)
	}
	return nil
}

func (p *RecordPO) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement == nil {
		return nil
	}
	if tx.Statement.Context == nil {
		return nil
	}
	if p.UpdatedBy == "" {
		p.UpdatedBy = current.UserID(tx.Statement.Context)
	}
	if len(p.UpdatedBy) == 0 {
		return nil
	}
	dest, ok := tx.Statement.Dest.(map[string]interface{})
	if ok {
		dest["updated_by"] = p.UpdatedBy
	}
	return nil
}

func (p *RecordPO) AfterDelete(tx *gorm.DB) error {
	if tx.Statement == nil || tx.Statement.Context == nil {
		return nil
	}
	if p.UpdatedBy == "" {
		p.UpdatedBy = current.UserID(tx.Statement.Context)
	}
	if len(p.UpdatedBy) == 0 {
		return nil
	}
	return tx.Model(tx.Statement.Model).
		Where(tx.Statement.Clauses["WHERE"].Expression).
		UpdateColumn("updated_by", p.UpdatedBy).Error
}

func (p *RecordPO) SetUpdatedBy(updatedBy string) {
	p.UpdatedBy = updatedBy
	p.Update("updated_by", updatedBy)
}

func (p *RecordPO) SetUpdatedAt(updatedAt time.Time) {
	p.UpdatedAt = updatedAt
	p.Update("updated_at", updatedAt)
}

func (p *RecordPO) SetDeletedAt(deletedAt gorm.DeletedAt) {
	p.DeletedAt = deletedAt
	p.Update("deleted_at", deletedAt)
}

func (p *RecordPO) SetRemark(remark string) {
	p.Remark = remark
	p.Update("remark", remark)
}

type CascadePO struct {
	RecordPO

	parentChanged bool   `gorm:"-"`
	ParentID      string `gorm:"column:parent_id;size:32;comment:父级ID"`
	ParentIDs     string `gorm:"column:parent_ids;size:500;comment:祖辈ID"`
}

func (p *CascadePO) ParentChanged() bool {
	return p.parentChanged
}

func (p *CascadePO) SetParentID(parentID string) {
	if len(p.ParentID) > 0 && p.ParentID == parentID {
		return
	}
	p.parentChanged = true
	p.ParentID = parentID
	p.Update("parent_id", parentID)
}

func (p *CascadePO) SetParentIDs(parentIDs string) {
	p.ParentIDs = parentIDs
	p.Update("parent_ids", parentIDs)
}

func (p *CascadePO) GetParentID() string {
	return p.ParentID
}

func (p *CascadePO) GetParentIDs() string {
	return p.ParentIDs
}
