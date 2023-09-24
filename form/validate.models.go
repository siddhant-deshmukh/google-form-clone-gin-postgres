package form

type NewForm struct {
	Title            string                `json:"title" validate:"required,min=5,max=100"`
	Description      string                `json:"description,omitempty" validate:"omitempty,max=300"`
	AuthorID         uint                  `json:"author_id" validate:"required"`
	Quiz_Setting     *New_Quiz_Setting     `json:"quiz_setting,omitempty"  validate:"omitempty"`
	Response_Setting *New_Response_Setting `json:"response_setting,omitempty"  validate:"omitempty"`
}
type New_Quiz_Setting struct {
	IsQuiz        bool `json:"is_quiz" validate:"boolean"`
	DefaultPoints uint `json:"default_points" validate:"required,min=1,max=100"`
}
type New_Response_Setting struct {
	CollectEmail bool `json:"collect_email" validate:"boolean"`
	AllowEditRes bool `json:"allow_edit_res" validate:"boolean"`
	SendResCopy  bool `json:"send_res_copy" validate:"boolean"`
}

type EditForm struct {
	Title            string                 `json:"title,omitempty" validate:"omitempty,min=5,max=100"`
	Description      string                 `json:"description,omitempty" validate:"omitempty,max=300"`
	Quiz_Setting     *Edit_Quiz_Setting     `json:"quiz_setting,omitempty"  validate:"omitempty"`
	Response_Setting *Edit_Response_Setting `json:"response_setting,omitempty"  validate:"omitempty"`
}
type Edit_Quiz_Setting struct {
	IsQuiz        bool `json:"is_quiz" validate:"boolean"`
	DefaultPoints uint `json:"default_points" validate:"required,min=1,max=100"`
}
type Edit_Response_Setting struct {
	CollectEmail bool `json:"collect_email" validate:"boolean"`
	AllowEditRes bool `json:"allow_edit_res" validate:"boolean"`
	SendResCopy  bool `json:"send_res_copy" validate:"boolean"`
}
