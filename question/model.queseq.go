package question

import (
	"time"

	"github.com/lib/pq"
)

type QueSeq struct {
	ID          uint          `json:"id"`
	AuthorID    uint          `json:"author_id" gorm:"not null" validate:"required,min=1"`
	FormID      uint          `json:"form_id" gorm:"index;not null" validate:"required,min=1"`
	QuestionSeq pq.Int64Array `json:"question_seq" gorm:"type:INTEGER[]; default: ARRAY[]::INTEGER[]"`
	CreatedAt   time.Time     `json:"-" validate:"-"`
	UpdatedAt   time.Time     `json:"-" validate:"-"`
}
