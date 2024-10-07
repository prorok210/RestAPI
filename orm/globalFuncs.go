package orm

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func Contains(slice []string, input string) bool {
	for _, str := range slice {
		if strings.Contains(str, input) {
			return true
		}
	}
	return false
}

// The function takes a string and returns it with the 1st character in large case
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Функция для поиска поля независимо от регистра
func FindFieldByName(t reflect.Type, fieldName string) (reflect.StructField, error) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if strings.EqualFold(field.Name, fieldName) {
			fmt.Println(reflect.ValueOf(field))
			return field, nil
		}
	}
	return reflect.StructField{}, fmt.Errorf("field %s not found", fieldName)
}
