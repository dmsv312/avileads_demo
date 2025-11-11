package repository

import (
	"avileads-web/models"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"reflect"
)

const (
	ActionCreate = 1
	ActionUpdate = 2
	ActionDelete = 3
)

func MarshalChanges(action string, newEntity, oldEntity interface{}) string {
	switch action {
	case "update":
		changes := make(map[string][2]interface{})

		newValue := reflect.ValueOf(newEntity).Elem()
		oldValue := reflect.ValueOf(oldEntity).Elem()
		entityType := newValue.Type()

		for fieldIndex := 0; fieldIndex < entityType.NumField(); fieldIndex++ {
			structField := entityType.Field(fieldIndex)
			if !structField.IsExported() || structField.Name == "Id" {
				continue
			}

			newComparable := comparableValueOf(newValue.Field(fieldIndex))
			oldComparable := comparableValueOf(oldValue.Field(fieldIndex))

			if !reflect.DeepEqual(newComparable, oldComparable) {
				changes[structField.Name] = [2]interface{}{oldComparable, newComparable}
			}
		}

		if len(changes) == 0 {
			return "{}"
		}
		jsonBytes, _ := json.Marshal(changes)
		return string(jsonBytes)

	default:
		jsonBytes, _ := json.Marshal(newEntity)
		if newEntity == nil {
			jsonBytes, _ = json.Marshal(oldEntity)
		}
		return string(jsonBytes)
	}
}

func TypeName[T any]() string {
	var zero T
	return reflect.TypeOf(zero).Name()
}

func GetPK(entity interface{}) interface{} {
	return reflect.ValueOf(entity).Elem().FieldByName("Id").Interface()
}

func PkFromEntity(newEntity, oldEntity interface{}) int {
	if newEntity != nil {
		return GetPK(newEntity).(int)
	}
	return GetPK(oldEntity).(int)
}

func GetModelID(db orm.Ormer, modelName string) (int, error) {
	var auditModel models.AuditModel

	err := db.QueryTable(new(models.AuditModel)).
		Filter("name", modelName).
		One(&auditModel)
	if err != nil && err != orm.ErrNoRows {
		return 0, err
	}

	if auditModel.Id == 0 {
		auditModel.Name = modelName
		id, err := db.Insert(&auditModel)

		if err != nil {
			return 0, err
		}

		return int(id), nil
	}

	return auditModel.Id, nil
}

func comparableValueOf(fieldValue reflect.Value) interface{} {
	if !fieldValue.IsValid() {
		return nil
	}
	if fieldValue.Kind() == reflect.Pointer {
		if fieldValue.IsNil() {
			return nil
		}
		idField := fieldValue.Elem().FieldByName("Id")
		if idField.IsValid() {
			return idField.Interface()
		}
		return fieldValue.Pointer()
	}
	return fieldValue.Interface()
}
