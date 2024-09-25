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

		if err != nil {
			return fmt.Errorf("error getting table columns %s: %v", tableName, err)
		}
		// fmt.Println("dbColumns", dbColumns)

		// Comparing the structure with the model
		modelColumns, err := getModelColumns(modelType)
		if err != nil {
			return fmt.Errorf("error getting model fields %s: %v", modelType.Name(), err)
		}
		// fmt.Println("modelColumns", modelColumns)

		// Column comparison
		ok, err := compareColumns(dbColumns, modelColumns)
		if err != nil {
			return fmt.Errorf("error comparing columns %s: %v", tableName, err)
		}
		if !ok {
			fmt.Printf("Difference in table structure %s!\n", tableName)
		} else {
			fmt.Printf("Table %s matches the model\n", tableName)
		}
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

		// fmt.Println("UD", updateRule.String, deleteRule.String)

		columnDetails := map[string]string{
			"column_name":     columnName.String,
			"data_type":       dataType.String,
			"is_nullable":     isNullable.String,
			"constraint_type": constraintType.String,
			"update_rule":     updateRule.String,
			"delete_rule":     deleteRule.String,
		}

		columns[columnName.String] = columnDetails
		// fmt.Println("CD", columnDetails, columnName.String)
	}

	return columns, nil
}

// Getting a list of fields from a structure
func getModelColumns(modelType reflect.Type) (map[string]map[string]string, error) {
	columns := map[string]map[string]string{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
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
		fmt.Println("ormTag", ormTag)
		if ormTag != "" {
			ormTagSlice := strings.Split(ormTag, " ")
			value, ok := tagToSqlType[ormTagSlice[0]]
			if !ok {
				return nil, fmt.Errorf("unknown type %s. You can add it in tagToSqlType in constants.go", ormTagSlice[0])
			} else {
				column["data_type"] = value
			}
			if Contains(ormTagSlice, "serial_primary_key") {
				column["constraint_type"] = "PRIMARY KEY"
				column["is_identity"] = "true"
				column["is_nullable"] = "NO"
			}
			if Contains(ormTagSlice, "ref") {
				column["constraint_type"] = "FOREIGN KEY"
				column["is_identity"] = "true"
				if Contains(ormTagSlice, "on_update") {
					fmt.Println("on_update contains")
					updateRule := find_rule(ormTag, "on_update")
					column["update_rule"] = updateRule
				} else {
					fmt.Println("on_update NOOOO contains")
				}
				if Contains(ormTagSlice, "on_delete") {
					fmt.Println("on_delete contains")
					deleteRule := find_rule(ormTag, "on_delete")
					column["update_rule"] = deleteRule
				} else {
					fmt.Println("on_delete NOOOOO scontains")
				}
			}
			if Contains(ormTagSlice, "not_null") {
				column["is_nullable"] = "NO"
			}
			if Contains(ormTagSlice, "unique") {
				column["constraint_type"] = "UNIQUE"
			}

		}
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("no fields found in model")
	}

	return columns, nil
}

// Comparing two map of columns
func compareColumns(dbColumns, modelColumns map[string]map[string]string) (bool, error) {
	// Checking the presence of all columns from the model in the database
	for key, valueDb := range dbColumns {
		if _, ok := modelColumns[key]; !ok {
			return false, fmt.Errorf("column %s not found in model", key)
		}

		// Checking the presence of all columns from the database in the model
		if _, ok := dbColumns[key]; !ok {
			return false, fmt.Errorf("column %s not found in database", key)
		}

		valueModel := modelColumns[key]
		// Checking the similarity of the columns
		for key, value := range valueDb {
			if value != valueModel[key] {
				return false, fmt.Errorf("column %s does not match the model, db: %s and model: %s", key, value, valueModel[key])
			}
		}
	}

	return true, nil
}
