package main

import (
	"context"
	"fmt"
	"reflect"
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
	// Сделать return нужных структур, не забыть про поле tableName
	return nil
}

// Фабричная функция для создания объектов на основании TableName
func (table *BaseTable) newModel(fields map[string]interface{}) (BaseCell, error) {
	if modelType, ok := tableRegistry[table.TableName]; ok {
		// Создание нового экземпляра нужного типа
		modelValue := reflect.New(modelType).Elem()

		// Устанавливаем значение поля TableName напрямую через рефлексию
		tableNameField := modelValue.FieldByName("TableName")
		if tableNameField.IsValid() && tableNameField.CanSet() {
			tableNameField.SetString(table.TableName)
		}

		model := modelValue.Addr().Interface().(BaseCell)

		if model == nil {
			return nil, fmt.Errorf("model is nil")
		}

		// Заполнение модели данными из fields
		for fieldName, value := range fields {
			if fieldName == "Id" {
				value = uint(value.(int32))
			}
			reflect.ValueOf(model).Elem().FieldByName(fieldName).Set(reflect.ValueOf(value))
		}

		return model, nil
	}
	return nil, fmt.Errorf("model not found in tableRegistry map")
}

func (table *BaseTable) getById(id uint) (interface{}, error) {
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
		columnName := capitalizeFirstLetter(string(fd.Name))
		result[columnName] = values[i]
	}
	fmt.Println(result)
	obj, err := table.newModel(result)
	if err != nil {
		return nil, fmt.Errorf("error creating model: %v", err)
	}
	fmt.Println(obj)

	// Проверка на ошибки после завершения итерации
	if rows.Err() != nil {
		return nil, fmt.Errorf("string processing error: %v", rows.Err())
	}

	// Возвращаем пользователя и nil как ошибку, если всё прошло успешно
	// Дополнить, чтобы вместо obj возвращало obj.(*type) для возврата не интерфейса, а структуры на основе TableName
	convertObject(obj, table.TableName)
	return obj, nil
}
