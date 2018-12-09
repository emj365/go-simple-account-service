package models

import (
	"errors"
	"time"

	"github.com/emj365/account/lib"
)

type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name"`
	Password  string     `json:"password"`
	Salt      string     `json:"-"`
}

func GetAllUser() []User {
	users := []User{}
	GetDB().Select("id, created_at, updated_at, name").Find(&users)
	return users
}

func FindUserByID(u *User, userID uint) {
	GetDB().Where(userID).Select("id, created_at, updated_at, name").First(u)
}

func (u *User) Auth(password string) bool {
	hash := lib.HashPassword(password, u.Salt)

	if hash == u.Password {
		return true
	}

	return false
}

func (u *User) FindForAuth() bool {
	GetDB().Where(User{Name: u.Name}).Select("id, created_at, updated_at, name, password, salt").Find(&u)
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
