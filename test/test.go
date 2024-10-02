package main

import (
	"fmt"
	"reflect"
)

type User struct {
	BaseCell
	B string
}
type BaseCell struct {
	A string
}

func main() {
	a := User{}
	b := BaseCell{A: "3"}
	a.allFields()
	b.allFields()
}

func (bs BaseCell) allFields() {
	fields := make([]string, 0)
	val := reflect.ValueOf(bs)
	for i := 0; i < val.NumField(); i++ {
		fields = append(fields, val.Type().Field(i).Name)
	}
	fmt.Println(fields)
}
