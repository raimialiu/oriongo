package controllers

import (
	"oriongo/internal/origongo"
)

type (
	UsersController struct {
		app origongo.OrionGo
	}
	User struct {
		Username  string `gorm:"column:username"`
		Email     string `gorm:"column:email"`
		FirstName string `gorm:"column:first_name"`
		LastName  string `gorm:"column:last_name"`
		City      string `gorm:"column:city"`
	}
)

func (u *UsersController) GetUsers() []User {
	var users []User
	u.app.DB().Find(&users)

	return users
}
