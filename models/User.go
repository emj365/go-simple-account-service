package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Salt     string `json:"-"`
}

func GetAllUser() []User {
	users := []User{}
	GetDB().Select("Name").Find(&users)
	return users
}

func (u *User) FindForAuth() bool {
	GetDB().Where(User{Name: u.Name}).Select("ID, Name, Password, Salt").Find(&u)
	found := u.Name != ""
	return found
}

func (u *User) Create() error {
	if u.Name == "" || u.Password == "" || u.Salt == "" {
		return errors.New("Name, Password, Salt can not be empty")
	}

	GetDB().Create(u)
	u.Password = "*******"
	return nil
}

func (u *User) NameExistence() bool {
	count := 0
	GetDB().Model(User{}).Where("name = ?", u.Name).Count(&count)
	return count > 0
}
