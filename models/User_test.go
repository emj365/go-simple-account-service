package models

import (
	"log"
	"testing"
)

func TestGetAllUser(t *testing.T) {
	GetDB().Unscoped().Delete(&User{})

	user := User{Name: "mike", Password: "passwd"}
	user.Create()

	users := GetAllUser()
	if len(users) != 1 {
		log.Fatalf("TestGetAllUser get wrong almount, which is %v", len(users))
	}

	GetDB().First(&user, user.ID)
	if user.Name != "mike" {
		log.Fatalf("TestGetAllUser user.name: %s", user.Name)
	}
}

func TestCreate(t *testing.T) {
	GetDB().Unscoped().Delete(&User{})

	user := User{Name: "mike", Password: "passwd"}
	err := user.Create()

	if err != nil {
		log.Fatalf("TestCreate get error: %s", err)
	}

	if user.Password != "*******" {
		log.Fatalf("TestCreate user.Password: %s", user.Password)
	}

	if user.Salt != "" {
		log.Fatalf("TestCreate user.Salt: %s", user.Salt)
	}

	var count int
	GetDB().Model(User{}).Count(&count)
	if count != 1 {
		log.Fatalf("TestCreate get wrong users almount, which is %v", count)
	}
}
