package db

import "RestAPI/orm"

type TableUsers struct {
	orm.BaseTable
}

type TableTokens struct {
	orm.BaseTable
}

type TableImages struct {
	orm.BaseTable
}

// Create our models here
type User struct {
	TableName string
	ID        int    `json:"id" orm:"serial_primary_key"`
	Mobile    string `json:"mobile" orm:"varchar unique not_null"`
	Otp       string `json:"otp" orm:"varchar"`
	IsActive  bool   `orm:"bool not_null"`
	Name      string `json:"name" orm:"varchar"`
	Surname   string `json:"surname" orm:"varchar"`
	Age       int    `json:"age" orm:"int"`
	Email     string `json:"email" orm:"varchar unique"`
	Password  string `json:"password" orm:"varchar"`
	CreateAt  string `orm:"timestamp"`
	UpdateAt  string `orm:"timestamp"`
}

type Token struct {
	TableName        string
	ID               int    `json:"id" orm:"serial_primary_key"`
	UserID           string `json:"user_id" orm:"int ref_users(id)"`
	AccessTokenHash  string `json:"access_token_hash" orm:"varchar unique"`
	RefreshTokenHash string `json:"refresh_token_hash" orm:"varchar unique not_null"`
	ExpiredAt        string `json:"expired_at" orm:"timestamp"`
	CreateAt         string `json:"create_at" orm:"timestamp"`
	UpdateAt         string `json:"update_at" orm:"timestamp"`
}

type Image struct {
	TableName string
	ID        int    `json:"id" orm:"serial_primary_key"`
	UserID    string `json:"user_id" orm:"int ref_users(id)"`
	ImageURL  string `json:"image_url" orm:"varchar"`
}
