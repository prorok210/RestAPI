package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
)

// Конструктор для структуры User
func newUser(name string, email string) *User {
	fmt.Println("Создан новый пользователь:", name, email)
	return &User{
		TableName: "users",
		Name:      name,
		Email:     email,
	}
}

// Общая функция для извлечения полей
func extractFields(obj interface{}) ([]interface{}, []string) {
	fmt.Println("Извлечение полей из объекта:", obj)
	val := reflect.ValueOf(obj).Elem()
	typ := reflect.TypeOf(obj).Elem()

	var values []interface{}
	var columns []string

	// Проходим по полям структуры
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		// Пропускаем поле с именем "ID"
		if fieldName == "ID" {
			continue
		}

		// Добавляем имя поля в список колонок
		columns = append(columns, fieldName)

		// Добавляем значение поля в список значений
		values = append(values, val.Field(i).Interface())
	}
	fmt.Println("данные о пользователе успешно отправлены")
	return values, columns
}

// Инициализация БД
func InitDB() {

	conn, InitDBError = pgx.Connect(context.Background(), CONNECTIONDATA)
	fmt.Println(conn)
	if InitDBError != nil {
		log.Fatalf("Ошибка подключения к БД: %v", InitDBError)
	}

	fmt.Println("Успешное подключение к PostgreSQL через pgx!")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50),
		email VARCHAR(50)
	);`

	// Выполняем запрос создания таблицы
	_, err := conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	} else {
		fmt.Println("Таблица users успешно создана или уже существует.")
	}

	// SQL-запрос для вставки строки в таблицу
	insertSQL := `
	INSERT INTO users (name, email) VALUES ($1, $2);
	`

	// Вставляем строки
	_, err = conn.Exec(context.Background(), insertSQL, "Иван Иванов", "ivan@example.com")
	_, err = conn.Exec(context.Background(), insertSQL, "Петр Петров", "petr@example.com")
	_, err = conn.Exec(context.Background(), insertSQL, "Сидор Сидоров", "sidor@example.com")
}

func (table *BaseTable) GetAll() {
	selectSQL := fmt.Sprintf(`SELECT * FROM %s;`, table.TableName)
	rows, err := conn.Query(context.Background(), selectSQL)
	if err != nil {
		log.Fatalf("Ошибка выполнения SELECT: %v", err)
	}
	defer rows.Close()

	// Получаем описание полей (столбцов)
	fieldDescriptions := rows.FieldDescriptions()

	// Создаем срез для хранения значений каждой строки
	values := make([]interface{}, len(fieldDescriptions))
	valuePtrs := make([]interface{}, len(fieldDescriptions))

	// Итерируем по строкам
	for rows.Next() {
		// Заполняем valuePtrs указателями на значения
		for i := range fieldDescriptions {
			valuePtrs[i] = &values[i]
		}

		// Сканируем строку
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatalf("Ошибка сканирования строки: %v", err)
		}

		// Выводим данные строки
		fmt.Println("Данные строки:")
		for i, fd := range fieldDescriptions {

			fmt.Printf("%s: %v\n", string(fd.Name), values[i])
		}
	}

	// Проверка на ошибки после завершения итерации
	if rows.Err() != nil {
		log.Fatalf("Ошибка обработки строк: %v", rows.Err())
	}
	// Сделать return нужных структур, не забыть про поле tableName
}

func Create(obj interface{}) {
	fmt.Println("CREATE", obj)
	values, columns := extractFields(obj)                 // Получаем поля и их значения
	fmt.Println("Значения:", values, "Колонки:", columns) // Выводим их на экран
	columns = columns[1:]                                 // Получаем имя таблицы
	tableName, values := values[0], values[1:]            // Получаем значения
	fmt.Println("Имя таблицы:", tableName)                // Выводим его на экран

	columnsStr := "(" + strings.Join(columns, ", ") + ")"

	// Создаем плейсхолдеры для значений ($1, $2, ...)
	placeholders := make([]string, len(values))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
	insertSQL := fmt.Sprintf(`INSERT INTO %s %s VALUES %s;`, tableName, columnsStr, placeholdersStr)
	// Вставляем строки
	_, err := conn.Exec(context.Background(), insertSQL, values...)
	if err != nil {
		log.Fatalf("Ошибка вставки строки: %v", err)
	}
	fmt.Println("Строка успешно добавлена. Запрос:", insertSQL)
}
