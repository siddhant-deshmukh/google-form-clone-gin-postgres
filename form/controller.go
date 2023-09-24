package form

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/question"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/utils"
	"gorm.io/gorm"
)

func createForm(c *gin.Context) {
	var newForm NewForm
	user_id := c.MustGet("user_id").(uint)
	err := c.BindJSON(&newForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad form data format",
			"err":     err,
		})
		return
	}

	newForm.AuthorID = user_id

	if res_msg, err := utils.ValidateFieldWithStruct(newForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": res_msg,
			"err":     err,
		})
		return
	}

	var form Form
	mapstructure.Decode(newForm, &form)
	result := db.Create(&form)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While creating data",
			"err":     result.Error,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"newform": newForm,
		"form":    form,
	})
}

func getFormById(c *gin.Context) {
	var isAuthor bool
	var form Form
	user_id := c.MustGet("user_id").(uint)

	formId, err := utils.GetFieldFromUrl(c, "id")
	if err != nil {
		return
	}
	result := db.Model(&Form{ID: formId}).First(&form)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While creating result",
			"err":     result.Error,
		})
		return
	}
	isAuthor = user_id == form.AuthorID

	if isAuthor {
		c.JSON(http.StatusCreated, gin.H{
			"form": form,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"form": gin.H{
				"ID":           form.ID,
				"Title":        form.Title,
				"Description":  form.Description,
				"Quiz_Setting": form.Quiz_Setting,
			},
		})
	}
}

func editForm(c *gin.Context) {
	var form EditForm
	user_id := c.MustGet("user_id").(uint)

	formId, err := utils.GetFieldFromUrl(c, "id")
	if err != nil {
		return
	}

	err = c.BindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad form data format",
			"err":     err,
		})
		return
	}
	if res_msg, err := utils.ValidateFieldWithStruct(form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": res_msg,
			"err":     err,
		})
		return
	}

	jsonData, err := json.Marshal(form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While reading form",
			"err":     err,
		})
		return
	}

	var form_map map[string]interface{}
	err = json.Unmarshal(jsonData, &form_map)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While reading form",
			"err":     err,
		})
		return
	}
	keys := []string{}
	for key, value := range form_map {
		if eq := reflect.DeepEqual(value, Quiz_Setting{}); !eq {
			keys = append(keys, key)
		}
	}
	fmt.Println(keys)
	fmt.Println()
	fmt.Println(form_map)

	result := db.Model(Form{}).Select(keys).Where("id = ? AND author_id = ?", formId, user_id).Updates(form_map)
	// result := db.Model(Form{}).Select("title").Where "send_res_copy": true("id = ? AND author_id = ?", formId, user_id).Updates(map[string]interface{}{"title": "Meow"})

	// result := db.First(&form, formId)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While updating row",
			"err":     result.Error,
			"keys":    keys,
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "While updating row",
			"err":     result.Error,
		})
		return
	}
	// if user_id != form.AuthorID {
	// 	c.JSON(http.StatusForbidden, gin.H{
	// 		"message": "Permission Denied",
	// 	})
	// 	return
	// }
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Updated",
		"form":    form,
	})
}

func deleteForm(c *gin.Context) {
	var form Form
	user_id := c.MustGet("user_id").(uint)

	formId, err := utils.GetFieldFromUrl(c, "id")
	if err != nil {
		return
	}

	result := db.First(&form, formId)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While creating result",
			"err":     result.Error,
		})
		return
	}

	if user_id != form.AuthorID {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission Denied",
		})
		return
	}
}

func getQuestions(c *gin.Context) {
	user_id := c.MustGet("user_id").(uint)

	formId, err := utils.GetFieldFromUrl(c, "id")
	if err != nil {
		return
	}

	var question_seq question.QueSeq

	if result := db.Where("form_id = ?", formId).First(&question_seq); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":     result.Error.Error(),
			"message": "Something went wrong",
		})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"err":     "Not found",
			"message": "Questions list not found",
		})
		return
	}

	if question_seq.AuthorID != user_id {
		questions := make(map[int]question.QuestionWithOutAnswer)
		for _, qId := range question_seq.QuestionSeq {
			temp_ques := question.QuestionWithOutAnswer{}
			if result := db.Model(&question.Question{ID: uint(qId)}).First(&temp_ques); result.Error == nil && result.RowsAffected > 0 {
				questions[int(qId)] = temp_ques
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"questions_seq": question_seq.QuestionSeq,
			"questions":     questions,
		})
	} else {
		questions := make(map[int]question.Question)
		for _, qId := range question_seq.QuestionSeq {
			temp_ques := question.Question{ID: uint(qId)}
			if result := db.First(&temp_ques); result.Error == nil && result.RowsAffected > 0 {
				questions[int(qId)] = temp_ques
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"questions_seq": question_seq.QuestionSeq,
			"questions":     questions,
		})
	}

}

func postQuestions(c *gin.Context) {

	var questionSeq []int
	// var oldQuestionSeq []int

	user_id := c.MustGet("user_id").(uint)
	form_id, err := utils.GetFieldFromUrl(c, "id")
	if err != nil {
		return
	}

	if err := c.BindJSON(&questionSeq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect format of input",
			"err":     err.Error(),
		})
		return
	}

	// queSeq := question.QueSeq{AuthorID: user_id, FormID: form_id}

	var queSeq question.QueSeq
	if result := db.Where("author_id = ? AND form_id = ?", user_id, form_id).First(&queSeq); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something goes wrong",
			"err":     result.Error.Error(),
		})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Incorrect format of input",
			"err":     "not found",
			"q":       queSeq,
		})
		return
	}

	if !areArrayEqual(queSeq.QuestionSeq, questionSeq) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "doesn't match old sequance",
		})
		return
	}

	if result := db.Model(question.QueSeq{}).Where("author_id = ? AND form_id = ?", user_id, form_id).Update("question_seq", gorm.Expr("ARRAY[ "+IntArrayToString(questionSeq)+"]")); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something goes wrong",
			"err":     result.Error.Error(),
		})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Incorrect format of input",
			"err":     "not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Updated!",
	})
}

func SetFormTable(gormDB *gorm.DB) {
	db = gormDB
	db.AutoMigrate(&Form{})
}

func areArrayEqual(arr1 []int64, arr2 []int) bool {
	var arr3 []int
	for _, v := range arr1 {
		arr3 = append(arr3, int(v)) // cast each element of a to int and append it to b
	}
	sort.Ints(arr3)
	sort.Ints(arr2)

	for i, v := range arr3 {
		if arr2[i] != v {
			return false
		}
	}
	return true
}

func IntArrayToString(arr []int) string {
	if len(arr) <= 0 {
		return ""
	} else if len(arr) == 1 {
		return strconv.Itoa(arr[0])
	} else if len(arr) == 2 {
		return strconv.Itoa(arr[0]) + ", " + strconv.Itoa(arr[1])
	} else {

		str := strconv.Itoa(arr[0]) + ", "
		for i := 1; i < len(arr)-1; i++ {
			str += strconv.Itoa(arr[i]) + ", "
		}
		return str + strconv.Itoa(arr[len(arr)-1])
	}
}
