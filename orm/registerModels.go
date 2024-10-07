package orm

import (
	"fmt"
	"reflect"
)

var (
	TypeTable = map[string]reflect.Type{
		/*
			Регистрация таблиц. Происходит автоматически
			"<tablename>":    reflect.TypeOf(<ModelStructName>),
		*/
	}
	ContstructMap = map[string]func(string) (interface{}, error){
		/*
			Регистрация конструкторов. Происходит автоматически
			"<tablename>":   New<ModelStructName>, */
	}
)

func RegisterModel(tableName string, obj interface{}) error {

	v := reflect.ValueOf(obj)

	// Проверяем, что переданный объект — это структура
	// Проверяем, что переданный объект — это структура или указатель на структуру
	if v.Kind() != reflect.Struct && !(v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct) {
		return fmt.Errorf("object %v is not a struct or pointer to a struct", v)
	}
	// Сохраняем тип объекта в TypeTable по ключу - названию таблицы
	// Если это указатель, берём тип структуры, на которую он указывает
	structType := reflect.TypeOf(obj)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	TypeTable[tableName] = structType

	// Создаем конструктор для объекта и сохраняем его в ContstructMap по ключу - названию таблицы
	constructor := func(tableName string) (interface{}, error) {
		objValue := reflect.New(TypeTable[tableName]).Elem()
		// Устанавливаем поле "TableName", если оно существует
		field := objValue.FieldByName("TableName")
		if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
			field.SetString(tableName)
		} else {
			return nil, fmt.Errorf("field TableName is invalid or cannot be set")
		}

		// Возвращаем указатель на новый объект
		return objValue.Addr().Interface(), nil
	}
	ContstructMap[tableName] = constructor
	return nil
}
