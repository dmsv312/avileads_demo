package repository

import (
	"avileads-web/models"
	"avileads-web/repository/audit"
	helpers "avileads-web/repository/helpers"
	repository "avileads-web/repository/specification"
	"fmt"
	"github.com/astaxie/beego/orm"
	"reflect"
)

type BaseRepository[T any] struct {
	db      orm.Ormer
	meta    *audit.AuditMeta
	auditor *BaseRepository[models.AuditLog]
}

func New[T any](db orm.Ormer, meta *audit.AuditMeta) *BaseRepository[T] {
	return &BaseRepository[T]{db: db, meta: meta}
}

func (r *BaseRepository[T]) Create(model *T) (int64, error) {
	id, err := r.db.Insert(model)
	if err != nil {
		return 0, err
	}

	err = r.writeAudit("create", model, nil)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *BaseRepository[T]) Get(id any, relations ...interface{}) (*T, error) {
	var entity T
	qs := r.db.QueryTable(&entity).Filter("Id", id)

	if len(relations) > 0 {
		qs = qs.RelatedSel(relations...)
	}

	err := qs.One(&entity)
	return &entity, err
}

func (r *BaseRepository[T]) Update(model *T, cols ...string) error {
	oldEntity := new(T)
	err := r.db.QueryTable(new(T)).
		Filter("Id", helpers.GetPK(model)).
		One(oldEntity)

	if err != nil {
		return err
	}

	if _, err := r.db.Update(model, cols...); err != nil {
		return err
	}

	err = r.writeAudit("update", model, oldEntity)
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepository[T]) Delete(id any) error {
	oldEntity := new(T)
	err := r.db.QueryTable(new(T)).Filter("Id", id).One(oldEntity)
	if err != nil {
		return err
	}

	entity := new(T)
	reflect.ValueOf(entity).Elem().FieldByName("Id").Set(reflect.ValueOf(id))
	if _, err := r.db.Delete(entity); err != nil {
		return err
	}

	err = r.writeAudit("delete", nil, oldEntity)
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepository[T]) GetAll(specification repository.Specification, orderBy ...string) ([]T, error) {
	querySetter := r.db.QueryTable(new(T))

	if specification != nil {
		if cond := specification.Cond(); cond != nil {
			querySetter = querySetter.SetCond(cond)
		}
	}

	if len(orderBy) > 0 {
		querySetter = querySetter.OrderBy(orderBy...)
	}

	var entities []T
	_, err := querySetter.All(&entities)
	return entities, err
}

func (r *BaseRepository[T]) UpdateAll(condition *orm.Condition, setters orm.Params) (int64, error) {
	var oldEntities []T
	qs := r.db.QueryTable(new(T))
	if condition != nil {
		qs = qs.SetCond(condition)
	}

	_, err := qs.All(&oldEntities)
	if err != nil {
		return 0, err
	}

	affected, err := qs.Update(setters)
	if err != nil || affected == 0 {
		return affected, err
	}

	var newEntities []T
	_, err = qs.All(&newEntities)
	if err != nil {
		return affected, err
	}

	oldMap := make(map[int]*T, len(oldEntities))
	for i := range oldEntities {
		pk := helpers.GetPK(&oldEntities[i]).(int)
		oldMap[pk] = &oldEntities[i]
	}

	for i := range newEntities {
		pk := helpers.GetPK(&newEntities[i]).(int)
		err := r.writeAudit("update", &newEntities[i], oldMap[pk])
		if err != nil {
			return affected, err
		}
	}

	return affected, nil
}

func (r *BaseRepository[T]) DeleteAll(condition *orm.Condition) (int64, error) {
	querySetter := r.db.QueryTable(new(T))
	if condition != nil {
		querySetter = querySetter.SetCond(condition)
	}
	return querySetter.Delete()
}

func (r *BaseRepository[T]) First(specification repository.Specification) (*T, error) {
	var entity T
	err := r.db.QueryTable(new(T)).
		SetCond(specification.Cond()).
		One(&entity)
	return &entity, err
}

func (r *BaseRepository[T]) DeleteBatch(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	var oldEntities []T
	_, _ = r.db.QueryTable(new(T)).Filter("Id__in", ids).All(&oldEntities)

	affected, err := r.db.
		QueryTable(new(T)).
		Filter("Id__in", ids).
		Delete()

	if err != nil || affected == 0 {
		return affected, err
	}

	for i := range oldEntities {
		err := r.writeAudit("delete", nil, &oldEntities[i])
		if err != nil {
			return affected, err
		}
	}

	return affected, nil
}

func (r *BaseRepository[T]) writeAudit(action string, newEntity, oldEntity interface{}) error {
	diffOrSnapshot := helpers.MarshalChanges(action, newEntity, oldEntity)

	if action == "update" && diffOrSnapshot == "{}" {
		return nil
	}

	var actionID int
	switch action {
	case "create":
		actionID = helpers.ActionCreate
	case "update":
		actionID = helpers.ActionUpdate
	case "delete":
		actionID = helpers.ActionDelete
	default:
		return fmt.Errorf("unknown audit action: %s", action)
	}

	modelName := helpers.TypeName[T]()
	modelID, err := helpers.GetModelID(r.db, modelName)
	if err != nil {
		return err
	}

	if action == "create" && r.meta.RkId == 0 && modelName == "Vacancy_avito" {
		r.meta.RkId = helpers.PkFromEntity(newEntity, oldEntity)
	}

	_, err = r.db.Insert(&models.AuditLog{
		UserId:   r.meta.UserId,
		RkId:     r.meta.RkId,
		ModelId:  modelID,
		ActionId: actionID,
		EntityId: helpers.PkFromEntity(newEntity, oldEntity),
		Changes:  diffOrSnapshot,
	})

	return err
}
