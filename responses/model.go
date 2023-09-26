package responses

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/question"
	"gorm.io/gorm"
)

var db *gorm.DB

type Response struct {
	ID        uint       `json:"id"`
	UserEmail string     `json:"user_email,omitempty"`
	FormID    uint       `json:"form_id"`
	Answers   queAnswers `json:"answers" `
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}

type queAnswers map[string][]string

func (answers *queAnswers) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := queAnswers{}
	err := json.Unmarshal(bytes, &result)
	*answers = queAnswers(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (answers queAnswers) Value() (driver.Value, error) {
	return json.Marshal(answers)
}

func (res *Response) BeforeSave(tx *gorm.DB) (err error) {
	var questions []QuestionSnippet

	if result := tx.Model(&question.Question{}).Select("id", "is_required, ques_type").Where("form_id = ?", res.FormID).Find(&questions); result.Error != nil {
		err = result.Error
		fmt.Println(err.Error())
		return
	} else if result.RowsAffected == 0 {
		err = gorm.ErrRecordNotFound
		fmt.Println(err.Error())
		return
	}

	finalQuesAnswers := make(queAnswers)

	for _, que := range questions {
		// 'mcq','checkbox','short','long','dropdown','date','time','eliminated'
		qId := strconv.Itoa(int(que.ID))
		ans, foundRes := res.Answers[qId]

		fmt.Println(qId, que.IsRequired, "\t", foundRes, "\n ")

		if foundRes {
			switch que.QuesType {
			case "mcq", "dropdown":
				if len(ans) == 1 {
					finalQuesAnswers[qId] = ans
				} else if que.IsRequired {
					err = errors.New("doesn't answer required field " + qId)
					return
				}
			case "checkbox":
				if len(ans) >= 1 {
					finalQuesAnswers[qId] = ans
				} else if que.IsRequired {
					err = errors.New("doesn't answer required field " + qId)
					return
				}
			case "short", "long":
				if len(ans) == 1 {
					finalQuesAnswers[qId] = ans
				} else if que.IsRequired {
					err = errors.New("doesn't answer required field " + qId)
					return
				}
			}
		} else if que.IsRequired {
			err = errors.New("doesn't answer required field " + qId)
			return
		}
	}
	return
}

type NewResInput struct {
	UserEmail string     `json:"user_email,omitempty" validate:"omitempty,email"`
	FormID    uint       `json:"form_id" validate:"required"`
	Answers   queAnswers `json:"answers" validate:"required"`
}

type QuestionSnippet struct {
	ID         uint   `json:"id"`
	QuesType   string `json:"ques_type"`
	IsRequired bool   `json:"is_required"`
}

func SetResponseTable(gormDB *gorm.DB) {
	db = gormDB
	db.AutoMigrate(&Response{})
}
