package db

import "reflect"

var (
	TableRegistry = map[string]reflect.Type{
		/*
			Регистрация таблиц
			"<tablename>":    reflect.TypeOf(<ModelStructName>{TableName: "<tablename>"}),
		*/
	}

	TypeMap = map[string]reflect.Type{
		/*
			Регистрация типов
			"<tablename>":   reflect.TypeOf((*<ModelStructName>)(nil)).Elem(),
		*/
	}
	ContstructMap = map[string]func(string) interface{}{
		/*
			Регистрация конструкторов
			"<tablename>":   New<ModelStructName>, */
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
		IsActive:  isActive,
		Name:      name,
		Surname:   surname,
		Age:       age,
		Email:     email,
		Password:  password,
		CreateAt:  createAt,
		UpdateAt:  updateAt,
	}
}

func NewToken(tableName string, id int, userID string, accessTokenHash string, refreshTokenHash string, expiredAt string, createAt string, updateAt string) *Token {
	return &Token{
		TableName:        tableName,
		ID:               id,
		UserID:           userID,
		AccessTokenHash:  accessTokenHash,
		RefreshTokenHash: refreshTokenHash,
		ExpiredAt:        expiredAt,
		CreateAt:         createAt,
		UpdateAt:         updateAt,
	}
}

func NewImage(tableName string, id int, userID string, imageURL string) *Image {
	return &Image{
		TableName: tableName,
		ID:        id,
		UserID:    userID,
		ImageURL:  imageURL,
	}
}

func (obj User) ToFields() ([]interface{}, []string) {
	return ExtractFields(obj)
}

func (obj Token) ToFields() ([]interface{}, []string) {
	return ExtractFields(obj)
}

func (obj Image) ToFields() ([]interface{}, []string) {
	return ExtractFields(obj)
}

// Generic function to extract fields
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
