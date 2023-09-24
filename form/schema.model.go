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

// * Form schema
type Form struct {
	ID               uint                `json:"id" gorm:"primaryKey;autoIncrement;<-:create" `
	AuthorID         uint                `json:"author_id" gorm:"not null"`
	Title            string              `json:"title" gorm:"type:varchar(100);not null;check: char_length(title) > 5"`
	Description      string              `json:"description" gorm:"type:varchar(300);not null"`
	Quiz_Setting     *Quiz_Setting       `json:"quiz_setting" gorm:"type:jsonb"`
	Response_Setting *Response_Setting   `json:"response_setting" gorm:"type:jsonb"`
	CreatedAt        time.Time           `json:"-"`
	UpdatedAt        time.Time           `json:"-"`
	Questions        []question.Question `gorm:"foreignKey:FormID;references:ID"`
	QueSeq           question.QueSeq     `gorm:"foreignKey:id"`
}

func (form *Form) BeforeCreate(tx *gorm.DB) (err error) {
	if form.Quiz_Setting == nil {
		form.Quiz_Setting = &Quiz_Setting{IsQuiz: true, DefaultPoints: 1}
	}
	if form.Response_Setting == nil {
		form.Response_Setting = &Response_Setting{CollectEmail: false, AllowEditRes: false, SendResCopy: false}
	}
	return nil
}
func (form *Form) AfterCreate(tx *gorm.DB) (err error) {
	que := question.QueSeq{
		AuthorID:    form.AuthorID,
		FormID:      form.ID,
		QuestionSeq: []int64{},
	}

	if result := tx.Create(&que); result.Error != nil {
		return result.Error
	}
	return nil
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
