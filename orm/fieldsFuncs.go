package orm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
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

// Function for creating a new table object based on obj.TableName
func Create(obj interface{}) error {
	fmt.Println("CREATE", obj)
	// Getting the fields and their values
	values, columns := extractFields(obj)
	// Removing a column "TableName" and Id
	columns = columns[2:]
	// Getting the table name and delete it from the list of value
	tableName, values := values[0], values[2:]
	// Create a SQL-string with column names
	columnsStr := "(" + strings.Join(columns, ", ") + ")"

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
	// get all the fields of the structure and their values
	values, columns := extractFields(obj)
	fmt.Println(values, columns)
	// Removing a column "TableName" and "ID"
	columns = columns[2:]
	// Getting the table name and delete it from the list of values
	tableName, values := values[0], values[1:]
	// Getting the obj ID
	strID := fmt.Sprint(values[0])
	// Removing the ID from the list of values
	values = values[1:]

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

// converts an object to a type from typeMap
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
