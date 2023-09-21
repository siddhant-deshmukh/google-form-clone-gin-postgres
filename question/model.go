package question

import (
	"os/user"

	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/form"
	"gorm.io/gorm"
)

var db *gorm.DB

type Question struct {
	ID          uint      `json:"id"`
	IsRequired  bool      `json:"is_required,omitempty" gorm:"default:true; not null"`
	Type        string    `json:"type" gorm:"check:type IN ('mcq','checkbox','short','long','dropdown','date','time'); not null" validate:"oneof=mcq checkbox short long dropdown date time"`
	Title       string    `json:"title" gorm:"type:varchar(100); check:char_length(title) > 5; default:'Untitled question'; not null"`
	Description string    `json:"description" gorm:"default:''; not null"`
	Options     *[]string `json:"options,omitempty" gorm:"type:varchar(30)[]; check:array_length(options, 1) < 100" validate:"max=100"`
	CorrectAns  *[]string `json:"correct_ans,omitempty" gorm:"type:varchar(30)[]; check:array_length(correct_ans, 1) < 100" validate:"max=100"`
	Points      uint      `json:"points,omitempty" gorm:"check:points > 0"`
	AuthorID    uint      `json:"author_id" gorm:"not null"`
	FormID      uint      `json:"form_id" gorm:"not null"`
	User        user.User `gorm:"references:id"`
	Form        form.Form `gorm:"references:id"`
}
type NewQuestion struct {
	IsRequired  bool      `json:"is_required,omitempty"`
	Type        string    `json:"type" validate:"oneof=mcq checkbox short long dropdown date time"`
	Title       string    `json:"title,omitempty" validate:"omitempty,max=100,min=5"`
	Description string    `json:"description,omitempty" validate:"omitempty,max=100"`
	Options     *[]string `json:"options,omitempty" validate:"omitempty,max=100"`
	CorrectAns  *[]string `json:"correct_ans,omitempty" validate:"omitempty,max=100"`
	Points      uint      `json:"points,omitempty" validate:"omitempty,min=1"`
	FormID      uint      `json:"form_id" validate:"required,min=1"`
}

func SetQuestionTable(gormDB *gorm.DB) {
	db = gormDB
	db.AutoMigrate(Question{})
}
