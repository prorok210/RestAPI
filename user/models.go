package user


type User struct {
	Mobile 		string  `json:"mobile"`
	Otp 		string 	`json:"otp"`
	isActive 	bool
	Name 		string 	`json:"name"`
	Surname 	string  `json:"surname"`
	Age 		int     `json:"age"`
	createAt 	string 
	updateAt 	string
}

var userStore = make(map[string]*User) // Удалить позже