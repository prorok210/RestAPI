package orm

import (
	"RestAPI/core"
	"RestAPI/db"
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

	// checking tables for compliance with structures
	err := CheckTables()
	if err != nil {
		return fmt.Errorf("error checking tables: %v", err)
	}

	log.Println("Successfully connected to the database.")

	return nil
}

func CreateTable(obj interface{}) error {
	data := reflect.TypeOf(obj)

	var tableName string

	// Checking the presence of the field
	field, found := data.FieldByName("TableName")
	if found {
		// Getting the field value
		userValue := reflect.ValueOf(obj)
		fieldValue := userValue.FieldByName("TableName")

		if fieldValue.IsValid() {
			tableName = fieldValue.String()

			fmt.Printf("Field '%s' was found in the structure. Value: %s\n", field.Name, tableName)
		} else {
			return fmt.Errorf("field '%s' was found, but its value is not available", field.Name)
		}
	} else {
		return fmt.Errorf("field 'TableName' was not found in the structure")
	}

	// Checking a table exists or not
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name=$1);"
	err := conn.QueryRow(context.Background(), query, tableName).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking table existence: %v", err)
	}

	if exists {
		fmt.Println("Table already exists.")
		return nil
	}

	sqlQuery := "CREATE TABLE IF NOT EXISTS " + tableName + " ("

	for i := 0; i < data.NumField(); i++ {
		field := data.Field(i)
		if field.Name == "TableName" {
			continue
		}
		ormTag := field.Tag.Get("orm")
		ormTag = strings.Replace(ormTag, "_", " ", -1)
		if ormTag == "" {
			return fmt.Errorf("field %s does not have a tag", field.Name)
		} else if strings.Contains(ormTag, "ref") {
			// Looking for the index of the substring "ref"
			start := strings.Index(ormTag, "ref")

			// Getting the substring from the index to the end of the string
			match := ormTag[start:]
			// Removing the substring from the tag
			ormTag = strings.Replace(ormTag, " "+match, "", -1)
			// Add the field name and value from the orm tag with reference
			sqlQuery += strings.ToLower(field.Name) + " " + ormTag + ", " + "FOREIGN KEY (" + strings.ToLower(field.Name) + ") REFERENCES " + strings.TrimPrefix(match, "ref ") + ", "
		} else {
			// Add the field name and value from the orm tag with reference
			sqlQuery += strings.ToLower(field.Name) + " " + ormTag + ", "
		}
	}
	// Remove the last comma and space
	sqlQuery = strings.TrimSuffix(sqlQuery, ", ")
	sqlQuery += ");"

	fmt.Println(sqlQuery)

	_, err = conn.Exec(context.Background(), sqlQuery)
	if err != nil {
		return fmt.Errorf("error creating table %s", err)
	} else {
		fmt.Println("Table created successfully")
	}
	return nil
}

// Function for checking the similarity of tables in database and structures
func CheckTables() error {
	for tableName, modelType := range db.TableRegistry {
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
		if ok {
			fmt.Printf("Table %s matches the model\n", tableName)
		} else {
			fmt.Printf("Difference in table structure %s!\n", tableName)
		}
	}

	return nil
}

// Getting a list of table columns from a database
func getTableColumns(tableName string) (map[string]map[string]string, error) {
	// use service tables to get column data
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

		// Get orm tags
		ormTag := strings.ToLower(field.Tag.Get("orm"))
		fmt.Println("ormTag", ormTag)
		if ormTag != "" {
			// every attribute is separated by a space
			ormTagSlice := strings.Split(ormTag, " ")
			// first attribute is the type. Translating attributes into sql text
			value, ok := tagToSqlType[ormTagSlice[0]]
			if !ok {
				return nil, fmt.Errorf("unknown type %s. You can add it in tagToSqlType in constants.go", ormTagSlice[0])
			} else {
				column["data_type"] = value
			}
			// Filling the column with data
			if Contains(ormTagSlice, "serial_primary_key") {
				column["constraint_type"] = "PRIMARY KEY"
				column["is_identity"] = "true"
				column["is_nullable"] = "NO"
			}
			if Contains(ormTagSlice, "ref") {
				column["constraint_type"] = "FOREIGN KEY"
				column["is_identity"] = "true"
				if Contains(ormTagSlice, "on_update") {
					updateRule := find_rule(ormTag, "on_update")
					column["update_rule"] = updateRule
				}
				if Contains(ormTagSlice, "on_delete") {
					fmt.Println("on_delete contains")
					deleteRule := find_rule(ormTag, "on_delete")
					column["delete_rule"] = deleteRule
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
				fmt.Println("FIELDS", key, value, valueModel[key])
				return false, fmt.Errorf("column %s does not match the model, db: %s and model: %s", key, value, valueModel[key])
			}
		}
	}

	return true, nil
}

func find_rule(tag string, rule string) string {
	// find rule in tag
	start := strings.Index(tag, rule)
	end := start + strings.Index(tag[start:], " ")
	// if end not found
	if end == start-1 {
		end = len(tag)
	}
	result := strings.ToUpper(strings.Replace(tag[start:end], "_", " ", -1))
	// for SET NULL case
	if len(strings.Split(result, " ")) < 3 {
		return ""
	}
	return strings.Join(strings.Split(result, " ")[2:], "")
}
