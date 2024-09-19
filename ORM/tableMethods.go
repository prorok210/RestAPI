package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
)

// Function to get all values ​​from the database
func (table *BaseTable) GetAll() error {
	selectSQL := fmt.Sprintf(`SELECT * FROM %s;`, table.TableName)
	rows, err := conn.Query(context.Background(), selectSQL)
	if err != nil {
		log.Printf("SELECT error: %v\n", err)
		return fmt.Errorf("SELECT error: %v\n", err)
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
			log.Printf("Line scan error: %v\n", err)
			return fmt.Errorf("Line scan error: %v\n", err)
		}

		for i, fd := range fieldDescriptions {

			fmt.Printf("%s: %v\n", string(fd.Name), values[i])
		}
	}

	if rows.Err() != nil {
		log.Printf("Ошибка обработки строк: %v\n", rows.Err())
		return fmt.Errorf("Ошибка обработки строк: %v\n", rows.Err())
	}
	// Сделать return нужных структур, не забыть про поле tableName
	return nil
}

// Фабричная функция для создания объектов на основании TableName
func (table *BaseTable) newModel(fields map[string]interface{}) BaseCell {
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
			fmt.Println("Неизвестная таблица:", table.TableName)
			return nil
		}

		// Заполнение модели данными из fields
		for fieldName, value := range fields {
			if fieldName == "Id" {
				value = uint(value.(int32))
			}
			reflect.ValueOf(model).Elem().FieldByName(fieldName).Set(reflect.ValueOf(value))
		}

		return model
	}
	return nil
}

func (table *BaseTable) getById(id uint) (interface{}, error) {
	selectSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1;`, table.TableName) // Используем placeholder для безопасности
	rows, err := conn.Query(context.Background(), selectSQL, id)
	if err != nil {
		log.Printf("SELECT error: %v\n", err)
		return nil, fmt.Errorf("SELECT error: %v\n", err)
	}
	defer rows.Close()

	// Проверяем, есть ли строки в результате запроса
	if !rows.Next() {
		log.Printf("Row with id = %d not found\n", id)
		return nil, fmt.Errorf("Row with id = %d not found\n", id)
	}

	// Получаем метаданные столбцов (названия столбцов)
	fieldDescriptions := rows.FieldDescriptions()

	// Получаем значения столбцов
	values, err := rows.Values()
	if err != nil {
		log.Printf("Error getting row values: %v\n", err)
		return nil, fmt.Errorf("Error getting row values: %v\n", err)
	}

	// Создаем map для хранения данных
	result := make(map[string]interface{})

	// Итерируем по столбцам и значениям, записывая их в map
	for i, fd := range fieldDescriptions {
		columnName := capitalizeFirstLetter(string(fd.Name))
		result[columnName] = values[i]
	}
	fmt.Println(result)
	obj := table.newModel(result)
	fmt.Println(obj)

	// Проверка на ошибки после завершения итерации
	if rows.Err() != nil {
		log.Printf("string processing error: %v\n", rows.Err())
		return nil, fmt.Errorf("string processing error: %v\n", rows.Err())
	}

	// Возвращаем пользователя и nil как ошибку, если всё прошло успешно
	// Дополнить, чтобы вместо obj возвращало obj.(*type) для возврата не интерфейса, а структуры на основе TableName
	convertObject(obj, table.TableName)
	return obj, nil
}
