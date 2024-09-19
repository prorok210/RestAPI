package main

// User structure constructor
func newUser(name string, email string) *User {
	return &User{
		TableName: "users",
		Id:        0, // The system will set the id after the SQL get func to the database
		Name:      name,
		Email:     email,
	}
}

// User structure constructor
func newDialog(owner string, opponent string) *Dialog {
	return &Dialog{
		TableName: "dialogs",
		Id:        0, // The system will set the id after the SQL get func to the database
		Owner:     owner,
		Opponent:  opponent,
	}
}
