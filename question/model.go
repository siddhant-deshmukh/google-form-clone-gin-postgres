package question

import (
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var db *gorm.DB

type Question struct {
	ID          uint           `json:"id"`
	IsRequired  bool           `json:"is_required" gorm:"default:true; not null"`
	QuesType    string         `json:"ques_type" gorm:"check:ques_type IN ('mcq','checkbox','short','long','dropdown','date','time','eliminated');default:'mcq'; not null"`
	Title       string         `json:"title" gorm:"type:varchar(100); check:char_length(title) > 0; default:'Untitled question'; not null"`
	Description string         `json:"description" gorm:"type:varchar(200); default:''; not null"`
	Options     pq.StringArray `json:"options" gorm:"type:text[]; check:array_length(options, 1) < 100"`
	CorrectAns  pq.StringArray `json:"correct_ans" gorm:"type:text[]; check:array_length(correct_ans, 1) < 100"`
	Points      uint           `json:"points" gorm:"check:points > 0; default: 1 "`
	AuthorID    uint           `json:"author_id" gorm:"not null"`
	FormID      uint           `json:"form_id" gorm:"not null"`
	CreatedAt   time.Time      `json:"-" validate:"-"`
	UpdatedAt   time.Time      `json:"-" validate:"-"`

	IndexAt uint `json:"index_at" gorm:"-"`
}
type NewQuestion struct {
	IsRequired  bool           `json:"is_required,omitempty" validate:"boolean"`
	QuesType    string         `json:"ques_type,omitempty" validate:"omitempty,oneof=mcq checkbox short long dropdown date time"`
	Title       string         `json:"title,omitempty" validate:"omitempty,max=100,min=1"`
	Description string         `json:"description,omitempty" validate:"omitempty,max=100"`
	Options     pq.StringArray `json:"options,omitempty" validate:"omitempty,max=100"`
	CorrectAns  pq.StringArray `json:"correct_ans,omitempty" validate:"omitempty,max=100"`
	Points      uint           `json:"points,omitempty" validate:"omitempty,min=1"`
	FormID      uint           `json:"form_id,omitempty" validate:"required,min=0"`
	AuthorID    uint           `json:"author_id,omitempty" validate:"required,min=0"`
	IndexAt     uint           `json:"index_at,omitempty" validate:"required,min=0"`
}
type EditQuestion struct {
	IsRequired  bool     `json:"is_required,omitempty" validate:"boolean"`
	QuesType    string   `json:"ques_type,omitempty" validate:"omitempty,oneof=mcq checkbox short long dropdown date time"`
	Title       string   `json:"title,omitempty" validate:"omitempty,max=100,min=1"`
	Description string   `json:"description,omitempty" validate:"omitempty,max=100"`
	Options     []string `json:"options" validate:"omitempty,max=100"`
	CorrectAns  []string `json:"correct_ans" validate:"omitempty,max=100"`
	Points      uint     `json:"points,omitempty" validate:"omitempty,min=1"`
}
type QuestionWithOutAnswer struct {
	ID          uint     `json:"id"`
	IsRequired  bool     `json:"is_required"`
	QuesType    string   `json:"ques_type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Options     []string `json:"options"`
	Points      uint     `json:"points"`
	FormID      uint     `json:"form_id,omitempty"`
}

func (que *Question) AfterCreate(tx *gorm.DB) (err error) {
	currQueSeq := QueSeq{AuthorID: que.AuthorID}
	result := tx.Where("form_id = ?", que.FormID).First(&currQueSeq)
	if result.Error != nil {
		return result.Error
	}

	if currQueSeq.QuestionSeq == nil {
		currQueSeq.QuestionSeq = []int64{int64(que.ID)}
	} else {
		length := len(currQueSeq.QuestionSeq)
		if length == 0 {
			currQueSeq.QuestionSeq = append([]int64{int64(que.ID)}, currQueSeq.QuestionSeq...)
		} else if length > int(que.IndexAt) {
			currQueSeq.QuestionSeq = append(currQueSeq.QuestionSeq[:int(que.IndexAt)+1], currQueSeq.QuestionSeq[int(que.IndexAt):]...)
			currQueSeq.QuestionSeq[que.IndexAt] = int64(que.ID)
		} else {
			currQueSeq.QuestionSeq = append(currQueSeq.QuestionSeq, int64(que.ID))
		}
	}
	// result = tx.Model(&QueSeq{}).Where("id = ?", currQueSeq.ID).Update("question_seq", []int{int(que.ID)})
	result = tx.Save(&currQueSeq)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (que *Question) BeforeDelete(tx *gorm.DB) (err error) {
	result := tx.Model(&QueSeq{}).Where("form_id = ?", que.FormID).Update("question_seq", gorm.Expr(fmt.Sprintf("array_remove(question_seq, %d)", que.ID)))
	return result.Error
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
