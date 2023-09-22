package form

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/question"
	"gorm.io/gorm"
)

var db *gorm.DB

type Form struct {
	ID               uint                `gorm:"primaryKey;autoIncrement;<-:create" json:"id"`
	AuthorID         uint                `json:"author_id" gorm:"not null"`
	Title            string              `gorm:"type:varchar(100);not null;check: char_length(title) > 5" json:"title" validate:"min=3,max=100,required"`
	Description      string              `gorm:"type:varchar(300);not null" json:"description" validate:"max=300"`
	Quiz_Setting     Quiz_Setting        `gorm:"type:jsonb" json:"quiz_setting"`
	Response_Setting Response_Setting    `gorm:"type:jsonb" json:"response_setting"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	Questions        []question.Question `gorm:"foreignKey:FormID;references:ID"`
	QueSeq           question.QueSeq     `gorm:"foreignKey:id"`
}

type Quiz_Setting struct {
	IsQuiz        bool `gorm:"default:true" json:"is_quiz"`
	DefaultPoints uint `gorm:"default:1; check: default_points > 0" json:"default_points" validate:"omitempty,min=1,max=100"`
}

func (quizeS *Quiz_Setting) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := Quiz_Setting{}
	err := json.Unmarshal(bytes, &result)
	*quizeS = Quiz_Setting(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (quizeS Quiz_Setting) Value() (driver.Value, error) {
	return json.Marshal(quizeS)
}

type Response_Setting struct {
	CollectEmail bool `gorm:"default:true" json:"collect_email"`
	AllowEditRes bool `gorm:"default:true" json:"allow_edit_res"`
	SendResCopy  bool `gorm:"default:true" json:"send_res_copy"`
}

func (resS *Response_Setting) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := Response_Setting{}
	err := json.Unmarshal(bytes, &result)
	*resS = Response_Setting(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (resS Response_Setting) Value() (driver.Value, error) {
	return json.Marshal(resS)
}

func SetFormTable(gormDB *gorm.DB) {
	db = gormDB
	db.AutoMigrate(&Form{})
}

type EditForm struct {
	Title            string                `json:"title,omitempty" validate:"omitempty,min=5,max=100"`
	Description      string                `json:"description,omitempty" validate:"omitempty,max=300"`
	Quiz_Setting     Edit_Quiz_Setting     `json:"quiz_setting,omitempty"  validate:"omitempty"`
	Response_Setting Edit_Response_Setting `json:"response_setting,omitempty"  validate:"omitempty"`
}
type Edit_Quiz_Setting struct {
	IsQuiz        bool `json:"is_quiz,omitempty" validate:"omitempty"`
	DefaultPoints uint `json:"default_points,omitempty" validate:"omitempty,min=1,max=100"`
}
type Edit_Response_Setting struct {
	CollectEmail bool `json:"collect_email,omitempty" validate:"omitempty"`
	AllowEditRes bool `json:"allow_edit_res,omitempty" validate:"omitempty"`
	SendResCopy  bool `json:"send_res_copy,omitempty" validate:"omitempty"`
}

func (form *Form) AfterCreate(tx *gorm.DB) (err error) {
	que := question.QueSeq{
		AuthorID:    form.AuthorID,
		FormID:      form.ID,
		QuestionSeq: []uint{},
	}

	if result := tx.Create(&que); result.Error != nil {
		return result.Error
	}
	return nil
}
