package question

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/utils"
	"gorm.io/gorm"
)

func newQuestion(c *gin.Context) {
	// binding the request body into newQues
	var newQues NewQuestion
	if err := c.BindJSON(&newQues); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad format of newQues",
		})
		return
	}

	user_id := c.MustGet("user_id").(uint) // gettting user_id from the url
	newQues.AuthorID = user_id

	// validating if the request body is in correct format
	if err_msg, err := utils.ValidateFieldWithStruct(newQues); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err_msg,
		})
		return
	}

	// validating if the formid belongs to the user rquesting this as a safety measure
	author_id, err := utils.GetFormAuthor(newQues.FormID, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Internal Server error (checking form author)",
		})
		return
	} else if author_id != user_id {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   err.Error(),
			"message": "Not allowed",
		})
		return
	}

	// create new question instance
	var question Question
	mapstructure.Decode(newQues, &question)
	results := db.Create(&question)
	if results.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something goes wrong while creating question",
			"error":   results.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Successful",
		"question": question,
		"newQues":  newQues,
	})
}

func editQuestion(c *gin.Context) {
	var editQue EditQuestion

	if err := c.BindJSON(&editQue); err != nil { // binding all the fields
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad format of question",
		})
		return
	}

	qId, err := utils.GetFieldFromUrl(c, "qId") // getting question id from the url
	if err != nil {
		return
	}

	// validate the struct by checking the validate rules mentioned in the validate struct tag key
	if err_msg, err := utils.ValidateFieldWithStruct(editQue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err_msg,
		})
		return
	}

	user_id := c.MustGet("user_id").(uint) // getting user_id from the AuthUserMiddleware

	// converting edit question to a map
	jsonData, err := json.Marshal(editQue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While reading form",
			"err":     err,
		})
		return
	}
	var que_map map[string]interface{}
	err = json.Unmarshal(jsonData, &que_map)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While reading form",
			"err":     err,
		})
		return
	}
	keys := []string{}
	for key, value := range que_map {
		if value != nil {
			keys = append(keys, key)
			if key == "options" {
				que_map[key] = gorm.Expr("ARRAY[ ? ]", utils.ArrayToString(editQue.Options))
			} else if key == "correct_ans" {
				que_map[key] = gorm.Expr("ARRAY[ ? ]", utils.ArrayToString(editQue.CorrectAns))
			}
		}
	}
	// fmt.Println("\n\n", keys, "\n\n", que_map, "\n\n ")

	question := Question{ID: qId, AuthorID: user_id}
	result := db.Model(Question{ID: qId, AuthorID: user_id}).Select(keys).Updates(que_map).Scan(&question)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While updating question",
			"err":     result.Error,
		})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "While updating question",
			"err":     result.Error,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Sucessfull",
		})
	}
}

func deleteQuestion(c *gin.Context) {
	qId, err := utils.GetFieldFromUrl(c, "qId") // getting question id from the url
	if err != nil {
		return
	}

	user_id := c.MustGet("user_id").(uint) // getting user_id from the AuthUserMiddleware

	question := Question{ID: qId, AuthorID: user_id}

	if result := db.First(&question); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While getting question",
			"err":     result.Error.Error(),
		})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "not autherized or invalid question id",
		})
		return
	}

	result := db.Where("author_id = ?", user_id).Delete(&question)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occured",
			"err":     result.Error.Error(),
		})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "not autherized or invalid question id",
		})
	} else {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Successfully deleted",
		})
	}

}

func getQuestions(c *gin.Context) {
	qId, err := utils.GetFieldFromUrl(c, "qId") // getting question id from the url
	if err != nil {
		return
	}

	user_id := c.MustGet("user_id").(uint) // getting user_id from the AuthUserMiddleware

	question := Question{ID: qId}
	result := db.Find(&question)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occured",
			"err":     result.Error.Error(),
		})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "invalid question id",
		})
	} else {
		if question.AuthorID == user_id {
			c.JSON(http.StatusAccepted, gin.H{
				"message":  "Successfully deleted",
				"question": question,
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"message": "Successfully deleted",
				"question": QuestionWithOutAnswer{
					ID:          question.ID,
					IsRequired:  question.IsRequired,
					QuesType:    question.QuesType,
					Title:       question.Title,
					Description: question.Description,
					Points:      question.Points,
					FormID:      question.FormID,
					Options:     question.Options,
				},
			})
		}
	}
}
