package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// Функция для проверки соответствия таблиц в БД и структур в коде
func CheckTables() error {
	// Проходим по каждой структуре модели
	for tableName, modelType := range tableRegistry {
		fmt.Printf("Проверка таблицы %s...\n", tableName)

		// Получаем текущую структуру таблицы из базы данных
		dbColumns, err := getTableColumns(tableName)

		fmt.Println("dbColumns", dbColumns)
		if err != nil {
			return fmt.Errorf("ошибка получения столбцов таблицы %s: %v", tableName, err)
		}

		// Сравниваем структуру с моделью
		modelColumns, err := getModelColumns(modelType)
		if err != nil {
			return fmt.Errorf("ошибка получения столбцов модели %s: %v", modelType.Name(), err)
		}
		fmt.Println("modelColumns", modelColumns)

		// Сравнение столбцов
		if !compareColumns(dbColumns, modelColumns) {
			fmt.Printf("Различие в структуре таблицы %s!\n", tableName)
		} else {
			fmt.Printf("Таблица %s соответствует модели.\n", tableName)
		}
	}

	return nil
}

// Получение списка столбцов таблицы из базы данных
func getTableColumns(tableName string) (map[string]string, error) {
	query := fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_name='%s';", tableName)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	columns := make(map[string]string)
	for rows.Next() {
		var columnName string
		var dataType string
		err := rows.Scan(&columnName, &dataType)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}
		columns[columnName] = dataType
	}

	return columns, nil
}

// Получение списка полей из модели (структуры)
func getModelColumns(modelType reflect.Type) (map[string]string, error) {
	columns := make(map[string]string)

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		ormTag := field.Tag.Get("orm")
		if ormTag != "" {
			// Извлечение типа из ORM тега
			ormTag = strings.Split(ormTag, " ")[0]

			columns[strings.ToLower(field.Name)] = tagToSqlType[ormTag]
		}
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("не найдено полей в модели")
	}

	return columns, nil
}

// Сравнение двух списков столбцов
func compareColumns(dbColumns, modelColumns map[string]string) bool {
	// Проверка наличия всех столбцов из модели в БД
	for col, modelType := range modelColumns {
		if dbType, exists := dbColumns[col]; !exists {
			fmt.Printf("Столбец %s отсутствует в БД.\n", col)
			return false
		} else if !strings.Contains(dbType, modelType) {
			fmt.Printf("Несоответствие типа столбца %s: в БД %s, в модели %s.\n", col, dbType, modelType)
			return false
		}
	}

	// Проверка на лишние столбцы в БД
	for col := range dbColumns {
		if _, exists := modelColumns[col]; !exists {
			fmt.Printf("Лишний столбец %s найден в БД.\n", col)
			return false
		}
	}

	return true
}
