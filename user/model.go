package user

import (
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/form"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	ID       uint
	Name     string      `gorm:"type:varchar(20); not null"`
	Email    string      `gorm:"type:varchar(100);uniqueIndex; not null;check:email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$' "`
	Password string      `gorm:"type:varchar(100); not null"`
	Forms    []form.Form `gorm:"foreignKey:AuthorID;references:ID"`
}
type UserCreateForm struct {
	Name     string `gorm:"type:varchar(20)"`
	Email    string `gorm:"type:varchar(100);uniqueIndex"`
	Password string `gorm:"type:varchar(100)"`
}

func SetUserTable(gormDB *gorm.DB) {
	db = gormDB
	db.AutoMigrate(&User{})
}