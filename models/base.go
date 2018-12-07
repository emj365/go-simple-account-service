package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // dependency
)

var db *gorm.DB //database

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Salt     string `json:"-"`
}

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}

	db.Debug().AutoMigrate(&User{}) //Database migration
	db.Debug().Model(&User{}).AddUniqueIndex("idx_user_name", "name")
}

// CloseDB close db connection
func CloseDB() {
	db.Close()
}

// GetDB returns a handle to the DB object
func GetDB() *gorm.DB {
	return db.Debug()
}
