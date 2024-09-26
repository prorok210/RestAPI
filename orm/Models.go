package orm

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
	Id        int    `orm:"serial_primary_key"`
	UserId    int    `orm:"int ref_users(id) on_update_cascade on_delete_cascade"`
	Text      string `orm:"varchar not_null"`
}

// User structure
type User struct {
	TableName string
	Id        int    `orm:"serial_primary_key"`
	Name      string `orm:"varchar not_null"`
	Email     string `orm:"varchar not_null"`
}

// Table of dialogs
type TableDialogs struct {
	BaseTable
}

// Dialog structure
type Dialog struct {
	TableName string
	Id        int    `orm:"serial_primary_key"`
	Owner     string `orm:"varchar not_null"`
	Opponent  string `orm:"varchar not_null"`
}
