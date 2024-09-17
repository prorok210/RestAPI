package main

// Интерфейс для строк из таблиц
type Cell interface {
	ToFields() ([]interface{}, []string) // Метод для получения значений и колонок для вставки
}

// Базовая модель таблицы
type BaseModel struct {
	TableName string
}

// Таблица пользователей
type TableUsers struct {
	BaseModel
}

// Объект из таблицы пользователей
type User struct {
	Name  string
	Email string
}

// Таблица диалогов
type TableDialogs struct {
	BaseModel
}

// Объект из таблицы диалогов
type Dialog struct {
	name string
}
