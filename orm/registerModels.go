package orm

import (
	"fmt"
	"reflect"
)

var (
	TypeTable = map[string]reflect.Type{
		/*
			Регистрация таблиц
			"<tablename>":    reflect.TypeOf(<ModelStructName>),
		*/
	}
	ContstructMap = map[string]func(string) interface{}{
		/*
			Регистрация конструкторов
			"<tablename>":   New<ModelStructName>, */
	}
)

func RegisterModel(tableName string, obj interface{}) error {

	v := reflect.ValueOf(obj)

	// Проверяем, что переданный объект — это структура
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("object %v is not a struct", v)
	}
	TypeTable[tableName] = reflect.TypeOf(obj)

	constructor := func(tableName string) interface{} {
		objValue := reflect.New(TypeTable[tableName])
		obj := objValue.Elem()
		obj.FieldByName("TableName").SetString(tableName)
		return &obj
	}
	ContstructMap[tableName] = constructor
	return nil
}

func ExtractFields(obj interface{}) ([]interface{}, []string) {
	// Используем reflect.Indirect для получения значения структуры
	val := reflect.Indirect(reflect.ValueOf(obj))
	typ := val.Type() // Получаем тип значения, после преобразования

	var values []interface{}
	var columns []string

	// Проходим по полям структуры
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		// Добавляем имя поля в список столбцов
		columns = append(columns, fieldName)

		// Добавляем значение поля в список значений
		values = append(values, val.Field(i).Interface())
	}
	return values, columns
}
