package user

type User struct {
	Mobile   string `json:"mobile"`
	Otp      string `json:"otp"`
	isActive bool
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
	createAt string
	updateAt string
}

type Tokens struct {
	UserID           string
	AccessTokenHash  string
	RefreshTokenHash string
	ExpiredAt        string
	CreateAt         string
	UpdateAt         string
}

var userStore = make(map[string]*User) // Удалить позже
