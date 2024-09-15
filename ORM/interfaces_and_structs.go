package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type BaseModel struct {
	TableName string
}

func (table *BaseModel) GetAll(conn *pgx.Conn) {
	selectSQL := fmt.Sprintf(`SELECT id, name, email FROM %s;`, table.TableName)
	fmt.Println(selectSQL)
	rows, err := conn.Query(context.Background(), selectSQL)
	if err != nil {
		log.Fatalf("Ошибка выполнения SELECT: %v", err)
	}
	defer rows.Close()

	fmt.Println("Данные из таблицы users:")

	// Итерируем по строкам и выводим данные
	for rows.Next() {
		var id int
		var name, email string

		// Считываем данные каждой строки
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			log.Fatalf("Ошибка сканирования строки: %v", err)
		}

		// Выводим данные строки
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", id, name, email)
	}

	// Проверка на ошибки после завершения итерации
	if rows.Err() != nil {
		log.Fatalf("Ошибка обработки строк: %v", rows.Err())
	}
}

// Таблицы
type TableUsers struct {
	BaseModel
}

type User struct {
	name  string
	email string
}

type TableDialogs struct {
	BaseModel
}
