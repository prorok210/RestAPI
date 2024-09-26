package db

import "reflect"

var (
	TableRegistry = map[string]reflect.Type{
		"users": reflect.TypeOf(User{TableName: "users"}),
		/*
			Регистрация таблиц
			"<tablename>":    reflect.TypeOf(<ModelStructName>{TableName: "<tablename>"}),
		*/

	}

	TypeMap = map[string]reflect.Type{
		"users": reflect.TypeOf((*User)(nil)).Elem(),
		/*
			Регистрация типов
			"<tablename>":   reflect.TypeOf((*<ModelStructName>)(nil)).Elem(),
		*/

	}
)

// Структура модели
/*
func NewConstructor(<TableName> string, fields []string) *<ModelStruct> {
	return &<ModelStruct>{
		TableName: <TableName>,
		<field1>:   fields[0],
		<field2>:   fields[1],
		...
	}
}
*/

/*
func (obj *<ModelStruct>) ToFields() ([]interface{}, []string) {
	return extractFields(<ModelStruct>)
}
*/

func NewUser(tableName string, id int, mobile string, otp string, isActive bool, name string, surname string, age int, email string, password string, createAt string, updateAt string) *User {
	return &User{
		TableName: tableName,
		ID:        id,
		Mobile:    mobile,
		Otp:       otp,
		isActive:  isActive,
		Name:      name,
		Surname:   surname,
		Age:       age,
		Email:     email,
		Password:  password,
		createAt:  createAt,
		updateAt:  updateAt,
	}
}

func (obj User) ToFields() ([]interface{}, []string) {
	return ExtractFields(obj)
}

// Generic function to extract fields
func ExtractFields(obj interface{}) ([]interface{}, []string) {
	// Getting the value and type of the object
	val := reflect.ValueOf(obj).Elem()
	typ := reflect.TypeOf(obj).Elem()

	var values []interface{}
	var columns []string

	// We go to fields of the structure
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		// Adding a field name to the list of columns
		columns = append(columns, fieldName)

		// Adding a field value to the list of values
		values = append(values, val.Field(i).Interface())
	}
	return values, columns
}
