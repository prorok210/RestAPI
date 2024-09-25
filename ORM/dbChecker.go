package main

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// Function for checking the similarity of tables in database and structures
func CheckTables() error {
	for tableName, modelType := range tableRegistry {
		fmt.Printf("Checking table %s...\n", tableName)
		fmt.Println("modelType", modelType)

		// Getting the current table structure from the database
		dbColumns, err := getTableColumns(tableName)

		fmt.Println("dbColumns", dbColumns)
		if err != nil {
			return fmt.Errorf("error getting table columns %s: %v", tableName, err)
		}

		// Comparing the structure with the model
		// modelColumns, err := getModelColumns(modelType)
		// if err != nil {
		// 	return fmt.Errorf("error getting model fields %s: %v", modelType.Name(), err)
		// }
		// fmt.Println("modelColumns", modelColumns)

		// // Column comparison
		// if !compareColumns(dbColumns, modelColumns) {
		// 	fmt.Printf("Difference in table structure %s!\n", tableName)
		// } else {
		// 	fmt.Printf("Table %s matches the model.\n", tableName)
		// }
	}

	return nil
}

// Getting a list of table columns from a database
func getTableColumns(tableName string) (map[string]map[string]string, error) {
	query := fmt.Sprintf(`
		SELECT
			c.column_name,
			c.data_type,
			c.is_nullable,
			c.column_default,
			c.character_maximum_length,
			c.numeric_scale,
			c.is_identity,
			tc.constraint_type,
			kcu.constraint_name,
			rc.update_rule,
			rc.delete_rule
		FROM
			information_schema.columns AS c
		LEFT JOIN 
			information_schema.key_column_usage AS kcu 
			ON c.table_name = kcu.table_name AND c.column_name = kcu.column_name
		LEFT JOIN 
			information_schema.table_constraints AS tc 
			ON kcu.constraint_name = tc.constraint_name
		LEFT JOIN 
			information_schema.referential_constraints AS rc 
			ON tc.constraint_name = rc.constraint_name
		LEFT JOIN 
			information_schema.check_constraints AS cc 
			ON cc.constraint_name = tc.constraint_name
		WHERE 
			c.table_name = '%s';`, tableName)

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer rows.Close()

	columns := make(map[string]map[string]string)
	for rows.Next() {
		var columnName, dataType, isNullable, columnDefault, constraintType, constraintName, updateRule, deleteRule sql.NullString
		var charMaxLength, numericScale sql.NullInt64
		var isIdentityStr sql.NullString

		err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault, &charMaxLength, &numericScale, &isIdentityStr, &constraintType, &constraintName, &updateRule, &deleteRule)
		if err != nil {
			return nil, fmt.Errorf("line reading error: %v", err)
		}

		// Преобразуем isIdentity в строку "true" или "false"
		isIdentity := false
		if isIdentityStr.Valid && isIdentityStr.String == "YES" {
			isIdentity = true
		}

		columnDetails := map[string]string{
			"column_name":          columnName.String,
			"data_type":            dataType.String,
			"is_nullable":          isNullable.String,
			"column_default":       columnDefault.String,
			"character_max_length": fmt.Sprintf("%v", charMaxLength.Int64),
			"numeric_scale":        fmt.Sprintf("%v", numericScale.Int64),
			"is_identity":          fmt.Sprintf("%v", isIdentity),
			"constraint_type":      constraintType.String,
			"constraint_name":      constraintName.String,
			"update_rule":          updateRule.String,
			"delete_rule":          deleteRule.String,
		}

		columns[columnName.String] = columnDetails
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
