package mate

import "github.com/ericnts/cherry/base"

type Service[R Repo[E], E base.PO] struct {
	Resource

	Repo R
}

func (r *Service[R, E]) Create(e E) (int64, error) {
	if err := r.Repo.CheckIndex(e); err != nil {
		return 0, err
	}
	return r.Repo.Create(e)
}

func (r *Service[R, E]) FindByQuery(q base.Query) ([]E, error) {
	query, args := q.Query()
	if len(query) == 0 {
		query, args = base.GetQuery(q)
	}
	return r.Repo.FindExport(query, args, q.GetPreloads(), q.GetOrder())
}

func (r *Service[R, E]) FindByPage(q base.PageQuery) ([]E, int64, error) {
	query, args := q.Query()
	if len(query) == 0 {
		query, args = base.GetQuery(q)
	}
	return r.Repo.PageExport(q.Limit(), q.Offset(), query, args, q.GetPreloads(), q.GetOrder())
}

func (r *Service[R, E]) GetByID(id string, preloads ...interface{}) (E, error) {
	return r.Repo.GetByID(id, preloads...)
}

func (r *Service[R, E]) Update(e E) (count int64, err error) {
	if err = r.Repo.CheckVerify(e); err != nil {
		return
	}
	if err = r.Repo.CheckIndex(e); err != nil {
		return
	}
	err = r.Transaction(func() error {
		if rpo, ok := any(e).(base.RPO); ok {
			if err = r.Repo.ClearAutoUpdateAssociationByID(rpo.GetID()); err != nil {
				return err
			}
		}
		count, err = r.Repo.Update(e)
		return err
	})
	return
}

func (r *Service[R, E]) Delete(e E, ids ...string) (count int64, err error) {
	idMap := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if len(id) > 0 {
			idMap[id] = struct{}{}
		}
	}
	if rpo, ok := any(e).(base.RPO); ok {
		if id := rpo.GetID(); len(id) > 0 {
			idMap[id] = struct{}{}
		}
	}
	newIDs := make([]string, 0, len(idMap))
	for key := range idMap {
		newIDs = append(newIDs, key)
	}
	if len(newIDs) == 0 {
		return 0, nil
	}
	if err = r.Repo.CheckVerify(e, newIDs...); err != nil {
		return
	}
	err = r.Transaction(func() error {
		if err = r.Repo.ClearAutoDeleteAssociationByID(newIDs...); err != nil {
			return err
		}
		count, err = r.Repo.DeleteByID(newIDs...)
		return err
	})
	return
}
