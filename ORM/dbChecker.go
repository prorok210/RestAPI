package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// Function for checking the similarity of tables in database and structures
func CheckTables() error {
	for tableName, modelType := range tableRegistry {
		fmt.Printf("Checking table %s...\n", tableName)

		// Getting the current table structure from the database
		dbColumns, err := getTableColumns(tableName)

		fmt.Println("dbColumns", dbColumns)
		if err != nil {
			return fmt.Errorf("error getting table columns %s: %v", tableName, err)
		}

		// Comparing the structure with the model
		modelColumns, err := getModelColumns(modelType)
		if err != nil {
			return fmt.Errorf("error getting model fields %s: %v", modelType.Name(), err)
		}
		fmt.Println("modelColumns", modelColumns)

		// Column comparison
		if !compareColumns(dbColumns, modelColumns) {
			fmt.Printf("Difference in table structure %s!\n", tableName)
		} else {
			fmt.Printf("Table %s matches the model.\n", tableName)
		}
	}

	return nil
}

// Getting a list of table columns from a database
func getTableColumns(tableName string) (map[string]string, error) {
	query := fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_name='%s';", tableName)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer rows.Close()

	columns := make(map[string]string)
	for rows.Next() {
		var columnName string
		var dataType string
		err := rows.Scan(&columnName, &dataType)
		if err != nil {
			return nil, fmt.Errorf("line reading error: %v", err)
		}
		columns[columnName] = dataType
	}

	return columns, nil
}

// Getting a list of fields from a structure
func getModelColumns(modelType reflect.Type) (map[string]string, error) {
	columns := make(map[string]string)

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		ormTag := field.Tag.Get("orm")
		if ormTag != "" {
			// Extract type from ORM tag
			ormTag = strings.Split(ormTag, " ")[0]

			columns[strings.ToLower(field.Name)] = tagToSqlType[ormTag]
		}
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("no fields found in model")
	}

	return columns, nil
}

// Comparing two map of columns
func compareColumns(dbColumns, modelColumns map[string]string) bool {
	// Checking the presence of all columns from the model in the database
	for col, modelType := range modelColumns {
		if dbType, exists := dbColumns[col]; !exists {
			fmt.Printf("Column %s is missing from the database.\n", col)
			return false
		} else if !strings.Contains(dbType, modelType) {
			fmt.Printf("Column type mismatch %s: in database %s, in model %s.\n", col, dbType, modelType)
			return false
		}
	}

	// Check for extra columns in the database
	for col := range dbColumns {
		if _, exists := modelColumns[col]; !exists {
			fmt.Printf("An extra column %s was found in the database.\n", col)
			return false
		}
	}

	return true
}
