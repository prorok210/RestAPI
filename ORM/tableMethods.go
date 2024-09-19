package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"unicode"
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
func (bt *BaseTable) newModel() BaseCell {
	if modelType, ok := tableRegistry[bt.TableName]; ok {
		// Создание нового экземпляра нужного типа
		modelValue := reflect.New(modelType).Elem().Addr().Interface().(BaseCell)
		return modelValue
	}
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

// Пример функции get
func (bt *BaseTable) Get(fields map[string]interface{}) BaseCell {
	// Создание нового объекта на основе TableName
	model := bt.newModel()

	fmt.Println(model)

	if model == nil {
		fmt.Println("Неизвестная таблица:", bt.TableName)
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

func (table *BaseTable) getById(id uint) (BaseCell, error) {
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
	obj := table.Get(result)
	fmt.Println(obj)
	if err != nil {
		log.Printf("Line scan error: %v\n", err)
		return nil, fmt.Errorf("Line scan error: %v\n", err)
	}

	// Проверка на ошибки после завершения итерации
	if rows.Err() != nil {
		log.Printf("string processing error: %v\n", rows.Err())
		return nil, fmt.Errorf("string processing error: %v\n", rows.Err())
	}

	// Возвращаем пользователя и nil как ошибку, если всё прошло успешно
	return obj, nil
}
