package orm

type BaseCell interface {
	ToFields() ([]interface{}, []string)
}

// Base Table structure
type BaseTable struct {
	TableName string
}

type TableUsers struct {
	BaseTable
}
