package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (user *User) ToFields() ([]interface{}, []string) {
	return extractFields(user)
}

// Generic function to extract fields
func extractFields(obj interface{}) ([]interface{}, []string) {
	fmt.Println("Extracting fields from an object:", obj)
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
	fmt.Println("Data extracted successfully:")
	return values, columns
}

func InitDB() error {
	var InitDBError error
	conn, InitDBError = pgx.Connect(context.Background(), CONNECTIONDATA)
	if InitDBError != nil {
		log.Printf("Database connect error: %v", InitDBError)
		return fmt.Errorf("Database connect error: %v", InitDBError)
	}

	fmt.Println("Successfully connected to the database.")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50),
		email VARCHAR(50)
	);`

	// Execute a table creation request
	_, err := conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return fmt.Errorf("Error creating table: %v", err)
	} else {
		fmt.Println("Table created successfully or already exists.")
	}

	return nil
}

// Function for creating a new table object based on obj.TableName
func Create(obj interface{}) error {
	fmt.Println("CREATE", obj)
	values, columns := extractFields(obj)      // Getting the fields and their values
	columns = columns[1:]                      // Removing a column "TableName"
	tableName, values := values[0], values[1:] // Getting the table name and delete it from the list of values

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
	_, err := conn.Exec(context.Background(), insertSQL, values...)
	if err != nil {
		log.Printf("Row insertion error: %v", err)
		return fmt.Errorf("Row insertion error: %v", err)
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
		log.Printf("Row update error: %v", err)
	}
	fmt.Println("Row updated successfully")
	return nil
}

// Функция, которая принимает интерфейс
func ProcessUser(i interface{}) {
	// Приведение интерфейса к структуре User
	user, ok := i.(User)
	if ok {
		// Приведение успешно
		fmt.Printf("User: %s, Email: %s\n", user.Name, user.Email)
	} else {
		// Приведение не удалось
		fmt.Println("Приведение к User не удалось")
	}
}
