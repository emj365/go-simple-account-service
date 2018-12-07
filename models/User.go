package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Salt     string `json:"-"`
}

func FindAllUsers() []User {
	users := []User{}
	GetDB().Select("Name").Find(&users)
	return users
}

func GetUserForAuth(name string) User {
	user := User{}
	GetDB().Where(User{Name: name}).Select("Name, Password, Salt").Find(&user)
	return user
}

func CreateUser(user *User) {
	GetDB().Create(user)
	user.Password = "*******"
}

func CheckUserAlreadyExist(name string) bool {
	count := 0
	GetDB().Model(User{}).Where(User{Name: name}).Count(&count)
	return count > 0
}
