package db

import "RestAPI/orm"

func Register() error {
	er := orm.RegisterModel("users", User{})
	if er != nil {
		return er
	}
	er = orm.RegisterModel("tokens", Token{})
	if er != nil {
		return er
	}
	orm.RegisterModel("images", Image{})
	if er != nil {
		return er
	}
	return nil
}
