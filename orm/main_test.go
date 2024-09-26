package orm

import (
	"context"
	"fmt"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTable_Success(t *testing.T) {
	// Создаем mock соединение с базой данных
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	// Устанавливаем поведение для проверки существования таблицы
	mock.ExpectQuery("SELECT EXISTS").WithArgs("users").
		WillReturnRows(mock.NewRows([]string{"exists"}).AddRow(false))

	// Устанавливаем поведение для создания таблицы
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").
		WillReturnResult(pgxmock.NewResult("CREATE", 1))
	user := User{TableName: "users"}
	err = CreateTable(user)
	assert.NoError(t, err)
}

func TestCreateTable_TableAlreadyExists(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	// Таблица уже существует
	mock.ExpectQuery("SELECT EXISTS").WithArgs("users").
		WillReturnRows(mock.NewRows([]string{"exists"}).AddRow(true))

	user := User{TableName: "users"}
	err = CreateTable(user)
	assert.NoError(t, err)
}

func TestCreateTable_TableNameFieldNotFound(t *testing.T) {
	// Тестируем случай, когда поле TableName отсутствует
	type NoTableName struct {
		Id   int    `orm:"serial_primary_key"`
		Name string `orm:"varchar not_null"`
	}

	err := CreateTable(NoTableName{})
	assert.EqualError(t, err, "field 'TableName' was not found in the structure")
}

func TestCreateTable_CheckTableExistsError(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	// Ошибка при проверке существования таблицы
	mock.ExpectQuery("SELECT EXISTS").WithArgs("users").
		WillReturnError(fmt.Errorf("query error"))

	user := User{TableName: "users"}
	err = CreateTable(user)
	assert.EqualError(t, err, "error checking table existence: query error")
}

func TestCreateTable_CreateTableError(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	// Таблица не существует
	mock.ExpectQuery("SELECT EXISTS").WithArgs("users").
		WillReturnRows(mock.NewRows([]string{"exists"}).AddRow(false))

	// Ошибка при создании таблицы
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").
		WillReturnError(fmt.Errorf("creation error"))

	user := User{TableName: "users"}
	err = CreateTable(user)
	assert.EqualError(t, err, "error creating table creation error")
}
