package question

import (
	"log"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

type Question struct {
	ID          uint      `json:"id"`
	IsRequired  bool      `json:"is_required,omitempty" gorm:"default:true; not null"`
	QuesType    string    `json:"ques_type" gorm:"default:'mcq';check:ques_type IN ('mcq','checkbox','short','long','dropdown','date','time'); not null" validate:"omitempty,oneof=mcq checkbox short long dropdown date time"`
	Title       string    `json:"title" gorm:"type:varchar(100); check:char_length(title) > 0; default:'Untitled question'; not null"  validate:"omitempty,max=100"`
	Description string    `json:"description" gorm:"default:''; not null"  validate:"omitempty,max=100"`
	Options     *[]string `json:"options,omitempty" gorm:"type:varchar(30)[]; check:array_length(options, 1) < 100" validate:"omitempty,max=100"`
	CorrectAns  *[]string `json:"correct_ans,omitempty" gorm:"type:varchar(30)[]; check:array_length(correct_ans, 1) < 100" validate:"omitempty,max=100"`
	Points      uint      `json:"points,omitempty" gorm:"check:points > 0"`
	AuthorID    uint      `json:"author_id" gorm:"not null" validate:"required,min=1"`
	FormID      uint      `json:"form_id" gorm:"not null" validate:"required,min=1"`
	// User        user.User `json:"-" validate:"-" gorm:"foreignKey:AuthorID;references:ID"`
	// Form        form.Form `json:"-" validate:"-" gorm:"foreignKey:FormID;references:ID"`
	CreatedAt time.Time `json:"-" validate:"-"`
	UpdatedAt time.Time `json:"-" validate:"-"`
}
type NewQuestion struct {
	IsRequired  bool      `json:"is_required,omitempty"`
	Type        string    `json:"type" validate:"oneof=mcq checkbox short long dropdown date time"`
	Title       string    `json:"title,omitempty" validate:"omitempty,max=100,min=1"`
	Description string    `json:"description,omitempty" validate:"omitempty,max=100"`
	Options     *[]string `json:"options,omitempty" validate:"omitempty,max=100"`
	CorrectAns  *[]string `json:"correct_ans,omitempty" validate:"omitempty,max=100"`
	Points      uint      `json:"points,omitempty" validate:"omitempty,min=1"`
	FormID      uint      `json:"form_id" validate:"required,min=1"`
	AuthorID    uint      `json:"author_id" validate:"required,min=1"`
}

func SetQuestionTable(gormDB *gorm.DB) {
	db = gormDB
	err1 := db.AutoMigrate(&Question{})
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	err2 := db.AutoMigrate(&QueSeq{})
	if err2 != nil {
		log.Fatal(err2.Error())
	}
}
