package main

// User structure constructor
func newUser(name string, email string) *User {
	return &User{
		TableName: "users",
		ID:        0, // The system will set the id after the SQL get func to the database
		Name:      name,
		Email:     email,
	}
}
