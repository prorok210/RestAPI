package main

// Объект из таблицы пользователей
type User struct {
	TableName string
	ID        uint
	Name      string
	Email     string
}

// Базовая модель таблицы
type BaseTable struct {
	TableName string
}

// Таблица пользователей
type TableUsers struct {
	BaseTable
}

// Таблица диалогов
type TableDialogs struct {
	BaseTable
}

// Объект из таблицы диалогов
type Dialog struct {
	name string
}
