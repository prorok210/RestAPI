package orm

type BaseCell interface {
	ToFields() ([]interface{}, []string)
}
