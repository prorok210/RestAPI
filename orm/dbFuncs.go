package orm

import (
	"RestAPI/core"
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
)

func InitDB() error {
	var InitDBError error
	conn, InitDBError = pgx.Connect(context.Background(), core.CONNECTIONDATA)
	if InitDBError != nil {
		return fmt.Errorf("database connect error: %v", InitDBError)
	}

	// Проверка наличия таблиц в базе данных и их соответствия моделям
	err := CheckTables()
	if err != nil {
		return fmt.Errorf("error checking tables: %v", err)
	}

	log.Println("Successfully connected to the database.")

	return nil
}

func CreateTable(tableName string, obj interface{}) error {
	data := reflect.TypeOf(obj)

	// Проверка на то, что таблица уже существует
	var exists bool

	// Делаем пробный запрос к БД
	query := "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name=$1);"
	err := conn.QueryRow(context.Background(), query, tableName).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking table existence: %v", err)
	}

	// Если получили ответ, значит, таблица уже существует
	if exists {
		fmt.Printf("Table %s already exists.", tableName)
		return nil
	}

	// Начало SQL-запроса для создания таблицы
	sqlQuery := "CREATE TABLE IF NOT EXISTS " + tableName + " ("

	for i := 0; i < data.NumField(); i++ {
		field := data.Field(i)
		// Поля TableName в БД не существует
		if field.Name == "TableName" {
			continue
		}
		// Парсим тег orm, который описывает поле таблицы в БДы
		ormTag := field.Tag.Get("orm")
		ormTag = strings.Replace(ormTag, "_", " ", -1)
		if ormTag == "" {
			return fmt.Errorf("field %s does not have a tag", field.Name)
		} else if strings.Contains(ormTag, "ref") {
			// Поиск начала подстроки ref
			start := strings.Index(ormTag, "ref")

			// Т,к ref идет в конце, берем подстроку от ref до конца строки
			match := ormTag[start:]
			// Удаляем ref из ormTag
			ormTag = strings.Replace(ormTag, " "+match, "", -1)
			// Создаем sql-запрос с внешним ключом
			sqlQuery += strings.ToLower(field.Name) + " " + ormTag + ", " + "FOREIGN KEY (" + strings.ToLower(field.Name) + ") REFERENCES " + strings.TrimPrefix(match, "ref ") + ");"
		} else {
			// Создаем sql-запрос
			sqlQuery += strings.ToLower(field.Name) + " " + ormTag + ");"
		}
	}

	// Делаем запрос к БД на создание таблицы
	_, err = conn.Exec(context.Background(), sqlQuery)
	if err != nil {
		return fmt.Errorf("error creating table %s", err)
	} else {
		fmt.Println("Table created successfully")
	}
	return nil
}

// Функция проверки соответствия таблиц в базе данных и моделей
func CheckTables() error {
	for tableName, modelType := range TypeTable {
		fmt.Printf("Checking table %s...\n", tableName)

		// Получение текущей структуры таблицы из базы данных
		dbColumns, err := getTableColumns(tableName)

		if err != nil {
			return fmt.Errorf("error getting table columns %s: %v", tableName, err)
		}
		// fmt.Println("dbColumns", dbColumns)

		// Сравнение конструкции с моделью
		modelColumns, err := getModelColumns(modelType)

		if err != nil {
			return fmt.Errorf("error getting model fields %s: %v", modelType.Name(), err)
		}
		// fmt.Println("modelColumns", modelColumns)

		// Сравнение столбцов
		ok, err := compareColumns(dbColumns, modelColumns)
		if err != nil {
			return fmt.Errorf("error comparing columns %s: %v", tableName, err)
		}
		if ok {
			fmt.Printf("Table %s matches the model\n", tableName)
		} else {
			fmt.Printf("Difference in table structure %s!\n", tableName)
		}
	}

	return nil
}

// Получение списка столбцов таблицы из базы данных
func getTableColumns(tableName string) (map[string]map[string]string, error) {
	// использование служебных таблиц для получения метаданных столбца
	query := fmt.Sprintf(`
	SELECT
    	c.column_name,
    	c.data_type,
    	c.is_nullable,
    	tc.constraint_type,
    	rc.update_rule,
    	rc.delete_rule
	FROM
    	information_schema.columns AS c
	LEFT JOIN 
		information_schema.key_column_usage AS kcu 
		ON c.table_name = kcu.table_name AND c.column_name = kcu.column_name
	LEFT JOIN 
		information_schema.table_constraints AS tc
		ON kcu.constraint_name = tc.constraint_name AND kcu.table_name = tc.table_name
	LEFT JOIN 
		information_schema.referential_constraints AS rc 
		ON tc.constraint_name = rc.constraint_name
	WHERE 
		c.table_name = '%s';`, tableName)

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer rows.Close()

	columns := make(map[string]map[string]string)
	for rows.Next() {
		var columnName, dataType, isNullable, constraintType, updateRule, deleteRule sql.NullString

		err := rows.Scan(&columnName, &dataType, &isNullable, &constraintType, &updateRule, &deleteRule)
		if err != nil {
			return nil, fmt.Errorf("line reading error: %v", err)
		}

		columnDetails := map[string]string{
			"column_name":     columnName.String,
			"data_type":       dataType.String,
			"is_nullable":     isNullable.String,
			"constraint_type": constraintType.String,
			"update_rule":     updateRule.String,
			"delete_rule":     deleteRule.String,
		}

		columns[columnName.String] = columnDetails
	}

	return columns, nil
}

// Получение данных об атрибутах модели
func getModelColumns(modelType reflect.Type) (map[string]map[string]string, error) {
	columns := map[string]map[string]string{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		// Дефолтные значения
		columns[strings.ToLower(field.Name)] = map[string]string{
			"column_name":     strings.ToLower(field.Name),
			"data_type":       "",
			"is_nullable":     "YES",
			"constraint_type": "",
			"update_rule":     "",
			"delete_rule":     ""}
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		column := columns[strings.ToLower(field.Name)]

		ormTag := strings.ToLower(field.Tag.Get("orm"))
		if ormTag != "" {
			ormTagSlice := strings.Split(ormTag, " ")
			// Первый атрибут - всегда тип. Т,к SQL свои типы, нужно их привести к общему через tagToSqlType. Его нужно дополнять
			value, ok := tagToSqlType[ormTagSlice[0]]
			if !ok {
				if ormTagSlice[0] == "string" {
					return nil, fmt.Errorf("unknown type %s. May be you meant varchar?", ormTagSlice[0])
				}
				return nil, fmt.Errorf("unknown type %s. You can add it in tagToSqlType in constants.go", ormTagSlice[0])
			} else if ok {
				column["data_type"] = value
			}
			// Заполняем поля в зависимости от тегов
			if contains(ormTagSlice, "serial_primary_key") {
				// Если есть тег serial_primary_key, то это первичный ключ, он является идентификатором и ненулевой
				column["constraint_type"] = "PRIMARY KEY"
				column["is_identity"] = "true"
				column["is_nullable"] = "NO"
			}
			if contains(ormTagSlice, "ref") {
				// Если есть тег ref, то это внешний ключ, он является идентификатором, но может быть нулевым (забавно, но sql возвращает именно так)
				column["constraint_type"] = "FOREIGN KEY"
				column["is_identity"] = "true"
				// Далее идет добавление правил обновления и удаления
				if contains(ormTagSlice, "on_update") {
					updateRule := findRule(ormTag, "on_update")
					column["update_rule"] = updateRule
				}
				if contains(ormTagSlice, "on_delete") {
					fmt.Println("on_delete contains")
					deleteRule := findRule(ormTag, "on_delete")
					column["delete_rule"] = deleteRule
				}
			}
			// Затем добавляем остальные поля. Всё приводим к формату, который возвращает SQL для сравнения модели и таблицы
			if contains(ormTagSlice, "not_null") {
				column["is_nullable"] = "NO"
			}
			if contains(ormTagSlice, "unique") {
				column["constraint_type"] = "UNIQUE"
			}

		}
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("no fields found in model")
	}

	return columns, nil
}

func contains(slice []string, input string) bool {
	for _, str := range slice {
		if strings.Contains(str, input) {
			return true
		}
	}
	return false
}

func findRule(tag string, rule string) string {
	// Найти правило (удаления или обновления) в теге
	start := strings.Index(tag, rule)
	end := start + strings.Index(tag[start:], " ")
	// Если правило находится в конце строки т,е не был найден пробел т,к разделителями являются '_'
	if end == start-1 {
		end = len(tag)
	}
	result := strings.ToUpper(strings.Replace(tag[start:end], "_", " ", -1))
	// Если не задано действие (Т,е есть слова ON UPDATE или ON DELETE, но нет действия по типу CASCADE, SET NULL и т.д)
	if len(strings.Split(result, " ")) < 3 {
		return ""
	}
	// Для случая, когда правило из нескольких слов. Например, NO ACTION или SET NULL
	return strings.Join(strings.Split(result, " ")[2:], "")
}

// Сравнение характеристик таблиц в БД и модели
func compareColumns(dbColumns, modelColumns map[string]map[string]string) (bool, error) {
	// Проверка наличия всех столбцов таблицы в модели
	for key, valueDb := range dbColumns {
		valueModel, ok := modelColumns[key]
		if !ok {
			return false, fmt.Errorf("column %s not found in model", key)
		}

		// Проверка соответствия столбцов
		for key, value := range valueDb {
			// SQL возвращает то пустое значение, то NO ACTION
			if key == "update_rule" || key == "delete_rule" {
				if (value == "NO ACTION" || value == "") && (valueModel[key] == "" || valueModel[key] == "NO ACTION") {
					continue
				}
			}
			if value != valueModel[key] {
				return false, fmt.Errorf(`in column "%s" %s does not match the model, db: %s and model: %s`, valueModel["column_name"], key, value, valueModel[key])
			}
		}
	}

	return true, nil
}
