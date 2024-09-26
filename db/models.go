package db

// Create our models here
type User struct {
	TableName string
	ID        int    `json:"id" orm:"serial_primary_key"`
	Mobile    string `json:"mobile" orm:"varchar unique not_null"`
	Otp       string `json:"otp" orm:"int"`
	isActive  bool   `orm:"bool not_null"`
	Name      string `json:"name" orm:"varchar"`
	Surname   string `json:"surname" orm:"varchar"`
	Age       int    `json:"age" orm:"int"`
	Email     string `json:"email" orm:"varchar unique"`
	Password  string `json:"password" orm:"varchar"`
	createAt  string `orm:"timestamp"`
	updateAt  string `orm:"timestamp"`
}

type Tokens struct {
	UserID           string
	AccessTokenHash  string
	RefreshTokenHash string
	ExpiredAt        string
	CreateAt         string
	UpdateAt         string
}

type Image struct {
	UserID   string
	ImageURL string
}
