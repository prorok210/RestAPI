package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5"
)

func (user *User) ToFields() ([]interface{}, []string) {
	return extractFields(user)
}

func (dialog *Dialog) ToFields() ([]interface{}, []string) {
	return extractFields(dialog)
}

func (message *Message) ToFields() ([]interface{}, []string) {
	return extractFields(message)
}

// Generic function to extract fields
func extractFields(obj interface{}) ([]interface{}, []string) {
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

func InitDB() error {
	var InitDBError error
	conn, InitDBError = pgx.Connect(context.Background(), CONNECTIONDATA)
	if InitDBError != nil {
		return fmt.Errorf("database connect error: %v", InitDBError)
	}

	err := CheckTables()
	if err != nil {
		return fmt.Errorf("error checking tables: %v", err)
	}

	fmt.Println("Successfully connected to the database.")

	creatingTables := []string{}

	createTableSQL1 := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50),
		email VARCHAR(50)
	);`

	createTableSQL2 := `
	CREATE TABLE IF NOT EXISTS dialogs (
		id SERIAL PRIMARY KEY,
		owner VARCHAR(50),
		opponent VARCHAR(50)
	);`

	creatingTables = append(creatingTables, createTableSQL1, createTableSQL2)
	// Execute a table creation request
	for _, createTableSQL := range creatingTables {
		_, err := conn.Exec(context.Background(), createTableSQL)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		} else {
			fmt.Println("Table created successfully or already exists.")
		}
	}

	return nil
}

func CreateTable(obj interface{}) error {
	data := reflect.TypeOf(obj)

	var tableName string

	// Проверяем наличие поля
	field, found := data.FieldByName("TableName")
	if found {
		// Получаем значение поля
		userValue := reflect.ValueOf(obj)
		fieldValue := userValue.FieldByName("TableName")

		if fieldValue.IsValid() {
			tableName = fieldValue.String()

			fmt.Printf("Поле '%s' найдено в структуре. Значение: %s\n", field.Name, tableName)
		} else {
			return fmt.Errorf("поле '%s' найдено, но его значение недоступно.", field.Name)
		}
	} else {
		return fmt.Errorf("Поле '%s' не найдено в структуре.\n", "TableName")
	}

	sqlQuery := "CREATE TABLE IF NOT EXISTS " + tableName + " ("

	for i := 0; i < data.NumField(); i++ {
		field := data.Field(i)
		if field.Name == "TableName" {
			continue
		}
		ormTag := field.Tag.Get("orm")
		if ormTag == "" {
			return fmt.Errorf("field %s does not have a tag", field.Name)
		} else if strings.Contains(ormTag, "ref") {
			// Ищем индекс подстроки "ref"
			start := strings.Index(ormTag, "ref")

			match := ormTag[start:]

			ormTag = strings.Replace(ormTag, " "+match, "", -1)

			sqlQuery += strings.ToLower(field.Name) + " " + ormTag + ", " + "FOREIGN KEY (" + strings.ToLower(field.Name) + ") REFERENCES " + strings.TrimPrefix(match, "ref ") + ", "
		} else {
			sqlQuery += strings.ToLower(field.Name) + " " + ormTag + ", "
		}
	}
	sqlQuery = strings.TrimSuffix(sqlQuery, ", ")
	sqlQuery += ");"

	fmt.Println(sqlQuery)

	_, err := conn.Exec(context.Background(), sqlQuery)
	if err != nil {
		return fmt.Errorf("error creating table:", err)
	} else {
		fmt.Println("Table created successfully or already exists.")
	}
	return nil
}

// Function for creating a new table object based on obj.TableName
func Create(obj interface{}) error {
	fmt.Println("CREATE", obj)
	values, columns := extractFields(obj)      // Getting the fields and their values
	columns = columns[2:]                      // Removing a column "TableName" and Id
	tableName, values := values[0], values[2:] // Getting the table name and delete it from the list of value

	columnsStr := "(" + strings.Join(columns, ", ") + ")" // Create a SQL-string with column names

	// Create a slice of placeholders
	placeholders := make([]string, len(values))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	// Create a string with placeholders for values
	placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
	// Create a SQL-string
	insertSQL := fmt.Sprintf(`INSERT INTO %s %s VALUES %s;`, tableName, columnsStr, placeholdersStr)
	// Send the request
	fmt.Println(insertSQL)
	_, err := conn.Exec(context.Background(), insertSQL, values...)
	if err != nil {
		return fmt.Errorf("row insertion error: %v", err)
	}
	fmt.Println("Row inserted successfully")
	return nil
}

// Function for updating information in the database
func Update(obj interface{}) error {
	values, columns := extractFields(obj) // get all the fields of the structure and their values
	fmt.Println(values, columns)
	columns = columns[2:]                      // Removing a column "TableName" and "ID"
	tableName, values := values[0], values[1:] // Getting the table name and delete it from the list of values

	strID := fmt.Sprint(values[0]) // Getting the obj ID
	values = values[1:]            // Removing the ID from the list of values

	// Create a string with column names and values for SQL-query
	updateData := ""
	for i := range columns {
		updateData += fmt.Sprintf(`%s = '%s', `, strings.ToLower(columns[i]), values[i])
	}

	// Remove the last comma and space
	updateData = strings.TrimSuffix(updateData, ", ")

	insertSQL := fmt.Sprintf(`UPDATE %s SET %s WHERE id = %s;`, tableName, updateData, strID)

	fmt.Println(insertSQL)
	// Passing strings
	_, err := conn.Exec(context.Background(), insertSQL)
	if err != nil {
		return fmt.Errorf("row update error: %v", err)
	}
	fmt.Println("Row updated successfully")
	return nil
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// convertObject function
func convertObject(obj interface{}, tableName string) (interface{}, error) {
	newType, ok := typeMap[tableName]
	if !ok {
		return nil, fmt.Errorf("type %s not found in typeMap", tableName)
	}

	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	if objType != newType {
		if !newType.AssignableTo(objType) {
			return nil, fmt.Errorf("type %s is not assignable to %s", newType, objType)
		}

		newObj := reflect.New(newType).Elem()
		newObj.Set(objValue.Convert(newType))
		return newObj.Interface(), nil
	}

	return obj, nil
}
