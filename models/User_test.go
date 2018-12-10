package models

import (
	"testing"
)

var userName = "mike"

func deleteUsers() {
	GetDB().Unscoped().Delete(&User{})
}

func createUser() (User, error) {
	user := User{Name: userName, Password: "whatever"}
	err := user.Create()
	return user, err
}

func TestGetAllUser(t *testing.T) {
	deleteUsers()
	createUser()

	users := GetAllUser()

	if want := 1; len(users) != want {
		t.Errorf("len(users) == %v, want %v", len(users), want)
	}

	user := &users[0]
	if want := userName; user.Name != want {
		t.Errorf("user.Name == %s, want %s", user.Name, want)
	}
}

func TestFindUserByID(t *testing.T) {
	deleteUsers()
	user, _ := createUser()

	foundUser := User{}
	FindUserByID(&foundUser, user.ID)

	if want := userName; foundUser.Name != want {
		t.Errorf("foundUser.Name == %s, want %s", foundUser.Name, want)
	}
}

func TestCreate(t *testing.T) {
	deleteUsers()
	user, err := createUser()

	if err != nil {
		t.Errorf("err == %s, want %v", err, nil)
	}

	if want := "*******"; user.Password != want {
		t.Errorf("user.Password == %s, want %s", err, want)
	}

	if want := ""; user.Salt != want {
		t.Errorf("user.Salt == %s, want %s", user.Salt, want)
	}

	var count int
	GetDB().Model(User{}).Count(&count)
	if want := 1; count != want {
		t.Errorf("count == %v, want %v", count, want)
	}
}

func TestNameExistence(t *testing.T) {
	deleteUsers()
	user, _ := createUser()

	if existence, want := user.NameExistence(), true; existence != want {
		t.Errorf("user.NameExistence == %v, want %v", existence, want)
	}

	user.Name = "none"
	if existence, want := user.NameExistence(), false; existence != want {
		t.Errorf("user.NameExistence == %v, want %v", existence, want)
	}
}
