package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
	connStr := GetDBData()

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Успешное подключение к PostgreSQL через pgx!")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50),
		email VARCHAR(50)
	);`

	// Выполняем запрос создания таблицы
	_, err = conn.Exec(context.Background(), createTableSQL)
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

	// Таблицу в объект языка
	users := &TableUsers{BaseModel{TableName: "users"}}
	// когда появится метод Create, тогда не нужно будет передавать conn
	// Либо сделать conn отдельным полем BaseModel
	users.GetAll(conn)
}
