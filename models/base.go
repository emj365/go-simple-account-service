package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB //database

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Salt     string `json:"Salt"`
}

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}

	db.Debug().AutoMigrate(&User{}) //Database migration
}

func CloseDB() {
	db.Close()
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
