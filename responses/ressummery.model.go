package responses

type ResSummery struct {
	ID         uint            `json:"id"`
	QuestionID uint            `json:"question_id" gorm:"index"`
	FormID     uint            `json:"form_id"`
	ResCount   uint            `json:"res_count"`
	Summery    map[string]uint `json:"summery" gorm:"type:jsonb"`
}
