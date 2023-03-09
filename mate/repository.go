package mate

import (
	"errors"
	"fmt"
	"github.com/ericnts/cherry/base"
	"github.com/ericnts/cherry/exception"
	"gorm.io/gorm"
	"reflect"
)

var _ Repo[base.PO] = (*Repository[base.PO])(nil)

type Repo[P base.PO] interface {
	Create(p P) (string, error)
	FindExport(query string, args []interface{}, preloads []interface{}, order ...string) ([]P, error)
	PageExport(limit, offset int, query string, args []interface{}, preloads []interface{}, order ...string) ([]P, int64, error)
	FindInIDs(ids []string, preloads ...interface{}) ([]P, error)
	All(preloads ...interface{}) ([]P, error)
	GetByID(id string, preloads ...interface{}) (P, error)
	Update(record P) (int64, error)
	DeleteByID(ids ...string) (int64, error)
	DeleteAssociationByID(association string, ids ...string) error
	ClearAutoUpdateAssociationByID(ids ...string) error
	ClearAutoDeleteAssociationByID(ids ...string) error
	// UnscopedDeleteByID 物理删除
	UnscopedDeleteByID(ids ...string) (int64, error)
	CheckIndex(po P) error
	CheckVerify(po P, ids ...string) error
}

type Repository[P base.PO] struct {
	Resource
}

func (r *Repository[P]) FindInIDs(ids []string, preloads ...interface{}) ([]P, error) {
	ps := make([]P, 0, 0)
	tx := r.DB().Where("id in ?", ids)
	for _, preload := range preloads {
		if pre, ok := preload.(string); ok && pre != "" {
			tx = tx.Preload(pre)
		} else if pre, ok := preload.(base.Preload); ok && pre.Query != "" {
			tx = tx.Preload(pre.Query, pre.Args...)
		}
	}
	err := tx.Find(&ps).Error
	return ps, err
}

func (r *Repository[P]) GetByID(id string, preloads ...interface{}) (P, error) {
	po := new(P)
	tx := r.DB()
	for _, preload := range preloads {
		if pre, ok := preload.(string); ok && pre != "" {
			tx = tx.Preload(pre)
		} else if pre, ok := preload.(base.Preload); ok && pre.Query != "" {
			tx = tx.Preload(pre.Query, pre.Args...)
		}
	}
	err := tx.First(po, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return *po, exception.DataNotFound
	}
	return *po, err
}

func (r *Repository[P]) FindExport(query string, args []interface{}, preloads []interface{}, order ...string) ([]P, error) {
	ps := make([]P, 0, 0)
	tx := r.DB().Where(query, args...)
	for _, preload := range preloads {
		if pre, ok := preload.(string); ok && pre != "" {
			tx = tx.Preload(pre)
		} else if pre, ok := preload.(base.Preload); ok && pre.Query != "" {
			tx = tx.Preload(pre.Query, pre.Args...)
		}
	}
	if len(order) > 0 && order[0] != "" {
		tx = tx.Order(order[0])
	} else {
		tx = tx.Order("updated_at desc")
	}
	err := tx.Find(&ps).Error
	return ps, err
}

func (r *Repository[P]) PageExport(limit, offset int, query string, args []interface{}, preloads []interface{}, order ...string) ([]P, int64, error) {
	var count int64
	ps := make([]P, 0, 0)
	tx := r.DB().Model(&ps).Where(query, args...).Count(&count)
	for _, preload := range preloads {
		if pre, ok := preload.(string); ok && pre != "" {
			tx = tx.Preload(pre)
		} else if pre, ok := preload.(base.Preload); ok && pre.Query != "" {
			tx = tx.Preload(pre.Query, pre.Args...)
		}
	}
	if len(order) > 0 {
		tx = tx.Order(order[0])
	} else {
		tx = tx.Order("updated_at desc")
	}
	err := tx.Limit(limit).Offset(offset).Find(&ps).Error
	return ps, count, err
}

func (r *Repository[P]) All(preloads ...interface{}) ([]P, error) {
	ps := make([]P, 0, 0)
	tx := r.DB().Model(&ps)
	for _, preload := range preloads {
		if pre, ok := preload.(string); ok && pre != "" {
			tx = tx.Preload(pre)
		} else if pre, ok := preload.(base.Preload); ok && pre.Query != "" {
			tx = tx.Preload(pre.Query, pre.Args...)
		}
	}
	err := tx.Order("updated_at desc").Find(&ps).Error
	return ps, err
}

func (r *Repository[P]) Create(po P) (string, error) {
	if cpo, ok := any(po).(base.CPO); ok {
		parentID := cpo.GetParentID()
		if len(parentID) != 0 {
			parent, err := r.GetByID(parentID)
			if err != nil {
				return "", err
			}
			parentIDs := any(parent).(base.CPO).GetParentIDs()
			if len(parentIDs) > 0 {
				cpo.SetParentIDs(fmt.Sprintf("%s,%s", parentIDs, parentID))
			} else {
				cpo.SetParentIDs(parentID)
			}
		}
	}

	tx := r.DB().Create(po)
	if rpo, ok := any(po).(base.RPO); ok {
		return rpo.GetID(), tx.Error
	}
	return "", tx.Error
}

func (r *Repository[P]) Update(po P) (int64, error) {
	var oldParentIDs, newParentIDs string
	cpo, ok := any(po).(base.CPO)
	if ok && cpo.ParentChanged() {
		oldPO, err := r.GetByID(cpo.GetID())
		if err != nil {
			return 0, err
		}
		oldCPO := any(oldPO).(base.CPO)
		parentID := cpo.GetParentID()
		oldParentID := oldCPO.GetParentID()
		if oldParentID != parentID {
			if len(oldParentID) == 0 {
				oldParentIDs = cpo.GetID()
			} else {
				oldParentIDs = fmt.Sprintf("%s,%s", oldCPO.GetParentIDs(), cpo.GetID())
			}

			if len(parentID) == 0 {
				cpo.SetParentIDs("")
				newParentIDs = cpo.GetID()
			} else {
				parent, err := r.GetByID(parentID)
				if err != nil {
					return 0, err
				}
				parentIDs := any(parent).(base.CPO).GetParentIDs()
				if len(parentIDs) > 0 {
					cpo.SetParentIDs(fmt.Sprintf("%s,%s", parentIDs, parentID))
				} else {
					cpo.SetParentIDs(parentID)
				}
				newParentIDs = fmt.Sprintf("%s,%s", cpo.GetParentIDs(), cpo.GetID())
			}
		}
	}

	var count int64
	err := r.Transaction(func() error {
		updateValues := po.GetChanges()
		if len(updateValues) == 0 {
			return nil
		}
		tx := r.DB().Model(po).Updates(updateValues)
		if tx.Error != nil {
			return tx.Error
		}
		count = tx.RowsAffected

		if oldParentIDs == newParentIDs {
			return nil
		}

		var ps []P
		tx = r.DB().Model(&ps).Where("find_in_set(?,parent_ids)", cpo.GetID()).
			Update("parent_ids", gorm.Expr("replace(`parent_ids`, ?, ?)", oldParentIDs, newParentIDs))
		count += tx.RowsAffected
		return tx.Error
	})

	return count, err
}

func (r *Repository[P]) UnscopedDeleteByID(ids ...string) (count int64, err error) {
	return r.delete(true, ids...)

}

func (r *Repository[P]) DeleteByID(ids ...string) (count int64, err error) {
	return r.delete(false, ids...)
}

func (r *Repository[P]) delete(unscoped bool, ids ...string) (count int64, err error) {
	if len(ids) == 0 {
		return
	}
	e := new(P)
	tx := r.DB()
	if unscoped {
		tx = tx.Unscoped()
	}
	if _, ok := any(*e).(base.CPO); ok {
		for _, id := range ids {
			tx = tx.Or("id = ? or find_in_set(?,parent_ids)", id, id)
		}
		tx.Delete(e)
		count = tx.RowsAffected
		err = tx.Error
	} else {
		tx = tx.Delete(e, ids)
		count = tx.RowsAffected
		err = tx.Error
	}
	return
}

func (r *Repository[P]) DeleteAssociationByID(association string, ids ...string) error {
	es, err := r.getSliceByID(ids...)
	if err != nil || es == nil {
		return err
	}
	return r.DB().Model(es).Association(association).Clear()
}

func (r *Repository[P]) ClearAutoUpdateAssociationByID(ids ...string) error {
	if len(ids) == 0 {
		return nil
	}
	var p P
	if _, ok := any(p).(base.RPO); !ok {
		return nil
	}
	associations := base.GetUpdateAssociation(reflect.New(reflect.TypeOf(p).Elem()).Interface().(base.PO))
	if len(associations) == 0 {
		return nil
	}
	es, err := r.getSliceByID(ids...)
	if err != nil || es == nil {
		return err
	}
	err = r.Transaction(func() error {
		for _, association := range associations {
			if err := r.DB().Model(es).Association(association).Clear(); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *Repository[P]) ClearAutoDeleteAssociationByID(ids ...string) error {
	if len(ids) == 0 {
		return nil
	}
	var p P
	if _, ok := any(p).(base.RPO); !ok {
		return nil
	}
	associations := base.GetDeleteAssociation(reflect.New(reflect.TypeOf(p).Elem()).Interface().(base.PO))
	if len(associations) == 0 {
		return nil
	}
	es, err := r.getSliceByID(ids...)
	if err != nil || es == nil {
		return err
	}
	err = r.Transaction(func() error {
		for _, association := range associations {
			if err := r.DB().Model(es).Association(association).Clear(); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *Repository[P]) getSliceByID(ids ...string) (interface{}, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var p P
	if _, ok := any(p).(base.RPO); !ok {
		return nil, nil
	}
	pType := reflect.TypeOf(p).Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(pType), 0, len(ids))
	for _, id := range ids {
		e := reflect.New(pType).Interface().(base.RPO)
		e.SetID(id)
		slice = reflect.Append(slice, reflect.ValueOf(e).Elem())
	}
	if _, ok := any(p).(base.CPO); ok {
		child := make([]string, 0, 10)
		tx := r.DB().Model(p).Select("id")
		for _, id := range ids {
			tx = tx.Or("find_in_set(?,parent_ids)", id)
		}
		if err := tx.Find(&child).Error; err != nil {
			return nil, err
		}
		for _, id := range child {
			e := reflect.New(pType).Interface().(base.RPO)
			e.SetID(id)
			slice = reflect.Append(slice, reflect.ValueOf(e).Elem())
		}
	}
	es := slice.Interface()
	return es, nil
}
func (r *Repository[P]) CheckIndex(po P) error {
	indices := base.GetIndexes(po)
	if len(indices) == 0 {
		return nil
	}
	for _, index := range indices {
		tx := r.DB().Model(po)
		if rpo, ok := any(po).(base.RPO); ok && len(rpo.GetID()) > 0 {
			tx.Where("id != ?", rpo.GetID())
		}
		for i, column := range index.Columns {
			tx.Where(fmt.Sprintf("%s = ?", column), index.Values[i])
		}
		var count int64
		if err := tx.Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			continue
		}
		return index.Err
	}
	return nil
}

func (r *Repository[P]) CheckVerify(po P, ids ...string) error {
	verify := base.GetVerify(po)
	if verify == nil {
		return nil
	}

	var tx *gorm.DB
	if len(ids) == 0 {
		tx = r.DB().Model(po)
		if rpo, ok := any(po).(base.RPO); ok {
			tx.Where("id = ?", rpo.GetID())
		}

	} else {
		tx = r.DB().Model(new(P)).Where("id in ?", ids)
	}

	for i := range verify.Columns {
		tx.Where(verify.Columns[i], verify.Values[i])
	}
	var count int64
	if err := tx.Count(&count).Error; err != nil {
		return err
	}

	if len(ids) == 0 {
		if count != 1 {
			return exception.DataInvalid
		}
	} else {
		if count != int64(len(ids)) {
			return exception.DataInvalid
		}
	}
	return nil
}
