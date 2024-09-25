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

type TableUsrees struct {
	BaseTable
}

type TableMessages struct {
	BaseTable
}

type Message struct {
	TableName string
	Id        int    `orm:"serial primary key"`
	UserId    int    `orm:"int ref users(id) on update cascade on delete cascade"`
	Text      string `orm:"varchar not null"`
}

// User structure
type User struct {
	TableName string
	Id        int    `orm:"serial primary key"`
	Name      string `orm:"varchar not null"`
	Email     string `orm:"varchar not null"`
}

// Table of dialogs
type TableDialogs struct {
	BaseTable
}

// Dialog structure
type Dialog struct {
	TableName string
	Id        int    `orm:"serial primary key"`
	Owner     string `orm:"varchar not null"`
	Opponent  string `orm:"varchar not null"`
}
