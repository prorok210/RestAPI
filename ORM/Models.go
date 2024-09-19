package main

type BaseCell interface {
	ToFields() ([]interface{}, []string)
}

// Base Table structure
type BaseTable struct {
	TableName string
}

// Table of users
type TableUsers struct {
	BaseTable
}

// User structure
type User struct {
	TableName string
	ID        uint
	Name      string
	Email     string
}

// Table of dialogs
type TableDialogs struct {
	BaseTable
}

// Dialog structure
type Dialog struct {
	name string
}
