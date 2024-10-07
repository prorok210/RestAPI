package orm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// Функция для создания нового объекта таблицы
func Create(obj interface{}) error {
	// Получаем поля и их значения из модели
	columns := extractFields(obj)
	// Удаляем столбец "TableName" и "ID"
	tableName := columns["TableName"].(string)
	delete(columns, "TableName")
	delete(columns, "ID")
	values := make([]interface{}, 0, len(columns))
	for _, value := range columns {
		values = append(values, value)
	}
	// Создаем шаблон SQL-запроса, который будем подставлять вместо плейсхолдеров
	columnsStr := "("
	for key := range columns {
		columnsStr += key + ", "
	}
	// Удаляем последнюю запятую и пробел и закрываем скобку
	columnsStr = strings.TrimSuffix(columnsStr, ", ") + ")"
	// Создаем слайс с плейсхолдерами
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	// Создаем строку с плейсхолдерами, на место которых будем подставлять значения
	placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
	// Создаем строку с SQL-запросом, подставляя название таблицы, столбцов и плейсхолдеры
	insertSQL := fmt.Sprintf(`INSERT INTO %s %s VALUES %s;`, tableName, columnsStr, placeholdersStr)
	// Выполняем запрос, подставляя значения в плейсхолдеры
	_, err := conn.Exec(context.Background(), insertSQL, values...)
	if err != nil {
		return fmt.Errorf("row insertion error: %v", err)
	}
	fmt.Println("Row inserted successfully")
	return nil
}

// Функция для обновления объекта в БД
func Update(obj interface{}) error {
	// Получаем все поля и их значения из модели. Ключ, название поля. По ключу доступно значение поля
	columns := extractFields(obj)
	tableName := columns["TableName"].(string)
	strID := fmt.Sprint(columns["ID"])
	// Удаляем столбец "TableName" и "ID"
	delete(columns, "TableName")
	delete(columns, "ID")

	// Создаем строку с обновляемыми данными в формате столбец = значение
	updateData := ""
	for column, value := range columns {
		updateData += fmt.Sprintf(`%s = '%s', `, strings.ToLower(column), value)
	}

	// Удаляем последнюю запятую и пробел
	updateData = strings.TrimSuffix(updateData, ", ")

	insertSQL := fmt.Sprintf(`UPDATE %s SET %s WHERE id = %s;`, tableName, updateData, strID)

	// Выполняем запрос
	_, err := conn.Exec(context.Background(), insertSQL)
	if err != nil {
		return fmt.Errorf("row update error: %v", err)
	}
	fmt.Println("Row updated successfully")
	return nil
}

func extractFields(obj interface{}) map[string]interface{} {
	// Используем reflect.Indirect для получения значения структуры
	val := reflect.Indirect(reflect.ValueOf(obj))

	// Проверяем, что obj — это структура или указатель на структуру
	if val.Kind() != reflect.Struct {
		return nil // или можно возвращать пустую карту или ошибку
	}
	typ := val.Type()

	var result = map[string]interface{}{}

	// Проходим по полям структуры
	for i := 0; i < val.NumField(); i++ {
		result[typ.Field(i).Name] = val.Field(i).Interface()
	}
	return result
}

// Переводим объект в тип из TypeTable, который соответствует переданной таблице (в таблице users хранятся объекты типа User и т.д.)
func convertObject(obj interface{}, tableName string) (interface{}, error) {
	newType, ok := TypeTable[tableName]
	if !ok {
		return nil, fmt.Errorf("type %s not found in typeMap", tableName)
	}

	// Содержимое переменной
	objValue := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	fmt.Println("TYPE", t, obj)
	// Если типы не равны и их нельзя привести друг к другу, возвращаем ошибку
	// _, ok = obj.(newType)
	// Указатель на новый объект типа newType
	newObj := reflect.New(newType).Elem()
	// Приводим значение objValue к типу newType и устанавливаем его в newObj. Проверка на совместимость была выше
	newObj.Set(objValue.Convert(newType))
	// Возвращаем тип интерфейса, чтобы можно было использовать в качестве объекта
	return newObj.Interface(), nil
}
