package orm

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

// Function to get all values ​​from the database
func (table *BaseTable) GetAll() error {
	selectSQL := fmt.Sprintf(`SELECT * FROM %s;`, table.TableName)
	rows, err := conn.Query(context.Background(), selectSQL)
	if err != nil {
		return fmt.Errorf("select error: %v", err)
	}
	defer rows.Close()

	// We get a description of the fields (columns)
	fieldDescriptions := rows.FieldDescriptions()

	// Create a slice to save the values ​​of row
	values := make([]interface{}, len(fieldDescriptions))
	valuePtrs := make([]interface{}, len(fieldDescriptions))

	for rows.Next() {
		// Fill valuePtrs with pointers to values
		for i := range fieldDescriptions {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return fmt.Errorf("line scan error: %v", err)
		}

		for i, fd := range fieldDescriptions {

			fmt.Printf("%s: %v\n", string(fd.Name), values[i])
		}
	}

	if rows.Err() != nil {
		return fmt.Errorf("line processing error: %v", rows.Err())
	}
	return nil
}

// Fabric function to create objects based on TableName
func (table *BaseTable) newModel(fields map[string]interface{}) (interface{}, error) {
	if modelType, ok := TypeTable[table.TableName]; ok {
		// Creating a new instance of the desired type
		model := reflect.New(modelType).Elem()

		// Setting the value of the TableName field using reflection
		tableNameField := model.FieldByName("TableName")
		if tableNameField.IsValid() && tableNameField.CanSet() {
			tableNameField.SetString(table.TableName)
		} else {
			return nil, fmt.Errorf("field TableName is invalid or cannot be set")
		}

		// Filling the model with data from fields
		for fieldName, value := range fields {
			fieldMeta, err := FindFieldByName(modelType, fieldName)
			fmt.Println("field", fieldName, "fieldMeta: ", fieldMeta)
			if err != nil {
				return nil, fmt.Errorf("field %s not found in model %s: %v", fieldName, modelType.Name(), err)
			}

			field := model.FieldByName(fieldMeta.Name)

			// Проверяем, что поле существует и его можно установить
			if !field.IsValid() || !field.CanSet() {
				return nil, fmt.Errorf("field %s is invalid or cannot be set", fieldName)
			}

			// Проверяем тип поля и значение, чтобы убедиться, что они совместимы
			fieldType := reflect.TypeOf(field.Interface())
			valueType := reflect.TypeOf(value)
			if t, ok := value.(time.Time); ok {
				// Если да, форматируем его по шаблону
				value = t.Format("2006-01-02 15:04:05")
				if strValue, ok := value.(string); ok {
					// Успешно привели значение к строке
					field.SetString(strValue)
				} else {
					return nil, fmt.Errorf("error converting time.Time to string")
				}
			} else {
				valueValue := reflect.ValueOf(value)

				// Проверка на возможность установки
				if !valueValue.IsValid() {
					return nil, fmt.Errorf("invalid value")
				}

				// Приведение типов
				if valueType.ConvertibleTo(fieldType) {
					// Преобразуем значение к типу поля и устанавливаем его
					field.Set(valueValue.Convert(fieldType))
				} else {
					// Если типы не совместимы
					return nil, fmt.Errorf("field %s type %s does not match value type %s", fieldName, fieldType, valueType)
				}
				// 	field.Set(reflect.ValueOf(value))
			}

			// Устанавливаем значение в поле
			fmt.Println("Set field: ", field)
		}

		return model, nil
	}
	return nil, fmt.Errorf("model not found in tableRegistry map")
}

// Function to get a data by id
func (table *BaseTable) GetById(id int) (interface{}, error) {
	selectSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1;`, table.TableName) // Используем placeholder для безопасности
	rows, err := conn.Query(context.Background(), selectSQL, id)
	if err != nil {
		return nil, fmt.Errorf("select error: %v", err)
	}
	defer rows.Close()

	// Проверяем, есть ли строки в результате запроса
	if !rows.Next() {
		return nil, fmt.Errorf("row with id = %d not found", id)
	}

	// Получаем метаданные столбцов (названия столбцов)
	fieldDescriptions := rows.FieldDescriptions()

	// Получаем значения столбцов
	values, err := rows.Values()
	if err != nil {
		return nil, fmt.Errorf("error getting row values: %v", err)
	}

	// Создаем map для хранения данных
	result := make(map[string]interface{})

	// Итерируем по столбцам и значениям, записывая их в map
	for i, fd := range fieldDescriptions {
		columnName := CapitalizeFirstLetter(string(fd.Name))
		result[columnName] = values[i]
	}
	fmt.Println("RESULT", result)
	obj, err := table.newModel(result)
	if err != nil {
		return nil, fmt.Errorf("error creating model: %v", err)
	}

	// Проверка на ошибки после завершения итерации
	if rows.Err() != nil {
		return nil, fmt.Errorf("string processing error: %v", rows.Err())
	}

	// Возвращаем пользователя и nil как ошибку, если всё прошло успешно
	// Дополнить, чтобы вместо obj возвращало obj.(*type) для возврата не интерфейса, а структуры на основе TableName
	convertObject(obj, table.TableName)
	return obj, nil
}
