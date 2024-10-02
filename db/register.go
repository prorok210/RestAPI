package db

import (
	"fmt"
	"reflect"
)

func registerModel(tableName string, obj interface{}) error {

	v := reflect.ValueOf(obj)

	// Проверяем, что переданный объект — это структура
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("object %v is not a struct", v)
	}
	TableRegistry[tableName] = reflect.TypeOf(obj)
	TypeMap[tableName] = reflect.TypeOf(obj).Elem()

	constructor := func(tableName string) interface{} {
		objValue := reflect.New(TableRegistry[tableName])
		obj := objValue.Elem()
		obj.FieldByName("TableName").SetString(tableName)
		return &obj
	}
	ContstructMap[tableName] = constructor
	return nil
}
